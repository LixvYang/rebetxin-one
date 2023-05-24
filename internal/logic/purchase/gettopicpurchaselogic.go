package purchase

import (
	"context"

	"github.com/lixvyang/rebetxin-one/common/errorx"
	"github.com/lixvyang/rebetxin-one/internal/svc"
	"github.com/lixvyang/rebetxin-one/internal/types"
	"github.com/lixvyang/rebetxin-one/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetTopicPurchaseLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetTopicPurchaseLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTopicPurchaseLogic {
	return &GetTopicPurchaseLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetTopicPurchaseLogic) GetTopicPurchase(req *types.GetTopicPurchaseReq) (resp *types.Topicpurchase, err error) {
	resp = new(types.Topicpurchase)
	tp, err := l.svcCtx.TopicPurchaseModel.FindOneByUidTid(l.ctx, req.Uid, req.Tid)
	if err != nil {
		if err == model.ErrNotFound {
			return resp, errorx.NewDefaultError("Not Found")
		}
		return resp, errorx.NewDefaultError(err.Error())
	}
	resp.Id = tp.Id
	resp.CreatedAt = tp.CreatedAt.String()
	resp.UpdatedAt = tp.UpdatedAt.String()
	resp.Uid = tp.Uid
	resp.Tid = tp.Tid
	resp.YesPrice = tp.YesPrice.String()
	resp.NoPrice = tp.NoPrice.String()

	return resp, nil
}
