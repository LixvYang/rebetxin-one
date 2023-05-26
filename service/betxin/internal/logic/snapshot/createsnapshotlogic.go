package snapshot

import (
	"context"

	"github.com/lixvyang/rebetxin-one/service/betxin/internal/svc"
	"github.com/lixvyang/rebetxin-one/service/betxin/internal/types"
	"github.com/lixvyang/rebetxin-one/service/betxin/model"

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
	_, err := l.svcCtx.SnapshotModel.FindOneByTraceId(l.ctx, req.TraceId)
	if err != nil {
		if err == model.ErrNotFound {
			_, err := l.svcCtx.SnapshotModel.Insert(l.ctx, &model.Snapshot{
				TraceId: req.TraceId,
				Uid:     req.Uid,
				Tid:     req.Tid,
			})
			if err != nil {
				logx.Errorw("SnapshotModel.Insert", logx.LogField{Key: "Err", Value: err.Error()})
				return err
			}
			return nil
		}
		return err
	}

	return nil
}
