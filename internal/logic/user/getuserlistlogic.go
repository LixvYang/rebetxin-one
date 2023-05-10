package user

import (
	"context"

	"github.com/lixvyang/rebetxin-one/internal/svc"
	"github.com/lixvyang/rebetxin-one/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserListLogic {
	return &GetUserListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserListLogic) GetUserList() (resp []types.User, err error) {
	// todo: add your logic here and delete this line

	return
}
