package snapshot

import (
	"net/http"

	"github.com/lixvyang/rebetxin-one/common/errorx"
	"github.com/lixvyang/rebetxin-one/service/betxin/internal/logic/snapshot"
	"github.com/lixvyang/rebetxin-one/service/betxin/internal/svc"
	"github.com/lixvyang/rebetxin-one/service/betxin/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func CreateSnapshotHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CreateSnapshotReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, errorx.NewDefaultParamsFailedError())
			return
		}

		l := snapshot.NewCreateSnapshotLogic(r.Context(), svcCtx)
		err := l.CreateSnapshot(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, errorx.NewDefaultError("Error"))
		} else {
			httpx.OkJsonCtx(r.Context(), w, errorx.NewSuccessJson("Success"))
		}
	}
}
