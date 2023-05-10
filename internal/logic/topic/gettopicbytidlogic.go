package topic

import (
	"context"
	"fmt"

	"github.com/lixvyang/rebetxin-one/common/errorx"
	"github.com/lixvyang/rebetxin-one/internal/svc"
	"github.com/lixvyang/rebetxin-one/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetTopicByTidLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetTopicByTidLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTopicByTidLogic {
	return &GetTopicByTidLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetTopicByTidLogic) GetTopicByTid(req *types.GetTopicByTidReq) (resp *types.GetTopicDataResp, err error) {
	topic, err := l.svcCtx.TopicModel.FindOneByTid(l.ctx, req.Tid)
	if err != nil {
		logx.Errorw("TopicModel.FindOneByTid", logx.LogField{Key: "Err: ", Value: err.Error()})
		return nil, errorx.NewDefaultError("FindOneByTid error.")
	}

	respData := new(types.GetTopicDataResp)
	respData.Cid = topic.Cid
	respData.CollectCount = topic.CollectCount
	respData.Content = topic.Content
	respData.CreatedAt = l.svcCtx.TimeToString(topic.CreatedAt)
	respData.DeletedAt = l.svcCtx.TimeToString(topic.DeletedAt.Time)
	respData.EndTime = l.svcCtx.TimeToString(topic.EndTime.Time)
	respData.Id = topic.Id
	respData.ImgUrl = topic.ImgUrl
	respData.Intro = topic.Intro
	respData.IsStop = topic.IsStop
	respData.NoPrice = topic.NoPrice.String()
	respData.NoRatio = topic.NoRatio.String()
	respData.ReadCount = topic.ReadCount
	respData.RefundEndTime = topic.RefundEndTime.Time.String()
	respData.Tid = topic.Tid
	respData.Title = topic.Title
	respData.TotalPrice = topic.TotalPrice.String()
	respData.UpdatedAt = topic.UpdatedAt.String()
	respData.YesPrice = topic.YesPrice.String()
	respData.YesRatio = topic.YesRatio.String()
	respData.Category = (*types.Category)(l.svcCtx.CategoryMap[topic.Cid])

	uid := fmt.Sprintf("%s", l.ctx.Value("uid"))
	
	if uid != "" {
		resp, err := l.svcCtx.TopicCollectModel.ListByUid(l.ctx, uid)
		if err != nil {
			logx.Errorw("TopicCollectRPC.GetTopicCollectByUid", logx.LogField{Key: "Err", Value: err.Error()})
		}
		for _, r := range resp {
			if r.Tid == topic.Tid {
				respData.IsCollect = 1
			}
		}
	}

	return respData, nil
}
