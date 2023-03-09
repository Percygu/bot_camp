package larkinfra

import (
	"context"
	"fmt"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	larktask "github.com/larksuite/oapi-sdk-go/v3/service/task/v1"
	"github.com/sirupsen/logrus"
)

func CreateTask(ctx context.Context, req *larktask.CreateTaskReq) (*larktask.CreateTaskResp, error) {
	resp, err := Client.Task.Task.Create(context.Background(), req)
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
