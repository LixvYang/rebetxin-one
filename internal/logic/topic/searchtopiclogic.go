package topic

import (
	"context"
	"fmt"
	"time"

	"github.com/lixvyang/rebetxin-one/internal/svc"
	"github.com/lixvyang/rebetxin-one/internal/types"
	"github.com/lixvyang/rebetxin-one/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchTopicLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSearchTopicLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchTopicLogic {
	return &SearchTopicLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SearchTopicLogic) SearchTopic(req *types.SearchTopicReq) (resp []types.GetTopicDataResp, err error) {
	topicList, err := l.svcCtx.TopicModel.Search(l.ctx, req.Title, req.Intro, req.Content)
	if err != nil {
		logx.Errorw("TopicModel.Search", logx.LogField{Value: err.Error(), Key: "Err"})
		return nil, err
	}
	topicDataList := l.getTopicDataList(topicList)

	return topicDataList, nil
}

func (l *SearchTopicLogic) getTopicDataList(args []*model.Topic) []types.GetTopicDataResp {
	uid := fmt.Sprintf("%s", l.ctx.Value("uid"))
	isCollectMap := make(map[string]bool)
	ok := false

	topicDataList := make([]types.GetTopicDataResp, len(args))
	if uid != "" {
		if time.Since(l.svcCtx.TopicCollectMap.QueryTime) > time.Minute*10 {
			resp, err := l.svcCtx.TopicCollectModel.ListByUid(l.ctx, uid)
			if err != nil {
				logx.Errorw("TopicCollectRPC.GetTopicCollectByUid", logx.LogField{Key: "Err: ", Value: err.Error()})
				return nil
			}

			m := make(map[string]bool)

			for _, tc := range resp {
				if tc.Status == 1 {
					m[tc.Tid] = true
				}
			}

			l.svcCtx.TopicCollectMap.TopicCollectMap[uid] = m
		}

		isCollectMap, ok = l.svcCtx.TopicCollectMap.TopicCollectMap["uid"]
		if !ok {
			resp, err := l.svcCtx.TopicCollectModel.ListByUid(l.ctx, uid)
			if err != nil {
				logx.Errorw("TopicCollectRPC.GetTopicCollectByUid", logx.LogField{Key: "Err: ", Value: err.Error()})
				return nil
			}

			m := make(map[string]bool)

			for _, tc := range resp {
				if tc.Status == 1 {
					m[tc.Tid] = true
				}
			}

			l.svcCtx.TopicCollectMap.TopicCollectMap["uid"] = m
		}
	}

	for i := 0; i < len(args); i++ {
		if uid != "" {
			if isCollectMap[args[i].Tid] {
				topicDataList[i].IsCollect = 1
			}
		}
		topicDataList[i].Category = (*types.Category)(l.svcCtx.CategoryMap[args[i].Cid])
		topicDataList[i].Cid = args[i].Cid
		topicDataList[i].CollectCount = args[i].CollectCount
		topicDataList[i].Content = args[i].Content
		topicDataList[i].CreatedAt = args[i].CreatedAt.String()
		topicDataList[i].DeletedAt = args[i].DeletedAt.Time.String()
		topicDataList[i].EndTime = args[i].EndTime.Time.String()
		topicDataList[i].Id = args[i].Id
		topicDataList[i].ImgUrl = args[i].ImgUrl
		topicDataList[i].Intro = args[i].Intro
		topicDataList[i].IsStop = args[i].IsStop
		topicDataList[i].NoPrice = args[i].NoPrice.String()
		topicDataList[i].NoRatio = args[i].NoRatio.String()
		topicDataList[i].ReadCount = args[i].ReadCount
		topicDataList[i].RefundEndTime = args[i].RefundEndTime.Time.String()
		topicDataList[i].Tid = args[i].Tid
		topicDataList[i].Title = args[i].Title
		topicDataList[i].TotalPrice = args[i].TotalPrice.String()
		topicDataList[i].UpdatedAt = args[i].UpdatedAt.String()
		topicDataList[i].YesPrice = args[i].YesPrice.String()
		topicDataList[i].YesRatio = args[i].YesRatio.String()
	}

	return topicDataList
}
