package purchase

import (
	"context"

	"github.com/lixvyang/rebetxin-one/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreateTopicPurchaseLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateTopicPurchaseLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateTopicPurchaseLogic {
	return &CreateTopicPurchaseLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateTopicPurchaseLogic) CreateTopicPurchase() error {
	// todo: add your logic here and delete this line

	return nil
}
