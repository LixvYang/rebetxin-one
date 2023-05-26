package topicAdmin

import (
	"context"

	"github.com/lixvyang/rebetxin-one/common/errorx"
	"github.com/lixvyang/rebetxin-one/service/betxin/internal/svc"
	"github.com/lixvyang/rebetxin-one/service/betxin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteTopicLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteTopicLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteTopicLogic {
	return &DeleteTopicLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteTopicLogic) DeleteTopic(req *types.DeleteTopicReq) error {
	topic, err := l.svcCtx.TopicModel.FindOneByTid(l.ctx, req.Tid)
	if err != nil {
		logx.Errorw("TopicModel.FindOneByTid", logx.LogField{Key: "Err", Value: err.Error()})
		return errorx.NewCodeError(1003, "DeleteTopic Error!")
	}

	err = l.svcCtx.TopicModel.Delete(l.ctx, topic.Id)
	if err != nil {
		logx.Errorw("TopicModel.Delete", logx.LogField{Key: "Err", Value: err.Error()})
		return errorx.NewCodeError(1003, "DeleteTopic Error!")

	}

	return nil
}
