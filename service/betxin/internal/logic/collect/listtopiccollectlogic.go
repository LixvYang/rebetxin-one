package collect

import (
	"context"
	"fmt"

	"github.com/lixvyang/rebetxin-one/common/convert"
	"github.com/lixvyang/rebetxin-one/service/betxin/internal/svc"
	"github.com/lixvyang/rebetxin-one/service/betxin/internal/types"

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

func (l *ListTopicCollectLogic) ListTopicCollect() ([]types.Topiccollect, error) {
	uid := fmt.Sprintf("%s", l.ctx.Value("uid"))
	topicList, err := l.svcCtx.TopicCollectModel.ListByUid(l.ctx, uid)
	if err != nil {
		logx.Errorw("TopicCollectModel.ListByUid", logx.LogField{Key: "Err", Value: err.Error()})
		return nil, err
	}

	respData := make([]types.Topiccollect, len(topicList))
	for i := 0; i < len(respData); i++ {
		respData[i].Topic = new(types.GetTopicDataResp)
		topic, _ := l.svcCtx.TopicModel.FindOneByTid(l.ctx, topicList[i].Tid)
		respData[i].Id = topicList[i].Id
		respData[i].CreatedAt = convert.TimeToString(topicList[i].CreatedAt)
		respData[i].Status = topicList[i].Status
		respData[i].Tid = topicList[i].Tid
		respData[i].Uid = topicList[i].Uid
		respData[i].UpdatedAt = convert.TimeToString(topicList[i].UpdatedAt)
		respData[i].Topic.Cid = topic.Cid
		respData[i].Topic.CollectCount = topic.CollectCount
		respData[i].Topic.Content = topic.Content
		respData[i].Topic.CreatedAt = convert.TimeToString(topic.CreatedAt)
		respData[i].Topic.DeletedAt = convert.TimeToString(topic.DeletedAt.Time)
		respData[i].Topic.EndTime = convert.TimeToString(topic.EndTime.Time)
		respData[i].Topic.Id = topic.Id
		respData[i].Topic.ImgUrl = topic.ImgUrl
		respData[i].Topic.Intro = topic.Intro
		respData[i].Topic.IsStop = topic.IsStop
		respData[i].Topic.NoPrice = topic.NoPrice.String()
		respData[i].Topic.NoRatio = topic.NoRatio.String()
		respData[i].Topic.ReadCount = topic.ReadCount
		respData[i].Topic.RefundEndTime = topic.RefundEndTime.Time.String()
		respData[i].Topic.Tid = topic.Tid
		respData[i].Topic.Title = topic.Title
		respData[i].Topic.TotalPrice = topic.TotalPrice.String()
		respData[i].Topic.UpdatedAt = topic.UpdatedAt.String()
		respData[i].Topic.YesPrice = topic.YesPrice.String()
		respData[i].Topic.YesRatio = topic.YesRatio.String()

		respData[i].Topic.IsCollect = 1
		respData[i].Topic.Category = new(types.Category)
		respData[i].Topic.Category.Id = l.svcCtx.CategoryMap[topic.Cid].Id
		respData[i].Topic.Category.CategoryName = l.svcCtx.CategoryMap[topic.Cid].CategoryName
	}

	return respData, nil
}
