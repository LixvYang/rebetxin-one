package topicAdmin

import (
	"context"
	"database/sql"

	"github.com/lixvyang/rebetxin-one/common/convert"
	"github.com/lixvyang/rebetxin-one/common/errorx"
	"github.com/lixvyang/rebetxin-one/service/betxin/internal/svc"
	"github.com/lixvyang/rebetxin-one/service/betxin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateTopicLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateTopicLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateTopicLogic {
	return &UpdateTopicLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateTopicLogic) UpdateTopic(req *types.UpdateTopicReq) error {
	topic, err := l.svcCtx.TopicModel.FindOneByTid(l.ctx, req.Tid)
	if err != nil {
		logx.Errorw("TopicModel.FindOneByTid", logx.LogField{Key: "Err", Value: err.Error()})
		return errorx.NewDefaultError("Err TopicModel.FindOneByTid")
	}

	topic.Cid = req.Cid
	topic.Title = req.Title
	topic.Content = req.Content
	topic.ImgUrl = req.ImgUrl
	topic.EndTime = sql.NullTime{Valid: true, Time: convert.StringToTime(req.EndTime)}
	topic.Intro = req.Intro
	topic.IsStop = req.IsStop
	topic.RefundEndTime = sql.NullTime{Valid: true, Time: convert.StringToTime(req.RefundEndTime)}

	err = l.svcCtx.TopicModel.Update(l.ctx, topic)
	if err != nil {
		logx.Errorw("TopicModel.Stop", logx.LogField{Key: "Err", Value: err.Error()})
		return errorx.NewDefaultError("Err TopicModel.FindOneByTid")
	}

	return nil
}
