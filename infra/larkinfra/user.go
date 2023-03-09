package larkinfra

import (
	"context"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	"github.com/sirupsen/logrus"
	"github/rotatebot/proto"
	"time"
)

type LarkUser struct {
	userAssessToken string
	refreshToken    string
}

func NewLarkUser() *LarkUser {
	return &LarkUser{}
}

func (u *LarkUser) GetMyToken() string {
	return u.userAssessToken
}

// todo: 接入回调
func (u *LarkUser) GetMyAssessTokenByRequest(ctx context.Context, code string) string {
	if len(u.userAssessToken) > 0 {
		return u.userAssessToken
	} else {
		u.requestMyAssessToken(ctx, code)
	}
	go func() {
		c := time.Tick(time.Minute)
		for _ = range c {
			u.refreshMyAssessToken(ctx, code)
		}
	}()
	return u.userAssessToken
}

func (c *LarkUser) refreshMyAssessToken(ctx context.Context, code string) {
	if len(c.refreshToken) == 0 {
		c.requestMyAssessToken(ctx, code)
		return
	}

	request := map[string]interface{}{
		"grant_type":    "refresh_token",
		"refresh_token": c.refreshToken,
	}
	resp, err := Client.Post(ctx, "/open-apis/authen/v1/refresh_access_token", request, larkcore.AccessTokenTypeApp)
	if err != nil {
		logrus.Error(err)
		return
	}
	body := &proto.TokenResp{}
	err = resp.JSONUnmarshalBody(body, &larkcore.Config{})
	if err != nil {
		logrus.Error(err)
		return
	}
}

func (c *LarkUser) requestMyAssessToken(ctx context.Context, code string) {
	request := map[string]interface{}{
		"grant_type": "authorization_code",
		"code":       code,
	}
	resp, err := Client.Post(ctx, "/open-apis/authen/v1/access_token", request, larkcore.AccessTokenTypeApp)
	if err != nil {
		logrus.Error(err)
		return
	}
	body := &proto.TokenResp{}
	err = resp.JSONUnmarshalBody(body, &larkcore.Config{})
	if err != nil {
		logrus.Error(err)
		return
	}
	c.refreshToken = body.Data.RefreshToken
	c.userAssessToken = body.Data.AccessToken
	return
}
