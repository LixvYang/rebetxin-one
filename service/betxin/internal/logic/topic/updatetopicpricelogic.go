package topic

import (
	"context"

	"github.com/lixvyang/rebetxin-one/service/betxin/internal/svc"
	"github.com/lixvyang/rebetxin-one/service/betxin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateTopicPriceLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateTopicPriceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateTopicPriceLogic {
	return &UpdateTopicPriceLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateTopicPriceLogic) UpdateTopicPrice(req *types.UpdateTopicPriceReq) error {
	// todo: add your logic here and delete this line

	return nil
}
