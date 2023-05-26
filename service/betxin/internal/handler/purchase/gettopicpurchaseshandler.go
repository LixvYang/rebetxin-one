package purchase

import (
	"net/http"

	"github.com/lixvyang/rebetxin-one/service/betxin/internal/logic/purchase"
	"github.com/lixvyang/rebetxin-one/service/betxin/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetTopicPurchasesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := purchase.NewGetTopicPurchasesLogic(r.Context(), svcCtx)
		resp, err := l.GetTopicPurchases()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
