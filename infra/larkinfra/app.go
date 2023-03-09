package larkinfra

import (
	"context"
	"encoding/json"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	"github.com/sirupsen/logrus"
)

type VisibilityRequest struct {
	AppId           string        `json:"app_id"`
	AddUsers        []*User       `json:"add_users"`
	DelUsers        []interface{} `json:"del_users"`
	IsVisiableToAll int           `json:"is_visiable_to_all"`
	AddDepartments  []string      `json:"add_departments"`
	DelDepartments  []interface{} `json:"del_departments"`
}
type User struct {
	OpenId string `json:"open_id"`
	UserId string `json:"user_id"`
}

// UpdateVisibility 应用可见性
func UpdateVisibility(ctx context.Context, userIDs []string) {
	users := make([]*User, 0)
	for _, d := range userIDs {
		users = append(users, &User{
			UserId: d,
		})
	}

	vis := &VisibilityRequest{
		AppId:           appid,
		AddUsers:        users,
		IsVisiableToAll: 0,
		AddDepartments:  []string{"od-00f02bd1248978fec7311904fba21f01", "79dc8e4d9bd8d4db"},
	}
	vb, _ := json.Marshal(vis)
	resp, err := Client.Post(ctx, "/open-apis/application/v3/app/update_visibility", string(vb), larkcore.AccessTokenTypeTenant)
	if err != nil {
		logrus.Error(err)
	}
	rb, _ := json.Marshal(resp)
	logrus.Infof("%s", rb)
}
