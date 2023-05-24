package purchase

import (
	"net/http"

	"github.com/lixvyang/rebetxin-one/common/errorx"
	"github.com/lixvyang/rebetxin-one/internal/logic/purchase"
	"github.com/lixvyang/rebetxin-one/internal/svc"
	"github.com/lixvyang/rebetxin-one/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetTopicPurchaseHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetTopicPurchaseReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := purchase.NewGetTopicPurchaseLogic(r.Context(), svcCtx)
		resp, err := l.GetTopicPurchase(&req)
		if err != nil {
			httpx.OkJsonCtx(r.Context(), w, errorx.NewDefaultError("Error"))
		} else {
			httpx.OkJsonCtx(r.Context(), w, errorx.NewSuccessJson(resp))
		}
	}
}
