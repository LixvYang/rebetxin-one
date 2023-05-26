package collect

import (
	"context"
	"fmt"

	"github.com/lixvyang/rebetxin-one/common/errorx"
	"github.com/lixvyang/rebetxin-one/service/betxin/internal/svc"
	"github.com/lixvyang/rebetxin-one/service/betxin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteTopicCollectLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteTopicCollectLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteTopicCollectLogic {
	return &DeleteTopicCollectLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteTopicCollectLogic) DeleteTopicCollect(req *types.DeleteTopicCollectReq) error {
	uid := fmt.Sprintf("%s", l.ctx.Value("uid"))
	tc, err := l.svcCtx.TopicCollectModel.FindOneByUidTid(l.ctx, uid, req.Tid)
	if err != nil {
		logx.Errorw("TopicCollectModel.FindOneByUidTid", logx.LogField{Key: "Err", Value: err.Error()})
		return errorx.NewDefaultError("Error")
	}

	err = l.svcCtx.TopicCollectModel.Delete(l.ctx, tc.Id)
	if err != nil {
		logx.Errorw("TopicCollectModel.Delete", logx.LogField{Key: "Err", Value: err.Error()})
		return errorx.NewDefaultError("Error")
	}
	return nil
}
