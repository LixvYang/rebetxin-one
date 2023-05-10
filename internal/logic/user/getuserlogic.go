package user

import (
	"context"
	"fmt"

	"github.com/lixvyang/rebetxin-one/common/errorx"
	"github.com/lixvyang/rebetxin-one/internal/svc"
	"github.com/lixvyang/rebetxin-one/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserLogic {
	return &GetUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserLogic) GetUser() (resp *types.User, err error) {
	uid := fmt.Sprintf("%s", l.ctx.Value("uid"))
	fmt.Println(uid)
	user, err := l.svcCtx.UserModel.FindOneByUid(l.ctx, uid)
	if err != nil {
		logx.Errorw("UserModel.FindOneByUid err: ", logx.LogField{Key: "err", Value: err.Error()})
		return nil, errorx.NewDefaultError("Error")
	}

	respData := new(types.User)
	respData.AvatarUrl = user.AvatarUrl.String
	respData.Biography = user.Biography.String
	respData.FullName = user.FullName.String
	respData.IdentityNumber = user.IdentityNumber
	respData.Uid = user.Uid
	respData.CreatedAt = user.CreatedAt.String()
	respData.UpdatedAt = user.UpdatedAt.String()

	return respData, nil
}
