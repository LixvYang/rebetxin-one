package collect

import (
	"context"
	"fmt"

	"github.com/lixvyang/rebetxin-one/internal/svc"
	"github.com/lixvyang/rebetxin-one/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListTopicCollectLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListTopicCollectLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListTopicCollectLogic {
	return &ListTopicCollectLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListTopicCollectLogic) ListTopicCollect() (resp []types.Topiccollect, err error) {
	uid := fmt.Sprintf("%s", l.ctx.Value("uid"))

	topicList, err := l.svcCtx.TopicCollectModel.ListByUid(l.ctx, uid)
	if err != nil {
		logx.Errorw("TopicCollectModel.ListByUid", logx.LogField{Key: "Err", Value: err.Error()})
		return nil, err
	}

	respData := make([]types.Topiccollect, len(topicList))
	for i := 0; i < len(respData); i++ {
		respData[i].Id = topicList[i].Id
		respData[i].CreatedAt = l.svcCtx.TimeToString(topicList[i].CreatedAt)
		respData[i].Status = topicList[i].Status
		respData[i].Tid = topicList[i].Tid
		respData[i].Uid = topicList[i].Uid
		respData[i].UpdatedAt = l.svcCtx.TimeToString(topicList[i].UpdatedAt)
	}

	return respData, nil
}
