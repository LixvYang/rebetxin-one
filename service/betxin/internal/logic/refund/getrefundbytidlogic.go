package refund

import (
	"context"

	"github.com/lixvyang/rebetxin-one/service/betxin/internal/svc"
	"github.com/lixvyang/rebetxin-one/service/betxin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetRefundByTidLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetRefundByTidLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRefundByTidLogic {
	return &GetRefundByTidLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetRefundByTidLogic) GetRefundByTid(req *types.GetRefundByTidReq) (resp *types.Refund, err error) {
	// todo: add your logic here and delete this line

	return
}
