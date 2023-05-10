package config

import "github.com/zeromicro/go-zero/core/service"

type Config struct {
	service.ServiceConf

	Mixin struct {
		Pin        string
		ClientId   string
		SessionId  string
		PinToken   string
		PrivateKey string
		AppSecret  string
	}
}
