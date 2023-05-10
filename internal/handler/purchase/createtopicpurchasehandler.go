package purchase

import (
	"net/http"

	"github.com/lixvyang/rebetxin-one/internal/logic/purchase"
	"github.com/lixvyang/rebetxin-one/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func CreateTopicPurchaseHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := purchase.NewCreateTopicPurchaseLogic(r.Context(), svcCtx)
		err := l.CreateTopicPurchase()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.Ok(w)
		}
	}
}
