package topicAdmin

import (
	"context"
	"database/sql"

	"github.com/gofrs/uuid"
	"github.com/lixvyang/rebetxin-one/common/errorx"
	"github.com/lixvyang/rebetxin-one/service/betxin/internal/svc"
	"github.com/lixvyang/rebetxin-one/service/betxin/internal/types"
	"github.com/lixvyang/rebetxin-one/service/betxin/model"
	"github.com/shopspring/decimal"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateTopicLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateTopicLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateTopicLogic {
	return &CreateTopicLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateTopicLogic) CreateTopic(req *types.CreateTopicReq) (resp *types.Topic, err error) {
	topic := l.createTopicPrepare(req)

	_, err = l.svcCtx.TopicModel.Insert(l.ctx, topic)
	if err != nil {
		logx.Errorw("TopicModel.Insert", logx.LogField{Key: "Err: ", Value: err.Error()})
		return nil, errorx.NewCodeError(1003, "Create Topic Error!")
	}

	return nil, err
}

func (l *CreateTopicLogic) createTopicPrepare(req *types.CreateTopicReq) (resp *model.Topic) {
	uuid, _ := uuid.NewV4()
	resp = new(model.Topic)

	resp = &model.Topic{
		Tid:           uuid.String(),
		Cid:           int64(req.Cid),
		Title:         req.Title,
		Intro:         req.Intro,
		Content:       req.Content,
		ImgUrl:        req.ImgUrl,
		RefundEndTime: sql.NullTime{Valid: true, Time: l.svcCtx.StringToTime(req.RefundEndTime)},
		EndTime:       sql.NullTime{Valid: true, Time: l.svcCtx.StringToTime(req.EndTime)},
	}
	resp.YesRatio, _ = decimal.NewFromString("50.00")
	resp.NoRatio, _ = decimal.NewFromString("50.00")
	return
}
