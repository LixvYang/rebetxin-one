package user

import (
	"github.com/lixvyang/rebetxin-one/common/errorx"
	"net/http"

	"github.com/lixvyang/rebetxin-one/internal/logic/user"
	"github.com/lixvyang/rebetxin-one/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetUserHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := user.NewGetUserLogic(r.Context(), svcCtx)
		resp, err := l.GetUser()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, errorx.NewSuccessJson(resp))
		}
	}
}
