package purchase

import (
	"context"
	"fmt"

	"github.com/lixvyang/rebetxin-one/service/betxin/internal/svc"
	"github.com/lixvyang/rebetxin-one/service/betxin/internal/types"
	"github.com/lixvyang/rebetxin-one/service/betxin/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListTopicPurchaseLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListTopicPurchaseLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListTopicPurchaseLogic {
	return &ListTopicPurchaseLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListTopicPurchaseLogic) ListTopicPurchase() (resp []types.GetTopicDataResp, err error) {
	uid := fmt.Sprintf("%s", l.ctx.Value("uid"))
	tpList, err := l.svcCtx.TopicPurchaseModel.ListByUid(l.ctx, uid)
	topicList := make([]*model.Topic, 0)
	for _, tp := range tpList {
		topic, _ := l.svcCtx.TopicModel.FindOneByTid(l.ctx, tp.Tid)
		topicList = append(topicList, topic)
	}
	respData := l.getTopicDataList(topicList, uid)

	return respData, nil
}

func (l *ListTopicPurchaseLogic) getTopicDataList(args []*model.Topic, uid string) []types.GetTopicDataResp {
	isCollectMap := make(map[string]bool)
	ok := false

	topicDataList := make([]types.GetTopicDataResp, len(args))
	if uid != "" {
		// if time.Since(l.svcCtx.TopicCollectMap.QueryTime) > time.Minute*10 {
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

		isCollectMap, ok = l.svcCtx.TopicCollectMap.TopicCollectMap[uid]
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
		topicDataList[i].Category = new(types.Category)
		topicDataList[i].Category.Id = l.svcCtx.CategoryMap[args[i].Cid].Id
		topicDataList[i].Category.CategoryName = l.svcCtx.CategoryMap[args[i].Cid].CategoryName
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
