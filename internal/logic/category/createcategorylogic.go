package category

import (
	"context"

	"github.com/lixvyang/rebetxin-one/common/errorx"
	"github.com/lixvyang/rebetxin-one/internal/svc"
	"github.com/lixvyang/rebetxin-one/internal/types"
	"github.com/lixvyang/rebetxin-one/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateCategoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateCategoryLogic {
	return &CreateCategoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateCategoryLogic) CreateCategory(req *types.CreateCategoryReq) (resp *types.CreateCategoryResp, err error) {
	_, err = l.svcCtx.CategoryModel.Insert(l.ctx, &model.Category{
		CategoryName: req.CategoryName,
	})
	if err != nil {
		logx.Errorw("CategoryModel.Insert", logx.LogField{Key: "Err", Value: err.Error()})
		return nil, errorx.NewCategoryError("Create category error!")
	}

	return &types.CreateCategoryResp{
		Data: "Success",
	}, nil
}
