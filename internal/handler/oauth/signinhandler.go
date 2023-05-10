package oauth

import (
	"github.com/lixvyang/rebetxin-one/common/errorx"
	"net/http"

	"github.com/lixvyang/rebetxin-one/internal/logic/oauth"
	"github.com/lixvyang/rebetxin-one/internal/svc"
	"github.com/lixvyang/rebetxin-one/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func SigninHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SigninReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, errorx.NewDefaultParamsFailedError())
			return
		}

		l := oauth.NewSigninLogic(r.Context(), svcCtx)
		resp, err := l.Signin(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, errorx.NewSuccessJson(resp))
		}
	}
}
