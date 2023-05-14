package logic

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/lixvyang/rebetxin-one/common/mixinpkg"
	"github.com/lixvyang/rebetxin-one/service/job/internal/svc"
)

type CreateTopicTaskHandler struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTopicTask(svcCtx *svc.ServiceContext) *CreateTopicTaskHandler {
	c := context.Background()
	return &CreateTopicTaskHandler{
		ctx:    c,
		svcCtx: svcCtx,
	}
}

// every one minute exec : if return err != nil , asynq will retry
func (l *CreateTopicTaskHandler) ProcessTask(ctx context.Context, t *asynq.Task) error {
	fmt.Printf("shcedule job demo -----> every one second exec \n")
	var p mixinpkg.NewTopicTask
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	// Email delivery code ...

	p.Exec(l.ctx)

	return nil
}
