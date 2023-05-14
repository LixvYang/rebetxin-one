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
	configFile  = flag.String("f", "etc/mqueue.yaml", "Specify the config file")
	Config      config.Config
)

func init() {

	flag.Parse()

	conf.MustLoad(*configFile, &Config, conf.UseEnv())

	// log、prometheus、trace、metricsUrl
	if err := Config.SetUp(); err != nil {
		panic(err)
	}

	store := &mixin.Keystore{
		ClientID:   Config.Mixin.ClientId,
		SessionID:  Config.Mixin.SessionId,
		PrivateKey: Config.Mixin.PrivateKey,
		PinToken:   Config.Mixin.PinToken,
	}

	MixinClient, err = mixin.NewFromKeystore(store)
	if _, err := MixinClient.UserMe(context.Background()); err != nil {
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
}
