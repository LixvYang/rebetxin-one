package topic

import (
	"context"

	"github.com/lixvyang/rebetxin-one/common/convert"
	"github.com/lixvyang/rebetxin-one/common/errorx"
	"github.com/lixvyang/rebetxin-one/service/betxin/internal/svc"
	"github.com/lixvyang/rebetxin-one/service/betxin/internal/types"

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
	respData.CreatedAt = convert.TimeToString(topic.CreatedAt)
	respData.DeletedAt = convert.TimeToString(topic.DeletedAt.Time)
	respData.EndTime = convert.TimeToString(topic.EndTime.Time)
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

	uid := req.Uid
	if uid != "undefined" {
		resp, err := l.svcCtx.TopicCollectModel.FindOneByUidTid(l.ctx, uid, req.Tid)
		if err != nil {
			logx.Errorw("TopicCollectRPC.GetTopicCollectByUid", logx.LogField{Key: "Err", Value: err.Error()})
		} else {
			if resp.Status == 1 {
				respData.IsCollect = 1
			}
		}
	}

	return respData, nil
}
