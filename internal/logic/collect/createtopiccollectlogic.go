package collect

import (
	"context"
	"fmt"

	"github.com/lixvyang/rebetxin-one/common/errorx"
	"github.com/lixvyang/rebetxin-one/internal/svc"
	"github.com/lixvyang/rebetxin-one/internal/types"
	"github.com/lixvyang/rebetxin-one/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateTopicCollectLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateTopicCollectLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateTopicCollectLogic {
	return &CreateTopicCollectLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateTopicCollectLogic) CreateTopicCollect(req *types.CreateTopicCollectReq) (err error) {
	uid := fmt.Sprintf("%s", l.ctx.Value("uid"))
	_, err = l.svcCtx.TopicCollectModel.Insert(l.ctx, &model.TopicCollect{
		Uid:    uid,
		Tid:    req.Tid,
		Status: 1,
	})
	if err != nil {
		logx.Errorw("TopicCollectModel.Insert", logx.LogField{Key: "Err ", Value: err.Error()})
		return errorx.NewDefaultError("Error")
	}

	return nil
}
