package main

import (
	"context"
	"flag"
	"os"

	"github.com/lixvyang/rebetxin-one/service/job/internal/config"
	"github.com/lixvyang/rebetxin-one/service/job/internal/logic"
	"github.com/lixvyang/rebetxin-one/service/job/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"

	"github.com/zeromicro/go-zero/core/conf"
)

var configFile = flag.String("a", "etc/mqueue.yaml", "Specify the config file")

func main() {
	flag.Parse()
	var c config.Config

	conf.MustLoad(*configFile, &c, conf.UseEnv())
	// quit := make(chan os.Signal, 1)
	// signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // 此处不会阻塞
	// log、prometheus、trace、metricsUrl
	if err := c.SetUp(); err != nil {
		panic(err)
	}

	svcContext := svc.NewServiceContext(c)
	ctx := context.Background()
	cronJob := logic.NewCronJob(ctx, svcContext)
	mux := cronJob.Register()

	if err := svcContext.AsynqServer.Run(mux); err != nil {
		logx.WithContext(ctx).Errorf("!!!CronJobErr!!! run err:%+v", err)
		os.Exit(1)
	}

	// <-quit // 阻塞在此，当接收到上述两种信号时才会往下执行
	// logx.Info("Shutdown Server ...")
	// logx.Info("Server exiting")
}
