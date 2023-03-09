package view

import (
	"context"
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
	handler := dispatcher.NewEventDispatcher(utils.BotVerifyToken, utils.BotEncryptedKey)
	handler.OnP2ChatMemberUserAddedV1(func(ctx context.Context, event *larkim.P2ChatMemberUserAddedV1) error {
		return app.NewJoinGroupHandler().Handle(ctx, event.Event)
	})
	r.POST("/webhook/bot2/event", sdkginext.NewEventHandlerFunc(handler))
}
