package larkinfra

import (
	"context"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	larkwiki "github.com/larksuite/oapi-sdk-go/v3/service/wiki/v2"
	"github.com/sirupsen/logrus"
)

func CreateSpaceMember(ctx context.Context, req *larkwiki.CreateSpaceMemberReq) (*larkwiki.CreateSpaceMemberResp, error) {
	// 发起请求
	resp, err := Client.Wiki.SpaceMember.Create(context.Background(), req)

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
	logrus.Infof(larkcore.Prettify(resp))
	return resp, nil
}

func CreateSpace(ctx context.Context, req *larkwiki.CreateSpaceReq, userToken string) (*larkwiki.CreateSpaceResp, error) {
	// 发起请求
	resp, err := Client.Wiki.Space.Create(context.Background(),
		req,
		larkcore.WithUserAccessToken(userToken))

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
	logrus.Infof(larkcore.Prettify(resp))
	return resp, nil
}
