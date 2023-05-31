package main

import (
	"flag"
	"fmt"

	"github.com/lixvyang/rebetxin-one/service/betxin/internal/config"
	"github.com/lixvyang/rebetxin-one/service/betxin/internal/handler"
	"github.com/lixvyang/rebetxin-one/service/betxin/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/betxin-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf, rest.WithCors(fmt.Sprintf("http://%s:4000", c.Host)))
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
