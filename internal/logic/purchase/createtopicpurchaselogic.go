package purchase

import (
	"context"

	"github.com/lixvyang/rebetxin-one/internal/svc"
	"github.com/lixvyang/rebetxin-one/internal/types"

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

// Create
func (l *CreateTopicPurchaseLogic) CreateTopicPurchase(req *types.CreateTopicPurchaseReq) error {
	// purchase, err := l.svcCtx.TopicPurchaseModel.FindOneByUidTid(l.ctx, req.Uid, req.Uid)
	// if err != nil {
	// 	if err == model.ErrNotFound {
	// 		logx.Errorf("TopicPurchaseModel.FindOneByUidTid", logx.LogField{Key: "Err", Value: err.Error()})
	// 		return errorx.NewDefaultError("Not found")
	// 	}
	// 	return errorx.NewDefaultError("Not found")
	// }

	// purchase.
	// l.svcCtx.TopicPurchaseModel.Update(l.ctx, )

	return nil
}
