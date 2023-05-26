package topicAdmin

import (
	"net/http"

	"github.com/lixvyang/rebetxin-one/common/errorx"

	"github.com/lixvyang/rebetxin-one/service/betxin/internal/logic/topicAdmin"
	"github.com/lixvyang/rebetxin-one/service/betxin/internal/svc"
	"github.com/lixvyang/rebetxin-one/service/betxin/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func UpdateTopicHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UpdateTopicReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, errorx.NewDefaultParamsFailedError())
			return
		}

		l := topicAdmin.NewUpdateTopicLogic(r.Context(), svcCtx)
		err := l.UpdateTopic(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.Ok(w)
		}
	}
}
