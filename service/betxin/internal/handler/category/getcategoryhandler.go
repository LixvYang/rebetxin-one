package category

import (
	"net/http"

	"github.com/lixvyang/rebetxin-one/common/errorx"

	"github.com/lixvyang/rebetxin-one/service/betxin/internal/logic/category"
	"github.com/lixvyang/rebetxin-one/service/betxin/internal/svc"
	"github.com/lixvyang/rebetxin-one/service/betxin/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetCategoryHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetCategoryReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, errorx.NewDefaultParamsFailedError())
			return
		}

		l := category.NewGetCategoryLogic(r.Context(), svcCtx)
		resp, err := l.GetCategory(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, errorx.NewSuccessJson(resp))
		}
	}
}
