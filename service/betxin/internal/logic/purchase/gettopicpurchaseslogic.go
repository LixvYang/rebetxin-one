package purchase

import (
	"context"

	"github.com/lixvyang/rebetxin-one/service/betxin/internal/svc"
	"github.com/lixvyang/rebetxin-one/service/betxin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetTopicPurchasesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetTopicPurchasesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTopicPurchasesLogic {
	return &GetTopicPurchasesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetTopicPurchasesLogic) GetTopicPurchases() (resp *types.GetTopicpurchasesResp, err error) {
	// todo: add your logic here and delete this line

	return
}
