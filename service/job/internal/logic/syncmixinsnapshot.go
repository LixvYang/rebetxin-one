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
	"github.com/lixvyang/rebetxin-one/service/job/internal/svc"
	"github.com/shopspring/decimal"
	"github.com/zeromicro/go-zero/core/logx"
)

var dtmServer = "consul://127.0.0.1:8500/dtmservice"

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
	Select int    `json:"select" comment:"0yes 1no"`
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
	snapshots, err := client.ReadSnapshots(c, "", time.Now(), "", 50)
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

func (l *SyncMixinSnapshotHandler) handlerNewMixinSnapshot(ctx context.Context, snapshot *mixin.Snapshot) error {
	// 1. 根据用户发送的Memo去判断是否进行下一步
	if snapshot.Memo == "" {
		logx.Infow("snapshot.Memo == \"\": ", logx.LogField{Key: "Error: ", Value: "Handle Mixin snapshot error!"})
		return nil
	}
	// 1.1 首先创建一个 mixinsnapshot 同步机器人账单
	// 1.1.1 同步失败了也没关系
	// _, err := l.svcCtx.MixinSnapshotRpc.AddMixinsnapshot(l.ctx, &mixinsnapshotsrv.AddMixinsnapshotReq{
	// 	SnapshotId:     snapshot.SnapshotID,
	// 	AssetId:        snapshot.AssetID,
	// 	OpponentId:     snapshot.OpponentID,
	// 	Amount:         snapshot.Amount.Truncate(8).InexactFloat64(),
	// 	TraceId:        snapshot.TraceID,
	// 	Memo:           snapshot.Memo,
	// 	Type:           snapshot.Type,
	// 	OpeningBalance: snapshot.OpeningBalance.Truncate(8).InexactFloat64(),
	// 	ClosingBalance: snapshot.ClosingBalance.Truncate(8).InexactFloat64(),
	// })
	// if err != nil {
	// 	logx.Errorw("MixinSnapshotRpc.AddMixinsnapshot: ", logx.LogField{Key: "Error: ", Value: err})
	// }

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

	// 4.1 是CNB 则下一步
	// 发送添加更新话题购买系统 这里的具体逻辑由话题购买系统做逻辑
	// TopicpurchaseBusiServer, err := l.svcCtx.Config.TopicpurchaseRPC.BuildTarget()
	// if err != nil {
	// 	logx.Errorw("TopicpurchaseRPC.BuildTarget err: ", logx.LogField{Key: "err: ", Value: err})
	// 	return err
	// }

	// TopicBusiServer, err := l.svcCtx.Config.TopicRPC.BuildTarget()
	// if err != nil {
	// 	logx.Errorw("TopicRPC.BuildTarget err: ", logx.LogField{Key: "err: ", Value: err})
	// 	return err
	// }

	// updateTopicPurchaseReq, updateTopicReq := createReq(&memo, snapshot)
	//这里只举了saga例子，tcc等其他例子基本没啥区别具体可以看dtm官网
	// gid := dtmgrpc.MustGenGid(dtmServer)
	// saga := dtmgrpc.NewSagaGrpc(dtmServer, gid).
	// Add(TopicpurchaseBusiServer+"/pb.topicpurchasesrv/UpdateTopicpurchase", TopicpurchaseBusiServer+"/pb.topicpurchasesrv/UpdateTopicpurchaseRollback", updateTopicPurchaseReq).
	// Add(TopicBusiServer+"/pb.topicsrv/UpdateTopic", TopicBusiServer+"/pb.topicsrv/UpdateTopicRollback", updateTopicReq)

	// err = saga.Submit()
	// dtmimp.FatalIfError(err)
	// if err != nil {
	// 	// TODO
	// 	// 需要将金额返还给用户
	// 	return fmt.Errorf("submit data to  dtm-server err  : %+v \n", err)
	// }

	// 接着更新话题系统 增加金额 Amount

	return nil
}

// func createReq(memo *Memo, snapshot *mixin.Snapshot) (*topicpurchasesrv.UpdateTopicpurchaseReq, *topicsrv.UpdateTopicReq) {
// 	amount := snapshot.Amount.String()
// 	updateTopicPurchaseReq := &topicpurchasesrv.UpdateTopicpurchaseReq{
// 		Action: 0, // 0 buy
// 		Uid:    snapshot.SnapshotID,
// 		Tid:    memo.Tid,
// 	}

// 	updateTopicReq := &topicsrv.UpdateTopicReq{
// 		TotalPrice: amount,
// 		Action:     0, // 0 buy
// 	}

// 	if memo.Select == 0 {
// 		updateTopicReq.YesPrice = amount
// 		updateTopicPurchaseReq.YesPrice = amount
// 	} else {
// 		updateTopicPurchaseReq.NoPrice = amount
// 		updateTopicReq.NoPrice = amount
// 	}

// 	return updateTopicPurchaseReq, updateTopicReq
// }
