package larkinfra

import (
	"math/rand"
	"time"

	lark "github.com/larksuite/oapi-sdk-go/v3"
)

// 训练营2
const (
	appid     = "cli_a4880a71a638d00e"
	appsecret = "Je8ygYnhXu2y4BlIUCPMbbojRdWtyyaX"
)

var Client *lark.Client

func init() {
	Client = lark.NewClient(appid, appsecret, lark.WithLogReqAtDebug(true),
		lark.WithEnableTokenCache(true))
	rand.Seed(time.Now().Unix())
}
