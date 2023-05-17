package snapshot

import (
	"net/http"

	"github.com/lixvyang/rebetxin-one/common/errorx"
	"github.com/lixvyang/rebetxin-one/internal/logic/snapshot"
	"github.com/lixvyang/rebetxin-one/internal/svc"
	"github.com/lixvyang/rebetxin-one/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetSnapshotHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetSnapshotReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, errorx.NewDefaultParamsFailedError())
			return
		}

		l := snapshot.NewGetSnapshotLogic(r.Context(), svcCtx)
		resp, err := l.GetSnapshot(&req)
		if err != nil {
			httpx.Error(w, errorx.NewCategoryError("Error"))
		} else {
			httpx.OkJsonCtx(r.Context(), w, errorx.NewSuccessJson(resp))
		}
	}
}
