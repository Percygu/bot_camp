package infra

import (
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"sync/atomic"
)

var cronEngine *cronBase

func init() {
	cronEngine = newCronBase()
}

type CronFun func()

type cronBase struct {
	isStart int32
	c       *cron.Cron
	jobs    map[string]CronFun
}

func newCronBase() *cronBase {
	c := cron.New()
	return &cronBase{
		isStart: 0,
		c:       c,
		jobs:    make(map[string]CronFun),
	}
}

// Register checkout tab in https://crontab.guru/#0_0_*_*_*
//  * * * * *
func Register(cronTab string, fn CronFun) {
	cronEngine.jobs[cronTab] = fn
	logrus.Info("Register Crontab: ", cronTab)
	_, _ = cronEngine.c.AddFunc(cronTab, fn)
}

// StartCronjob 非阻塞
func StartCronjob() {
	if atomic.CompareAndSwapInt32(&cronEngine.isStart, 0, 1) {
		logrus.Info("Start Cronjob")
		cronEngine.c.Start()
	}
}
