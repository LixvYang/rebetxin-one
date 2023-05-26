package user

import (
	"context"
	"database/sql"

	"github.com/lixvyang/rebetxin-one/common/errorx"
	"github.com/lixvyang/rebetxin-one/service/betxin/internal/svc"
	"github.com/lixvyang/rebetxin-one/service/betxin/internal/types"
	"github.com/lixvyang/rebetxin-one/service/betxin/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateUserLogic {
	return &CreateUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateUserLogic) CreateUser(req *types.CreateUserReq) (err error) {
	_, err = l.svcCtx.UserModel.Insert(l.ctx, &model.User{
		IdentityNumber: req.IdentityNumber,
		FullName:       sql.NullString{String: req.FullName, Valid: true},
		Uid:            req.Uid,
		AvatarUrl:      sql.NullString{String: req.AvatarUrl, Valid: true},
		Biography:      sql.NullString{String: req.Biography, Valid: true},
	})
	if err != nil {
		logx.Errorw("UserModel.Insert failed! ", logx.LogField{Key: "Err: ", Value: err})
		return errorx.NewDefaultError("CreateUser error!")
	}
	return nil
}
