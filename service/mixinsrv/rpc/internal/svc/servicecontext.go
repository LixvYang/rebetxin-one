package svc

import (
	"context"
	"log"

	"github.com/fox-one/mixin-sdk-go"
	"github.com/lixvyang/rebetxin-one/service/mixinsrv/rpc/internal/config"
)

type ServiceContext struct {
	Config      config.Config
	MixinClient *mixin.Client
	Pin string
}

func NewServiceContext(c config.Config) *ServiceContext {

	store := &mixin.Keystore{
		ClientID:   c.Mixin.ClientId,
		SessionID:  c.Mixin.SessionId,
		PrivateKey: c.Mixin.PrivateKey,
		PinToken:   c.Mixin.PinToken,
	}

	client, err := mixin.NewFromKeystore(store)
	if _, err := client.UserMe(context.Background()); err != nil {
		switch {
		case mixin.IsErrorCodes(err, mixin.Unauthorized, mixin.EndpointNotFound):
			// handle unauthorized error
			panic(err)
		case mixin.IsErrorCodes(err, mixin.InsufficientBalance):
			// handle insufficient balance error
			panic(err)
		default:
		}
	}

	if err != nil {
		log.Fatal(err)
	}

	return &ServiceContext{
		Config:      c,
		MixinClient: client,
		Pin: c.Mixin.Pin,
	}
}
