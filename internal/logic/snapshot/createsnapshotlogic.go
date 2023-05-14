package snapshot

import (
	"context"
	"fmt"

	"github.com/lixvyang/rebetxin-one/internal/svc"
	"github.com/lixvyang/rebetxin-one/internal/types"
	"github.com/lixvyang/rebetxin-one/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateSnapshotLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateSnapshotLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateSnapshotLogic {
	return &CreateSnapshotLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateSnapshotLogic) CreateSnapshot(req *types.CreateSnapshotReq) error {
	fmt.Println(*req)

	_, err := l.svcCtx.SnapshotModel.Insert(l.ctx, &model.Snapshot{
		TraceId: req.TraceId,
		Uid: req.Uid,
		Tid: req.Tid,
	})	
	if err != nil {
		logx.Errorw("SnapshotModel.Insert", logx.LogField{Key: "Err", Value: err.Error()})
	}
	

	return nil
}
