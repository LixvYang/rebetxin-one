package logic

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hibiken/asynq"
	"github.com/lixvyang/rebetxin-one/service/job/jobtype"
)

type TestHandler struct{}

type TestStruct struct {
	Name string
}

func (th *TestHandler) ProcessTask(ctx context.Context, t *asynq.Task) error {
	var p TestStruct
	json.Unmarshal(t.Payload(), &p)
	log.Printf("%+v", p)
	return nil
}

func NewTestHandler() *TestHandler {
	return &TestHandler{}
}

func NewTestStruct(name string) (*asynq.Task, error) {
	payload, err := json.Marshal(TestStruct{Name: name})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(jobtype.MixinTestTask, payload), nil
}
