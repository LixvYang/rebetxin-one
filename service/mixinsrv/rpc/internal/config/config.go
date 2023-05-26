package config

import "github.com/zeromicro/go-zero/zrpc"

type Config struct {
	zrpc.RpcServerConf
	
	Mixin struct {
		Pin        string
		ClientId   string
		SessionId  string
		PinToken   string
		PrivateKey string
		AppSecret  string
	}
}
