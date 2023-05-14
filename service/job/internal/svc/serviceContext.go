package svc

import (
	"github.com/fox-one/mixin-sdk-go"
	"github.com/hibiken/asynq"
	"github.com/lixvyang/rebetxin-one/service/job/internal/config"
	"github.com/zeromicro/go-zero/core/logx"
)

type ServiceContext struct {
	Config      config.Config
	AsynqServer *asynq.Server
	AsynqClient *asynq.Client
	MixinClient *mixin.Client

	// MixinSnapshotRpc mixinsnapshotsrv.Mixinsnapshotsrv
	// TopicpurchaseRpc topicpurchasesrv.Topicpurchasesrv
}

func NewServiceContext(c config.Config) *ServiceContext {
	store := &mixin.Keystore{
		ClientID:   c.Mixin.ClientId,
		SessionID:  c.Mixin.SessionId,
		PrivateKey: c.Mixin.PrivateKey,
		PinToken:   c.Mixin.PinToken,
	}

	mixinClient, err := mixin.NewFromKeystore(store)
	if err != nil {
		logx.Errorw("init mixinclient err: ", logx.LogField{Key: "Error: ", Value: err})
		panic(err)
	}

	return &ServiceContext{
		Config:      c,
		AsynqServer: newAsynqServer(c),
		AsynqClient: newAsynqClient(c),
		MixinClient: mixinClient,
		// MixinSnapshotRpc: mixinsnapshotsrv.NewMixinsnapshotsrv(zrpc.MustNewClient(c.MixinSnapshotRPC)),
		// TopicpurchaseRpc: topicpurchasesrv.NewTopicpurchasesrv(zrpc.MustNewClient(c.TopicpurchaseRPC)),
	}
}
