package svc

import (
	"github.com/fox-one/mixin-sdk-go"
	"github.com/hibiken/asynq"
	"github.com/lixvyang/rebetxin-one/service/betxin/model"
	"github.com/lixvyang/rebetxin-one/service/job/internal/config"
	"github.com/lixvyang/rebetxin-one/service/mixinsrv/rpc/mixinsrv"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config      config.Config
	AsynqServer *asynq.Server
	AsynqClient *asynq.Client
	MixinClient *mixin.Client

	SnapshotModel      model.SnapshotModel
	TopicPurchaseModel model.TopicpurchaseModel
	TopicModel         model.TopicModel

	MixinSrvRPC mixinsrv.Mixinsrv
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.Mysql.DNS)
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
		Config:             c,
		AsynqServer:        newAsynqServer(c),
		AsynqClient:        newAsynqClient(c),
		MixinClient:        mixinClient,
		SnapshotModel:      model.NewSnapshotModel(conn, c.CacheRedis),
		TopicPurchaseModel: model.NewTopicpurchaseModel(conn, c.CacheRedis),
		TopicModel:         model.NewTopicModel(conn, c.CacheRedis),
		MixinSrvRPC:        mixinsrv.NewMixinsrv(zrpc.MustNewClient(c.MixinSrvRPC)),
	}
}
