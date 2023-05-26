package collect

import (
	"net/http"

	"github.com/lixvyang/rebetxin-one/common/errorx"

	"github.com/lixvyang/rebetxin-one/service/betxin/internal/logic/collect"
	"github.com/lixvyang/rebetxin-one/service/betxin/internal/svc"
	"github.com/lixvyang/rebetxin-one/service/betxin/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func CreateTopicCollectHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CreateTopicCollectReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, errorx.NewDefaultParamsFailedError())
			return
		}

		l := collect.NewCreateTopicCollectLogic(r.Context(), svcCtx)
		err := l.CreateTopicCollect(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJson(w, errorx.NewSuccessJson("Success"))
		}
	}
}
