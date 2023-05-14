package mixinpkg

import (
	"context"
	"encoding/json"
	"time"

	"github.com/hibiken/asynq"
	"github.com/lixvyang/rebetxin-one/common/convert"
	"github.com/lixvyang/rebetxin-one/service/job/jobtype"
)

type RefundTask struct {
	OpponentId string
	Amount     string
	Memo       string
	AssetId    string
}

func (r *RefundTask) Exec(ctx context.Context) {
	SendTransfer(ctx, r.OpponentId, r.Memo, r.AssetId, convert.String2Decimal(r.Amount))
}

func NewRefundTask(opponentId, amount, memo, AssetId string) (*asynq.Task, error) {
	payload, err := json.Marshal(RefundTask{
		OpponentId: opponentId,
		Amount:     amount,
		Memo:       memo,
		AssetId:    AssetId,
	})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(jobtype.MixinRefundTask, payload, asynq.MaxRetry(5), asynq.Timeout(60*time.Minute)), nil
}
