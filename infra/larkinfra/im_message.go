package larkinfra

import (
	"context"
	"fmt"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	"github.com/sirupsen/logrus"
	"github/rotatebot/utils"
	"strconv"
	"time"
)

type LarkIm struct {
	targetChatID string //发送的目标群
}

func NewLarkInstance(targetChatID string) *LarkIm {
	return &LarkIm{targetChatID: targetChatID}
}

type MsgType string

const (
	TextMsg        MsgType = "text"
	PostMsg        MsgType = "post"
	InteractiveMsg MsgType = `interactive`
)

// SendLarkGroup 发送到群消息
func (l *LarkIm) SendLarkGroup(ctx context.Context, msg string, msgType MsgType) error {
	//https://open.feishu.cn/document/uAjLw4CM/ukTMukTMukTM/im-v1/message/create_json
	logrus.Infof("send message is %s", msg)
	req := larkim.NewCreateMessageReqBuilder().
		ReceiveIdType(`chat_id`).
		Body(larkim.NewCreateMessageReqBodyBuilder().
			ReceiveId(l.targetChatID).
			MsgType(string(msgType)).
			Content(msg).
			Build()).
		Build()
	resp, err := Client.Im.Message.Create(ctx, req)
	if err != nil {
		logrus.Errorf("%v", err)
		return err
	}

	if !resp.Success() {
		logrus.Errorf("%v", err)
		return err
	}

	logrus.Infof("send success")
	return nil
}

func FetchAllGroupMembers(ctx context.Context, chatID string, memberType utils.MemberType) []*larkim.ListMember {
	var resp = make([]*larkim.ListMember, 0)
	req := larkim.NewGetChatMembersReqBuilder().
		ChatId(chatID).
		MemberIdType(string(memberType)).
		PageSize(50).
		Build()
	ans, err := Client.Im.ChatMembers.Get(ctx, req)
	if err != nil {
		return resp
	}
	//fmt.Println(larkcore.Prettify(ans))
	hasmore := *ans.Data.HasMore
	for _, item := range ans.Data.Items {
		if item.Name == nil {
			continue
		}
		resp = append(resp, item)
	}
	if !hasmore || len(resp) == 0 {
		logrus.Infof("Finish Find Member for chatID=%s,got len=%d", chatID, len(resp))
		return resp
	}
	if hasmore {
		time.Sleep(time.Second)
		resp = append(resp, FetchAllGroupMembers(ctx, chatID, memberType)...)
	}
	return resp
}

func GetGroupMessage(ctx context.Context, startUnix, endUnix int64, chatID string) (*larkim.ListMessageRespData, error) {
	startUnix_, endUnix_ := strconv.FormatInt(startUnix, 10), strconv.FormatInt(endUnix, 10)
	req := larkim.NewListMessageReqBuilder().ContainerIdType("chat").
		ContainerId(chatID).
		StartTime(startUnix_).
		EndTime(endUnix_).
		PageSize(20).
		Build()
	resp, err := Client.Im.Message.List(ctx, req)
	if err != nil {
		logrus.Errorf("%+v", err)
		return nil, err
	}

	// 处理错误
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	// 服务端错误处理
	if !resp.Success() {
		fmt.Println(resp.Code, resp.Msg, resp.RequestId())
		return nil, fmt.Errorf(resp.Msg)
	}
	return resp.Data, nil
}
