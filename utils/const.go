package utils

import (
	"os/user"
	"strings"
)

const (
	// 测试用
	Fish    = "ou_88d59ec35b63b8c9212702b0804886c0"
	FishUID = "5efg94ff"
	PengGe  = "ou_1b73b5f130001a43c22387149a71dcc4"
	All     = "all"
)

const (
	BotOpenID               = "ou_35ac11238b92d3ab3ad98547d3c7fd47" // 机器人openid
	StudentOpenDepartmentID = "od-1aad208d5b2b4b4fd3dd8ee26d8abce6" // 学生部门id
	StudentDepartmentID     = "79dc8e4d9bd8d4db"
	MentorDepartmentID      = "1"
	MentorOpenDepartmentID  = "od-00f02bd1248978fec7311904fba21f01" // 导师团部门id
	BotEncryptedKey         = "yYaeTpNE9LntvFubmILegh4i8JUTsF4A"
	BotVerifyToken          = "TsbOh0kxubaQlOMaIDWeAg7rs2Y6d4z8"
)

type MemberType string

const (
	OpenIDType MemberType = "open_id"
	UserIDType MemberType = "user_id"
)

func GetWhiteOpenIDs() []string {
	if IsTestEnv() {
		return testWhiteOpenIDs
	}
	return whiteOpenIDs
}

func Get1V1MentorIDs() []string {
	if IsTestEnv() {
		return []string{
			"ou_88d59ec35b63b8c9212702b0804886c0", //牛哥
			"ou_9d30ad8e58e9137696bf502825d50518", //小林哥
			"ou_1b73b5f130001a43c22387149a71dcc4", //鹏哥
			"ou_20c871c0783883ced66213f17ed0cd64", // 小鱼
		}
	}
	return whiteOpenIDs
}

func GetWorkshopPrepareChatMentors() []string {
	return Get1V1MentorIDs()
}

//todo: whitelist 可配置
//OpenID:同一个 User ID 在不同应用中的 Open ID 不同
var testWhiteOpenIDs = []string{
	//"ou_6e75a4323ad0eef42ae349d71203969e", //狼哥
	//"ou_bcd4e27aad8d4cda70bb0d0b091e487c", //虎哥
	//"ou_9a1afb986757fa17ee0d4d6ba2289d72", //牛哥
}
var whiteOpenIDs = []string{
	"ou_88d59ec35b63b8c9212702b0804886c0", //牛哥
	"ou_9d30ad8e58e9137696bf502825d50518", //小林哥
	"ou_1b73b5f130001a43c22387149a71dcc4", //鹏哥
	"ou_20c871c0783883ced66213f17ed0cd64", // 小鱼
}

func GetTaskFocusIDs() []string {
	return focusOpenIDs
}

//todo:根据部门获取
var focusOpenIDs = []string{
	"ou_88d59ec35b63b8c9212702b0804886c0", //牛哥
	"ou_9d30ad8e58e9137696bf502825d50518", //小林哥
	"ou_1b73b5f130001a43c22387149a71dcc4", //鹏哥
	"ou_20c871c0783883ced66213f17ed0cd64", // 小鱼
}

const (
	testProduceChatID = "oc_d3f731544c00d1d467110f5a533c5522" // 导师群测试
	produceChatID     = "oc_3cb3aa4021eee68165c78e731476d3c1" // 导师团
	testTopicChatID   = "oc_c6f4a213c3f2b35bba6e17868af7c072" // 目标话题测试
	topicChatID       = "oc_9ab398ed766e74cc68706c1919b48471" // 系统设计workshop
	mainGroupChatID   = "oc_55c292dc1b798b0d7289c29ddc89d6cc" // 训练营大群
)

// GetRemindChatWhiteIDs 发言提醒白名单
func GetRemindChatWhiteIDs() map[string]struct{} {
	return map[string]struct{}{
		produceChatID:   {},
		topicChatID:     {},
		mainGroupChatID: {},
	}
}

//GetTestEnvRemindAllowList 发言提醒测试仅发名单
func GetTestEnvRemindAllowList() map[string]struct{} {
	return map[string]struct{}{
		testProduceChatID: {},
	}
}

// GetRemindPostChatID 测试提醒消息用
func GetRemindPostChatID() string {
	if IsTestEnv() {
		return testProduceChatID
	}
	return testProduceChatID
}

// GetJoinMainGroupChatID 测试进群消息用
func GetJoinMainGroupChatID() string {
	if IsTestEnv() {
		return testProduceChatID
	}
	return mainGroupChatID
}

// GetWorkshopProduceChatID workshop发题名单
func GetWorkshopProduceChatID() string {
	if IsTestEnv() {
		return testProduceChatID
	}
	return produceChatID
}

// GetWorkshopTopicGroupChatID 获取接受workshop任务群id
func GetWorkshopTopicGroupChatID() string {
	if IsTestEnv() {
		return testTopicChatID
	}
	return topicChatID
}

func IsTestEnv() bool {
	result := false
	u, _ := user.Current()
	if strings.Contains(u.Username, "bytedance") {
		result = true
	}
	return result
}
