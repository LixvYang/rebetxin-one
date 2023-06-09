package refund

import (
	"context"
	"fmt"

	"github.com/lixvyang/rebetxin-one/common/constant"
	"github.com/lixvyang/rebetxin-one/common/errorx"
	"github.com/lixvyang/rebetxin-one/service/betxin/internal/svc"
	"github.com/lixvyang/rebetxin-one/service/betxin/internal/types"
	"github.com/lixvyang/rebetxin-one/service/betxin/model"
	"github.com/lixvyang/rebetxin-one/service/mixinsrv/rpc/pb"
	"github.com/shopspring/decimal"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateRefundLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateRefundLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateRefundLogic {
	return &CreateRefundLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateRefundLogic) CreateRefund(req *types.CreateRefundReq) error {
	uid := fmt.Sprintf("%s", l.ctx.Value("uid"))

	// 首先查询，一下对应的Amount不能比用户原有的高
	userPurchase, err := l.svcCtx.TopicPurchaseModel.FindOneByUidTid(l.ctx, uid, req.Tid)
	if err != nil {
		return errorx.NewDefaultError("Error: 找不到用户的购买记录")
	}

	refundAmount, err := decimal.NewFromString(req.Amount)
	if err != nil {
		logx.Errorw("Error", logx.LogField{Key: "Error: ", Value: err.Error()})
		return errorx.NewDefaultError("Error: Refund Amount Format Error!")
	}

	totalAmount := refundAmount.Div(decimal.NewFromFloat(0.9))

	// 尝试先减少话题总价格和Yes价格或者No价格，然后更新用户购买表
	// 查询话题信息
	topic, err := l.svcCtx.TopicModel.FindOneByTid(l.ctx, req.Tid)
	if err != nil {
		logx.Errorw("TopicModel.FindOneByTid", logx.LogField{Key: "Error", Value: err.Error()})
		return errorx.NewDefaultError("Error: refund too much.")
	}

	if totalAmount.GreaterThan(topic.TotalPrice) {
		logx.Errorw("totalAmount.GreaterThan(topic.TotalPrice)", logx.LogField{Key: "Error", Value: err.Error()})
		return errorx.NewDefaultError("Error: refund too much.")
	}

	// 1. 判断数额是否合法，接着判断
	// 判断将 totalAmount 小于
	// 最后会扣除10%转账给用户
	switch req.Select {
	case 0:
		if userPurchase.YesPrice.LessThanOrEqual(totalAmount) {
			logx.Errorw("YesPrice.LessThanOrEqual", logx.LogField{Key: "Error: ", Value: err.Error()})
			return errorx.NewDefaultError("Error: refund too much.")
		}
	case 1:
		if userPurchase.NoPrice.LessThanOrEqual(totalAmount) {
			logx.Errorw("NoPrice.LessThanOrEqual", logx.LogField{Key: "Error: ", Value: err.Error()})
			return errorx.NewDefaultError("Error: refund too much.")
		}
	}

	err = l.HandleRefundLogic(req.Select, userPurchase, topic, refundAmount)
	if err != nil {
		logx.Errorw("HandleRefundLogic Error", logx.LogField{Key: "Error: ", Value: err.Error()})
		return errorx.NewDefaultError("Error: refund.")
	}
	// 给用户转账
	// 转账逻辑
	in := &pb.SendTransferReq{
		OpponentId: uid,
		AssetId:    constant.CNB_ASSET_ID,
		Amount:     refundAmount.String(),
		Memo:       fmt.Sprintf("%s - 退款", topic.Title),
	}
	l.svcCtx.MixinSrvRPC.SendTransfer(l.ctx, in)

	// 向其他用户转账
	feeAmount := totalAmount.Mul(decimal.NewFromFloat(0.1))
	tps, err := l.svcCtx.TopicPurchaseModel.ListByTid(l.ctx, req.Tid)
	if err != nil {
		logx.Errorw("TopicPurchaseModel.ListByTid", logx.LogField{Key: "Error: ", Value: err.Error()})
		return errorx.NewDefaultError("TopicPurchaseModel.ListByTid: refund.")
	}

	if len(tps) <= 1 {
		logx.Info("len(tps) <= 1")
		return nil
	}

	peerAmount := feeAmount.Div(decimal.NewFromInt(int64(len(tps) - 1))).String()
	for i := 0; i < len(tps); i++ {
		if tps[i].Uid == uid {
			continue
		}
		// 转账逻辑
		in := &pb.SendTransferReq{
			OpponentId: tps[i].Uid,
			AssetId:    constant.CNB_ASSET_ID,
			Amount:     peerAmount,
			Memo:       fmt.Sprintf("%s - 退款奖励", topic.Title),
		}
		l.svcCtx.MixinSrvRPC.SendTransfer(l.ctx, in)
	}

	return nil
}

func (l *CreateRefundLogic) HandleRefundLogic(Select int64, userPurchase *model.Topicpurchase, topic *model.Topic, refundAmount decimal.Decimal) error {
	switch Select {
	case 0:
		// 更新用户购买系统逻辑
		userPurchase.YesPrice = userPurchase.YesPrice.Sub(refundAmount)
		err := l.svcCtx.TopicPurchaseModel.Update(l.ctx, userPurchase)
		if err != nil {
			logx.Errorw("TopicPurchaseModel.Update", logx.LogField{Key: "Error: ", Value: err.Error()})
			return errorx.NewDefaultError("TopicPurchaseModel.Update Error!")
		}
		// 更新话题逻辑
		topic.YesPrice = topic.YesPrice.Sub(refundAmount)
		topic.TotalPrice = topic.TotalPrice.Sub(refundAmount)
		topic.YesRatio = topic.YesPrice.Div(topic.TotalPrice).Mul(decimal.NewFromInt(100))
		topic.NoRatio = topic.NoPrice.Div(topic.TotalPrice).Mul(decimal.NewFromInt(100))
		fmt.Println("topicModel: ", topic)
		err = l.svcCtx.TopicModel.Update(l.ctx, topic)
		if err != nil {
			logx.Errorw("TopicModel.Update", logx.LogField{Key: "Error: ", Value: err.Error()})
			return errorx.NewDefaultError("Topic.Update Error!")
		}
	case 1:
		// 更新用户购买系统逻辑
		userPurchase.NoPrice = userPurchase.NoPrice.Sub(refundAmount)
		err := l.svcCtx.TopicPurchaseModel.Update(l.ctx, userPurchase)
		if err != nil {
			logx.Errorw("TopicPurchaseModel.Update", logx.LogField{Key: "Error: ", Value: err.Error()})
			return errorx.NewDefaultError("TopicPurchaseModel.Update Error!")
		}
		topic.NoPrice = topic.NoPrice.Sub(refundAmount)
		topic.TotalPrice = topic.TotalPrice.Sub(refundAmount)
		topic.YesRatio = topic.YesPrice.Div(topic.TotalPrice).Mul(decimal.NewFromInt(100))
		topic.NoRatio = topic.NoPrice.Div(topic.TotalPrice).Mul(decimal.NewFromInt(100))
		fmt.Printf("topicModel: %+v", topic)
		err = l.svcCtx.TopicModel.Update(l.ctx, topic)
		if err != nil {
			logx.Errorw("TopicModel.Update", logx.LogField{Key: "Error: ", Value: err.Error()})
			return errorx.NewDefaultError("Topic.Update Error!")
		}
	}
	return nil
}
