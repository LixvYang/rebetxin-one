package purchase

import (
	"net/http"

	"github.com/lixvyang/rebetxin-one/common/errorx"
	"github.com/lixvyang/rebetxin-one/service/betxin/internal/logic/purchase"
	"github.com/lixvyang/rebetxin-one/service/betxin/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func ListTopicPurchaseHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := purchase.NewListTopicPurchaseLogic(r.Context(), svcCtx)
		resp, err := l.ListTopicPurchase()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, errorx.NewDefaultError("Error"))
		} else {
			httpx.OkJsonCtx(r.Context(), w, errorx.NewSuccessJson(resp))
		}
	}
}
