package main

import (
	"context"
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github/rotatebot/app"
	"github/rotatebot/config"
	"github/rotatebot/dao"
	"github/rotatebot/infra"
	"github/rotatebot/infra/larkinfra"
	"github/rotatebot/utils"
	"github/rotatebot/view"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

var campName string

func init() {
	flag.StringVar(&campName, "camp_name", "bot_camp2", "camp name")
}

func Init() error {
	utils.InitCampName(campName)
	config.InitConfig()
	if _, ok := config.GetGlobalConf().BotSvrConfig[campName]; !ok {
		log.Errorf("camp_name:%s not exists", campName)
		return fmt.Errorf("camp_name:%s not exists", campName)
	}
	botConfig := utils.GetBotCampConf()
	larkinfra.InitClient(botConfig.AppID, botConfig.AppSecret)
	return nil
}

//go:build=
func main() {
	flag.Parse()
	if err := Init(); err != nil {
		return
	}
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
	botConfig := utils.GetBotCampConf()
	if err := r.Run(":" + strconv.Itoa(botConfig.Port)); err != nil {
		panic(err)
	}
}
