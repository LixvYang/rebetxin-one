package category

import (
	"net/http"

	"github.com/lixvyang/rebetxin-one/common/errorx"

	"github.com/lixvyang/rebetxin-one/service/betxin/internal/logic/category"
	"github.com/lixvyang/rebetxin-one/service/betxin/internal/svc"
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
