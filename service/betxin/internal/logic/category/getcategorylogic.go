package category

import (
	"context"

	"github.com/lixvyang/rebetxin-one/common/errorx"
	"github.com/lixvyang/rebetxin-one/service/betxin/internal/svc"
	"github.com/lixvyang/rebetxin-one/service/betxin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCategoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCategoryLogic {
	return &GetCategoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCategoryLogic) GetCategory(req *types.GetCategoryReq) (resp *types.Category, err error) {
	category, err := l.svcCtx.CategoryModel.FindOne(l.ctx, req.Id)
	if err != nil {
		logx.Errorw("error: ", logx.LogField{Key: "CategoryModel.FindOne error", Value: err})
		return nil, errorx.NewCodeError(1003, "Find Category Error!")
	}

	return (*types.Category)(category), nil
}
