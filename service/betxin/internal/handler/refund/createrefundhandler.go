package refund

import (
	"net/http"

	"github.com/lixvyang/rebetxin-one/common/errorx"

	"github.com/lixvyang/rebetxin-one/service/betxin/internal/logic/refund"
	"github.com/lixvyang/rebetxin-one/service/betxin/internal/svc"
	"github.com/lixvyang/rebetxin-one/service/betxin/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func CreateRefundHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CreateRefundReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, errorx.NewDefaultParamsFailedError())
			return
		}

		l := refund.NewCreateRefundLogic(r.Context(), svcCtx)
		err := l.CreateRefund(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, errorx.NewDefaultError("Error"))
		} else {
			httpx.OkJson(w, errorx.NewSuccessJson("Success"))
		}
	}
}
