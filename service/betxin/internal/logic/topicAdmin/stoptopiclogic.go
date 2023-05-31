package topicAdmin

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

type StopTopicLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewStopTopicLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StopTopicLogic {
	return &StopTopicLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *StopTopicLogic) StopTopic(req *types.StopTopicReq) error {
	topic, err := l.svcCtx.TopicModel.FindOneByTid(l.ctx, req.Tid)
	if err != nil {
		logx.Errorw("TopicModel.FindOneByTid", logx.LogField{Key: "Err", Value: err.Error()})
		return errorx.NewCodeError(1003, "Error")
	}
	topic.IsStop = 1

	err = l.svcCtx.TopicModel.Update(l.ctx, topic)
	if err != nil {
		logx.Errorw("TopicModel.Stop", logx.LogField{Key: "Err", Value: err.Error()})
		return errorx.NewCodeError(1003, "Error")
	}

	purchases, _ := l.svcCtx.TopicPurchaseModel.ListByTid(l.ctx, topic.Tid)
	switch req.WhichWin {
	case 1:
	}

	WinPurchases := make([]model.Topicpurchase, 0)

	if req.WhichWin == 0 {
		for _, purchase := range purchases {
			if purchase.YesPrice.GreaterThan(decimal.NewFromInt(0)) {
				WinPurchases = append(WinPurchases, purchase)
			}
		}
	} else {
		for _, purchase := range purchases {
			if purchase.NoPrice.GreaterThan(decimal.NewFromInt(0)) {
				WinPurchases = append(WinPurchases, purchase)
			}
		}
	}

	amount := topic.TotalPrice.Div(decimal.NewFromInt(int64(len(WinPurchases))))

	for _, winpurchase := range WinPurchases {
		l.svcCtx.MixinSrvRPC.SendMessage(l.ctx, &pb.SendMessageReq{Content: fmt.Sprintf("您预测的话题 %s 正确,发放奖金", topic.Title), ReceiptId: winpurchase.Uid})

		if winpurchase.YesPrice.GreaterThan(decimal.NewFromInt(0)) {
			l.svcCtx.MixinSrvRPC.SendTransfer(l.ctx, &pb.SendTransferReq{
				OpponentId: winpurchase.Uid,
				AssetId:    constant.CNB_ASSET_ID,
				Amount:     amount.String(),
			})
		}
	}

	return nil
}
