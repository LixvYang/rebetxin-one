package logic

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/lixvyang/rebetxin-one/common/mixinpkg"
	"github.com/lixvyang/rebetxin-one/service/job/internal/svc"
)

type RefundTaskHandler struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRefundTask(svcCtx *svc.ServiceContext) *RefundTaskHandler {
	c := context.Background()
	return &RefundTaskHandler{
		ctx:    c,
		svcCtx: svcCtx,
	}
}

// every one minute exec : if return err != nil , asynq will retry
func (l *RefundTaskHandler) ProcessTask(ctx context.Context, t *asynq.Task) error {
	fmt.Printf("shcedule RefundTaskHandler job demo -----> \n")
	var p mixinpkg.RefundTask
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	// Email delivery code ...

	p.Exec(l.ctx)

	return nil
}
