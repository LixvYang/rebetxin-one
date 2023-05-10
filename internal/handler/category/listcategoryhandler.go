package category

import (
	"github.com/lixvyang/rebetxin-one/common/errorx"
	"net/http"

	"github.com/lixvyang/rebetxin-one/internal/logic/category"
	"github.com/lixvyang/rebetxin-one/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func ListCategoryHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := category.NewListCategoryLogic(r.Context(), svcCtx)
		resp, err := l.ListCategory()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, errorx.NewSuccessJson(resp))
		}
	}
}
