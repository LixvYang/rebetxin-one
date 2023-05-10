package collect

import (
	"github.com/lixvyang/rebetxin-one/common/errorx"
	"net/http"

	"github.com/lixvyang/rebetxin-one/internal/logic/collect"
	"github.com/lixvyang/rebetxin-one/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func ListTopicCollectHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := collect.NewListTopicCollectLogic(r.Context(), svcCtx)
		resp, err := l.ListTopicCollect()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, errorx.NewSuccessJson(resp))
		}
	}
}
