package logic

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/fox-one/mixin-sdk-go"
	"github.com/hibiken/asynq"
	"github.com/lixvyang/rebetxin-one/common/constant"
	"github.com/lixvyang/rebetxin-one/common/errorx"
	"github.com/lixvyang/rebetxin-one/model"
	"github.com/lixvyang/rebetxin-one/service/job/internal/svc"
	"github.com/shopspring/decimal"
	"github.com/zeromicro/go-zero/core/logx"
)

// SyncMixinSnapshotHandler
type SyncMixinSnapshotHandler struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSyncMixinSnapshotHandler(svcCtx *svc.ServiceContext) *SyncMixinSnapshotHandler {
	c := context.Background()
	return &SyncMixinSnapshotHandler{
		ctx:    c,
		svcCtx: svcCtx,
	}
}

// every one minute exec : if return err != nil , asynq will retry
func (l *SyncMixinSnapshotHandler) ProcessTask(ctx context.Context, _ *asynq.Task) error {
	fmt.Printf("shcedule job demo -----> every one second exec \n")
	createdAt, err := getTopSnapshotCreatedAt(l.svcCtx.MixinClient, ctx)
	if err != nil {
		return err
	}
	stats := &Stats{createdAt}
	l.sendTopCreatedAtToChannel(ctx, stats)
	return nil
}

type Memo struct {
	Tid    string `json:"tid"`
	Select int64  `json:"select" comment:"0yes 1no"`
}

type Stats struct {
	preCreatedAt time.Time
}

func (s *Stats) getPrevSnapshotCreatedAt() time.Time {
	return s.preCreatedAt
}

func (s *Stats) updatePrevSnapshotCreatedAt(time time.Time) {
	s.preCreatedAt = time
}

func getTopSnapshotCreatedAt(client *mixin.Client, c context.Context) (time.Time, error) {
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
	snapshots, err := client.ReadSnapshots(c, "", time.Now(), "", 50)
	if err != nil {
		return nil, err
	}
	return snapshots, nil
}

func (l *SyncMixinSnapshotHandler) sendTopCreatedAtToChannel(ctx context.Context, stats *Stats) {
	preCreatedAt := stats.getPrevSnapshotCreatedAt()
	snapshots, err := getTopHundredCreated(l.svcCtx.MixinClient, ctx)
	if err != nil {
		log.Printf("getTopHundredCreated error")
		return
	}
	var wg sync.WaitGroup
	for _, snapshot := range snapshots {
		if snapshot.CreatedAt.After(preCreatedAt) {
			stats.updatePrevSnapshotCreatedAt(snapshot.CreatedAt)
			if snapshot.Amount.Cmp(decimal.NewFromInt(0)) == 1 && snapshot.Type == "transfer" {
				wg.Add(1)
				go func(s *mixin.Snapshot) {
					defer wg.Done()
					_ = l.handlerNewMixinSnapshot(ctx, s)
				}(snapshot)
			}
		}
	}
	wg.Wait()
}

func (l *SyncMixinSnapshotHandler) handlerNewMixinSnapshot(ctx context.Context, snapshot *mixin.Snapshot) (err error) {
	sp := new(model.Snapshot)
	userPurchase := new(model.Topicpurchase)
	// 1. 根据用户发送的Memo去判断是否进行下一步
	if snapshot.Memo == "" {
		logx.Infow("snapshot.Memo == \"\": ", logx.LogField{Key: "Error: ", Value: "Handle Mixin snapshot error!"})
		return nil
	}

	// 2. 解码用户发送的Memo
	// 2.1 解码Memo失败 则退出
	memoMsg, err := base64.StdEncoding.DecodeString(snapshot.Memo)
	if err != nil {
		logx.Errorw("base64.StdEncoding.DecodeString: ", logx.LogField{Key: "Error: ", Value: err})
		return err
	}

	var memo Memo
	if err := json.Unmarshal(memoMsg, &memo); err != nil {
		return err
	}

	// 1.1 首先创建一个 snapshot 同步机器人账单
	// 1.1.1 同步失败了也没关系
	sp, err = l.svcCtx.SnapshotModel.FindOneByTraceId(l.ctx, snapshot.TraceID)
	if err != nil {
		if err == model.ErrNotFound {
			logx.Errorw("SnapshotModel.FindOneByTraceId", logx.LogField{Key: "Error", Value: err.Error()})
			l.svcCtx.SnapshotModel.Insert(l.ctx, &model.Snapshot{
				Uid:     snapshot.Sender,
				TraceId: snapshot.TraceID,
				Tid:     memo.Tid,
				End:     1,
			})
			sp, _ = l.svcCtx.SnapshotModel.FindOneByTraceId(l.ctx, snapshot.TraceID)
		} else {
			logx.Errorw("SnapshotModel.FindOneByTraceId", logx.LogField{Key: "Error", Value: err.Error()})
			return err
		}
	}
	sp.End = 1
	_ = l.svcCtx.SnapshotModel.Update(l.ctx, sp)

	// @2023.4.27 只使用CNB做价格  存储的也只有CNB的个数
	// type Memo struct {
	// 	Tid    string `json:"tid"`
	// 	Select int    `json:"select" comment:"0yes 1no"`
	// }
	//
	// 4.0 判断是不是CNB 不是CNB 则退出
	if snapshot.AssetID != "CNB" {
		logx.Info("snapshot.AssetID is not CNB.")
	}

	// 4 是CNB 则下一步
	// 4.1 更新话题系统
	// 4.2 更新话题购买系统
	// 4.3 通知用户
	topic, err := l.svcCtx.TopicModel.FindOneByTid(l.ctx, memo.Tid)
	if err != nil {
		logx.Errorw("TopicModel.FindOneByTid: ", logx.LogField{Key: "Error", Value: err.Error()})
		return err
	}

	userPurchase, err = l.svcCtx.TopicPurchaseModel.FindOneByUidTid(l.ctx, sp.Uid, sp.Tid)
	if err != nil {
		if err == model.ErrNotFound {
			l.svcCtx.TopicPurchaseModel.Insert(l.ctx, &model.Topicpurchase{
				Uid: sp.Uid,
				Tid: sp.Tid,
			})
			userPurchase, _ = l.svcCtx.TopicPurchaseModel.FindOneByUidTid(l.ctx, sp.Uid, sp.Tid)
		} else {
			logx.Errorw("TopicPurchaseModel.FindOneByUidTid", logx.LogField{Key: "Error", Value: err.Error()})
		}
	}

	return l.HandleUserPurchase(memo.Select, userPurchase, topic, snapshot)
}

func (l *SyncMixinSnapshotHandler) HandleUserPurchase(Select int64, userPurchase *model.Topicpurchase, topic *model.Topic, snapshot *mixin.Snapshot) (err error) {
	switch Select {
	case 0:
		// 更新用户购买系统逻辑
		userPurchase.YesPrice = userPurchase.YesPrice.Add(snapshot.Amount)
		err := l.svcCtx.TopicPurchaseModel.Update(l.ctx, userPurchase)
		if err != nil {
			logx.Errorw("TopicPurchaseModel.Update", logx.LogField{Key: "Error: ", Value: err.Error()})
			return errorx.NewDefaultError("TopicPurchaseModel.Update Error!")
		}
		topic.YesPrice = topic.YesPrice.Add(snapshot.Amount)
		topic.TotalPrice = topic.TotalPrice.Add(snapshot.Amount)
		topic.YesRatio = topic.YesPrice.Div(topic.TotalPrice).Mul(decimal.NewFromInt(100))
		topic.NoRatio = topic.NoPrice.Div(topic.TotalPrice).Mul(decimal.NewFromInt(100))
		err = l.svcCtx.TopicModel.Update(l.ctx, topic)
		if err != nil {
			logx.Errorw("TopicModel.Update", logx.LogField{Key: "Error: ", Value: err.Error()})
			return errorx.NewDefaultError("Topic.Update Error!")
		}
	case 1:
		// 更新用户购买系统逻辑
		userPurchase.NoPrice = userPurchase.YesPrice.Add(snapshot.Amount)
		err := l.svcCtx.TopicPurchaseModel.Update(l.ctx, userPurchase)
		if err != nil {
			logx.Errorw("TopicPurchaseModel.Update", logx.LogField{Key: "Error: ", Value: err.Error()})
			return errorx.NewDefaultError("TopicPurchaseModel.Update Error!")
		}
		topic.NoPrice = topic.NoPrice.Add(snapshot.Amount)
		topic.TotalPrice = topic.TotalPrice.Add(snapshot.Amount)
		topic.YesRatio = topic.YesPrice.Div(topic.TotalPrice).Mul(decimal.NewFromInt(100))
		topic.NoRatio = topic.NoPrice.Div(topic.TotalPrice).Mul(decimal.NewFromInt(100))
		err = l.svcCtx.TopicModel.Update(l.ctx, topic)
		if err != nil {
			logx.Errorw("TopicModel.Update", logx.LogField{Key: "Error: ", Value: err.Error()})
			return errorx.NewDefaultError("Topic.Update Error!")
		}
	default:
		return err
	}
	return nil
}
