package app

import (
	"context"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	"github.com/sirupsen/logrus"
	"github/rotatebot/infra"
	"github/rotatebot/infra/larkinfra"
	"github/rotatebot/proto"
	"github/rotatebot/utils"
	"strings"
	"time"
)

const dayOff = 7

// RemindTalk 群聊remind说话功能
func RemindTalk(ctx context.Context) {
	infra.Register("*/5 * * * *", func() {
		logrus.Info("Trigger Remind Check")
		var pageToken string
		utils.MaxLoopController(ctx, 500, func() (bool, error) {
			resp, err := larkinfra.ListBotChats(ctx, pageToken)
			if err != nil {
				logrus.Errorf("list bots err:%+v", err)
				return false, err
			}
			for _, chat := range resp.Data.Items {
				if msg, ok := checkNeedRemind(ctx, chat); !ok {
					logrus.Infof("chatID=%s,chat=%s skip remind with reason:%s", *chat.ChatId, *chat.Name, msg)
					continue
				}
				go sendRemindTalkMsg(ctx, chat)
			}
			// last thing
			pageToken = *resp.Data.PageToken
			if !*resp.Data.HasMore || len(pageToken) == 0 {
				logrus.Info("No More Data To Check,Finish")
				return true, nil
			}
			return false, nil
		})
	})
}

func checkNeedRemind(ctx context.Context, chat *larkim.ListChat) (string, bool) {
	now := time.Now()
	start := now.AddDate(0, 0, -dayOff)
	msg, err := larkinfra.GetGroupMessage(ctx, start.Unix(), now.Unix(), *chat.ChatId)
	if err != nil {
		// error放过
		return "获取消息失败", false
	}
	// 测试环境专属旁路逻辑，后续不再走
	if utils.IsTestEnv() {
		allowMap := utils.GetTestEnvRemindAllowList()
		if _, ok := allowMap[*chat.ChatId]; !ok {
			return "测试环境下不属于测试群", false
		}
	}

	// 校验是否是白名单
	whiteMap := utils.GetRemindChatWhiteIDs()
	if _, ok := whiteMap[*chat.ChatId]; ok {
		return "豁免，在白名单中", false
	}

	// 跳过测试群
	if strings.Contains(*chat.Name, "测试") {
		return "跳过测试群", false
	}

	// todo: 兜底【机器人创建群后先走到该定时任务，后发出了初始通知，则可能两条消息】
	// 目前没接口获取群创建时间

	if len(msg.Items) == 0 {
		return "", true
	}
	return "有范围内消息，不需要提醒", false
}

// sendRemindTalkMsg 提醒消息
func sendRemindTalkMsg(ctx context.Context, chat *larkim.ListChat) {
	logrus.Info("Send Msg To ", *chat.Name)
	msg := proto.NewRemindTalkPost(dayOff)
	im := larkinfra.NewLarkInstance(*chat.ChatId)
	_ = im.SendLarkGroup(ctx, msg, larkinfra.PostMsg)
}
