package main

import (
	"log"

	"github.com/hibiken/asynq"
	"github.com/lixvyang/rebetxin-one/service/job/internal/logic"
)

const redisAddr = "127.0.0.1:6379"

func main() {
	client := asynq.NewClient(asynq.RedisClientOpt{Addr: redisAddr})
	defer client.Close()

	task, err := logic.NewTestStruct("Lixin")
	if err != nil {
		panic(err)
	}

	info, err := client.Enqueue(task)
	if err != nil {
		panic(err)
	}
	log.Printf("Enqued task: id = %s queue=%s", info.ID, info.Queue)
}
