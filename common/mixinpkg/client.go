package mixinpkg

import (
	"context"
	"flag"
	"log"

	"github.com/fox-one/mixin-sdk-go"
	"github.com/lixvyang/rebetxin-one/common/mixinpkg/internal/config"
	"github.com/zeromicro/go-zero/core/conf"
)

var (
	MixinClient *mixin.Client
	err         error
	configFile  = flag.String("f", "etc/mixin.yaml", "Specify the config file")
)

func init() {

	flag.Parse()
	var c config.Config

	conf.MustLoad(*configFile, &c, conf.UseEnv())

	// log、prometheus、trace、metricsUrl
	if err := c.SetUp(); err != nil {
		panic(err)
	}

	store := &mixin.Keystore{
		ClientID:   c.Mixin.ClientId,
		SessionID:  c.Mixin.SessionId,
		PrivateKey: c.Mixin.PrivateKey,
		PinToken:   c.Mixin.PinToken,
	}

	MixinClient, err = mixin.NewFromKeystore(store)
	if _, err := MixinClient.UserMe(context.Background()); err != nil {
		switch {
		case mixin.IsErrorCodes(err, mixin.Unauthorized, mixin.EndpointNotFound):
			// handle unauthorized error
		case mixin.IsErrorCodes(err, mixin.InsufficientBalance):
			// handle insufficient balance error
		default:
		}
	}

	if err != nil {
		log.Fatal(err)
	}
}
