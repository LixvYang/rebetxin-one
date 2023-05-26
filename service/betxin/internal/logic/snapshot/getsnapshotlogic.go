package snapshot

import (
	"context"
	"fmt"

	"github.com/lixvyang/rebetxin-one/common/errorx"
	"github.com/lixvyang/rebetxin-one/service/betxin/internal/svc"
	"github.com/lixvyang/rebetxin-one/service/betxin/internal/types"
	"github.com/lixvyang/rebetxin-one/service/betxin/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSnapshotLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetSnapshotLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSnapshotLogic {
	return &GetSnapshotLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSnapshotLogic) GetSnapshot(req *types.GetSnapshotReq) (resp *types.GetSnapshotResp, err error) {
	fmt.Println(*req)
	snapshot, err := l.svcCtx.SnapshotModel.FindOneByTraceId(l.ctx, req.TraceId)
	if err != nil {
		if err == model.ErrNotFound {
			return &types.GetSnapshotResp{
				End: 0,
			}, errorx.NewDefaultError("Error")
		}
		return &types.GetSnapshotResp{
			End: 0,
		}, errorx.NewDefaultError("Error")
	}
	fmt.Println(222)

	if snapshot.End == 1 {
		fmt.Println(111)
		return &types.GetSnapshotResp{
			End: 1,
		}, nil
	}
	fmt.Println(333)

	return &types.GetSnapshotResp{
		End: 0,
	}, errorx.NewDefaultError("Error")
}
