package topic

import (
	"context"
	"fmt"
	"time"

	"github.com/lixvyang/rebetxin-one/common/token"
	"github.com/lixvyang/rebetxin-one/internal/svc"
	"github.com/lixvyang/rebetxin-one/internal/types"
	"github.com/lixvyang/rebetxin-one/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListTopicByCidLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListTopicByCidLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListTopicByCidLogic {
	return &ListTopicByCidLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

var (
	defaultPageSize int64 = 10
	defaultCursor   int64 = 10
)

func (l *ListTopicByCidLogic) ListTopicByCid(req *types.ListTopicByCidReq) (resp *types.ListTopicByCidResp, err error) {
	defaultCursor = model.DefaultCursor
	page := token.Token(req.PageToken).Decode()
	var (
		cursor   int64 = defaultCursor
		pageSize int64 = defaultPageSize
	)

	if req.PageToken != "" {
		// 解析分页
		if page.NextTimeAtUTC > time.Now().Unix() || time.Now().Unix()-page.NextTimeAtUTC > int64(time.Hour)*24 {
			logx.Errorw("bad page token", logx.LogField{Key: "Page Token", Value: "bad page token"})
			return nil, err
		}

		// invaild
		if page.PreID <= 0 || page.NextTimeAtUTC == 0 || page.NextTimeAtUTC > time.Now().Unix() || page.PageSize <= 0 {
			logx.Errorw("bad page token", logx.LogField{Key: "Page Token", Value: "bad page token"})
			return
		}
		cursor = page.PreID
		pageSize = page.PageSize
	}

	topicList, err := l.svcCtx.TopicModel.ListByCid(l.ctx, req.Cid, cursor, pageSize+1)
	if err != nil {
		logx.Errorw("TopicModel.ListByCid", logx.LogField{Key: "err", Value: err.Error()})
		return nil, err
	}

	topicDataList := l.getTopicDataList(topicList)

	var (
		hasPrePage   bool
		prePageToken string
	)

	if len(topicDataList) > int(pageSize) {
		hasPrePage = true
	}

	// if has pre page
	if hasPrePage {
		prePageInfo := token.Page{
			PreID:         topicDataList[len(topicDataList)-1].Id,
			NextTimeAtUTC: time.Now().Unix(),
			PageSize:      pageSize,
		}
		prePageToken = string(prePageInfo.Encode())

		return &types.ListTopicByCidResp{
			PrePageToken: prePageToken,
			List:         topicDataList[:len(topicDataList)-1],
		}, nil
	}

	return &types.ListTopicByCidResp{
		PrePageToken: prePageToken,
		List:         topicDataList,
	}, nil

	return
}

func (l *ListTopicByCidLogic) getTopicDataList(args []*model.Topic) []types.GetTopicDataResp {
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
