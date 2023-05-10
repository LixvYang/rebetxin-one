package refund

import (
	"context"
	"fmt"

	"github.com/lixvyang/rebetxin-one/common/convert"
	"github.com/lixvyang/rebetxin-one/common/errorx"
	"github.com/lixvyang/rebetxin-one/internal/svc"
	"github.com/lixvyang/rebetxin-one/internal/types"
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

	// 最后会扣除10%转账给用户
	if req.Select == 0 {
		if userPurchase.YesPrice.Mul(decimal.NewFromFloat(0.9)).LessThan(convert.String2Decimal(req.Amount)) {
			return errorx.NewDefaultError("Error: refund too much.")
		}
	} else {
		if userPurchase.NoPrice.Mul(decimal.NewFromFloat(0.9)).LessThan(convert.String2Decimal(req.Amount)) {
			return errorx.NewDefaultError("Error: refund too much.")
		}
	}

	// 尝试先减少话题总价格和Yes价格或者No价格，然后更新用户购买表
	// 发送添加更新话题购买系统 这里的具体逻辑由话题购买系统做逻辑

	// 最后再给用户转账

	// 接着查询购买此话题的用户有哪些，除了退款者，我们需要给其他人都进行对应额度(平均/按比例)的转账

	// 查询退款，如果曾经已经退过款，那就把退款增加
	// 如果没有退款，那就创建一个退款

	return nil
}
