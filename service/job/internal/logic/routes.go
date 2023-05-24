package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/fox-one/mixin-sdk-go"
	"github.com/hibiken/asynq"
	"github.com/lixvyang/rebetxin-one/common/constant"
	"github.com/lixvyang/rebetxin-one/service/job/internal/svc"
	"github.com/lixvyang/rebetxin-one/service/job/jobtype"
	"github.com/shopspring/decimal"
	"github.com/zeromicro/go-zero/core/logx"
)

type CronJob struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCronJob(ctx context.Context, svcCtx *svc.ServiceContext) *CronJob {
	return &CronJob{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// register job
func (l *CronJob) Register() *asynq.ServeMux {
	go l.InitCron()

	mux := asynq.NewServeMux()

	//scheduler job
	mux.Handle(jobtype.MixinSnapshotSyncRecord, NewSyncMixinSnapshotHandler(l.ctx, l.svcCtx))
	mux.Handle(jobtype.MixinTestTask, NewTestHandler())

	return mux
}

type Memo struct {
	Tid    string `json:"tid"`
	Select int64  `json:"select" comment:"0yes 1no"`
}

type Stats struct {
	lastProcessed time.Time
}

func (s *Stats) getPrevSnapshotCreatedAt() time.Time {
	return s.lastProcessed
}

func (s *Stats) updatePrevSnapshotCreatedAt(time time.Time) {
	s.lastProcessed = time
}

func getLastedSnapshot(client *mixin.Client, c context.Context) (time.Time, error) {
	snapshots, err := client.ReadSnapshots(c, constant.CNB_ASSET_ID, time.Now(), "", 50)
	fmt.Println(len(snapshots))
	if err != nil {
		return time.Now(), err
	}
	if len(snapshots) == 0 {
		return time.Now(), nil
	}
	return snapshots[0].CreatedAt, nil
}

func getTopHundredCreated(client *mixin.Client, c context.Context) ([]*mixin.Snapshot, error) {
	snapshots, err := client.ReadSnapshots(c, constant.CNB_ASSET_ID, time.Now(), "", 50)
	if err != nil {
		return nil, err
	}
	return snapshots, nil
}

func (l *CronJob) InitCron() {
	ctx := context.Background()
	createdAt, err := getLastedSnapshot(l.svcCtx.MixinClient, ctx)
	if err != nil {
	}
	stats := &Stats{createdAt}

	for {
		preCreatedAt := stats.getPrevSnapshotCreatedAt()

		snapshots, err := getTopHundredCreated(l.svcCtx.MixinClient, ctx)
		if err != nil {
			logx.Error("getTopHundredCreated error")
		}

		for _, snapshot := range snapshots {
			if snapshot.CreatedAt.After(preCreatedAt) {
				stats.updatePrevSnapshotCreatedAt(snapshot.CreatedAt)
				if snapshot.Amount.Cmp(decimal.NewFromInt(0)) == 1 && snapshot.Type == "transfer" {
					task, err := NewMixinSnapshot(snapshot)
					l.svcCtx.AsynqClient.Enqueue(task)
					if err != nil {
						logx.Errorw("AsynqClient.Enqueue Error: ", logx.LogField{Key: "Error:", Value: err.Error()})
					}
				}
			} else {
				break
			}
		}
		time.Sleep(3 * time.Second)
	}
}

func NewMixinSnapshot(s *mixin.Snapshot) (*asynq.Task, error) {
	payload, err := json.Marshal(*s)
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(jobtype.MixinSnapshotSyncRecord, payload), nil
}
