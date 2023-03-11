package view

import (
	"context"
	log "github.com/sirupsen/logrus"
	"github/rotatebot/app"
	"github/rotatebot/utils"

	"github.com/gin-gonic/gin"
	sdkginext "github.com/larksuite/oapi-sdk-gin"
	"github.com/larksuite/oapi-sdk-go/v3/event/dispatcher"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
)

func RouterInit() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	registerGroupOfficial(r)
	return r
}

func registerGroupOfficial(r *gin.Engine) {
	log.Infof(utils.GetBotCampConf().BotVerifyToken)
	log.Infof(utils.GetBotCampConf().BotEncryptedKey)
	handler := dispatcher.NewEventDispatcher(utils.GetBotCampConf().BotVerifyToken, utils.GetBotCampConf().BotEncryptedKey)
	handler.OnP2ChatMemberUserAddedV1(func(ctx context.Context, event *larkim.P2ChatMemberUserAddedV1) error {
		return app.NewJoinGroupHandler().Handle(ctx, event.Event)
	})
	r.POST(utils.GetBotCampConf().EventUrl, sdkginext.NewEventHandlerFunc(handler))
}
