package logic

import (
	"context"
	"time"

	"github.com/hibiken/asynq"
	"github.com/lixvyang/rebetxin-one/service/job/internal/svc"
	"github.com/lixvyang/rebetxin-one/service/job/jobtype"
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
	l.InitCron()

	mux := asynq.NewServeMux()

	//scheduler job
	mux.Handle(jobtype.MixinSnapshotSyncRecord, NewSyncMixinSnapshotHandler(l.svcCtx))
	mux.Handle(jobtype.MixinNewTopicTask, NewTopicTask(l.svcCtx))

	//queue job , asynq support queue job
	// wait you fill..

	return mux
}

func (l *CronJob) InitCron() {
	go func() {
		// 创建任务
		task := asynq.NewTask(jobtype.MixinSnapshotSyncRecord, nil)
		// 使用 time.Ticker 实现定时任务
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		for range ticker.C {
			_, err := l.svcCtx.AsynqClient.Enqueue(task)
			if err != nil {
				panic(err)
			}
		}
	}()
}
