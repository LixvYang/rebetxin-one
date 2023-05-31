package config

import (
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	service.ServiceConf
	Redis redis.RedisConf

	CacheRedis cache.CacheConf
	Mysql      struct {
		DNS string
	}

	Mixin struct {
		Pin        string
		ClientId   string
		SessionId  string
		PinToken   string
		PrivateKey string
		AppSecret  string
	}

	MixinSrvRPC zrpc.RpcClientConf
}
