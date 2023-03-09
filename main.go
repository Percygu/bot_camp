package main

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github/rotatebot/app"
	"github/rotatebot/config"
	"github/rotatebot/dao"
	"github/rotatebot/infra"
	"github/rotatebot/utils"
	"github/rotatebot/view"
	"os"
	"os/signal"
	"syscall"
)

//go:build=
func main() {
	config.InitConfig()
	log.SetReportCaller(true)
	log.Infof("Is Test Env:%t", utils.IsTestEnv())
	ctx, cancel := context.WithCancel(context.Background())
	infra.StartCronjob()
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGINT)
	go func() {
		<-c
		fmt.Println("\r- Ctrl+C pressed in Terminal")
		cancel()
		os.Exit(0)
	}()
	dao.InitDB()
	app.RegisterTask(ctx)
	r := view.RouterInit()
	if err := r.Run(":8080"); err != nil {
		panic(err)
	}
}
