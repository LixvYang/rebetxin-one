package svc

import (
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/lixvyang/rebetxin-one/service/job/internal/config"
)

func newAsynqServer(c config.Config) *asynq.Server {
	return asynq.NewServer(
		asynq.RedisClientOpt{Addr: c.Redis.Host, Password: c.Redis.Pass},
		asynq.Config{
			IsFailure: func(err error) bool {
				fmt.Printf("asynq server exec task IsFailure ======== >>>>>>>>>>>  err : %+v \n", err)
				return true
			},
			Concurrency: 1, //max concurrent process job task num
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
		},
	)
}
