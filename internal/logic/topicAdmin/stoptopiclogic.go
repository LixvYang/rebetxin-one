package topicAdmin

import (
	"context"

	"github.com/lixvyang/rebetxin-one/common/errorx"
	"github.com/lixvyang/rebetxin-one/internal/svc"
	"github.com/lixvyang/rebetxin-one/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type StopTopicLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewStopTopicLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StopTopicLogic {
	return &StopTopicLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *StopTopicLogic) StopTopic(req *types.StopTopicReq) error {
	topic, err := l.svcCtx.TopicModel.FindOneByTid(l.ctx, req.Tid)
	if err != nil {
		logx.Errorw("TopicModel.FindOneByTid", logx.LogField{Key: "Err", Value: err.Error()})
		return errorx.NewCodeError(1003, "Error")
	}
	topic.IsStop = 1

	err = l.svcCtx.TopicModel.Update(l.ctx, topic)
	if err != nil {
		logx.Errorw("TopicModel.Stop", logx.LogField{Key: "Err", Value: err.Error()})
		return errorx.NewCodeError(1003, "Error")
	}

	return nil
}
