package larkinfra

import (
	"context"
	"fmt"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	"github.com/sirupsen/logrus"
)

func CreateGroup(ctx context.Context, req *larkim.CreateChatReq) (*larkim.CreateChatResp, error) {

	// 发起请求
	// 如开启了SDK的Token管理功能，就无需在请求时调用larkcore.WithTenantAccessToken("-xxx")来手动设置租户Token了
	resp, err := Client.Im.Chat.Create(context.Background(), req)

	// 处理错误
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	// 服务端错误处理
	if !resp.Success() {
		logrus.Error(resp.Code, resp.Msg, resp.RequestId())
		return nil, err
	}

	// 业务处理
	fmt.Println(larkcore.Prettify(resp))
	return resp, nil
}

func ListBotChats(ctx context.Context, pageToken string) (*larkim.ListChatResp, error) {
	// 创建请求对象
	req := larkim.NewListChatReqBuilder().
		UserIdType(`open_id`).
		PageToken(pageToken).
		PageSize(100).
		Build()

	// 发起请求
	// 如开启了SDK的Token管理功能，就无需在请求时调用larkcore.WithTenantAccessToken("-xxx")来手动设置租户Token了
	resp, err := Client.Im.Chat.List(context.Background(), req)

	// 处理错误
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	// 服务端错误处理
	if !resp.Success() {
		logrus.Error(resp.Code, resp.Msg, resp.RequestId())
		return nil, fmt.Errorf(resp.Msg)
	}

	// 业务处理
	logrus.Infof(larkcore.Prettify(resp))
	return resp, nil
}
