package refund

import (
	"context"

	"github.com/lixvyang/rebetxin-one/service/betxin/internal/svc"
	"github.com/lixvyang/rebetxin-one/service/betxin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListRefundLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListRefundLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListRefundLogic {
	return &ListRefundLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListRefundLogic) ListRefund() (resp []types.Refund, err error) {
	// todo: add your logic here and delete this line

	return
}
