package mixinpkg

import (
	"context"
	"encoding/json"
	"time"

	"github.com/hibiken/asynq"
	"github.com/lixvyang/rebetxin-one/service/job/jobtype"
)

type NewTopicTask struct {
	Message     string
	Title       string
	Intro       string
	RecipientId string
	IconUrl     string
	Tid         string
	AssetId     string
}

func (t *NewTopicTask) Exec(ctx context.Context) {
	if t.Message != "" {
		SendMessage(ctx, t.RecipientId, t.Message)
	}
	SendCards(ctx, t.Title, t.Intro, t.IconUrl, t.Tid, t.RecipientId)
	SendButtonGroup(ctx, t.Tid, t.AssetId, t.RecipientId)
}

func InitNewTopicTask(message, recipientId, iconUrl, intro, title, AssetId, tid string) (*asynq.Task, error) {
	payload, err := json.Marshal(NewTopicTask{
		Message:     message,
		AssetId:     AssetId,
		Title:       title,
		Intro:       intro,
		Tid:         tid,
		RecipientId: recipientId,
		IconUrl:     iconUrl,
	})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(jobtype.MixinNewTopicTask, payload, asynq.MaxRetry(5), asynq.Timeout(60*time.Minute)), nil
}
