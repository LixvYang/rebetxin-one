package collect

import (
	"net/http"

	"github.com/lixvyang/rebetxin-one/common/errorx"

	"github.com/lixvyang/rebetxin-one/service/betxin/internal/logic/collect"
	"github.com/lixvyang/rebetxin-one/service/betxin/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func ListTopicCollectHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := collect.NewListTopicCollectLogic(r.Context(), svcCtx)
		resp, err := l.ListTopicCollect()
		if err != nil {
			httpx.Error(w, errorx.NewDefaultError("Err"))
		} else {
			httpx.OkJsonCtx(r.Context(), w, errorx.NewSuccessJson(resp))
		}
	}
}
