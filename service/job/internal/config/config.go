package config

import (
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

type Config struct {
	service.ServiceConf
	Redis redis.RedisConf

	Mixin struct {
		Pin        string
		ClientId   string
		SessionId  string
		PinToken   string
		PrivateKey string
		AppSecret  string
	}
}
