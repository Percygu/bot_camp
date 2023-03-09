package app

import (
	"context"
	"fmt"
	"github/rotatebot/infra"
	"github/rotatebot/infra/larkinfra"
	"github/rotatebot/proto"
	"github/rotatebot/utils"
	"log"
	"math/rand"
	"testing"
	"time"

	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	"github.com/sirupsen/logrus"
)

// TestListChatID 列出机器人在的聊天群
func TestListChatID(t *testing.T) {
	ctx := context.Background()
	req := larkim.NewListChatReqBuilder().UserIdType("open_id").PageSize(100).Build()
	ans, err := larkinfra.Client.Im.Chat.List(ctx, req)
	if err != nil {
		panic(err)
	}
	if !ans.Success() {
		fmt.Println(ans.Code, ans.Msg, ans.RequestId())
		panic(ans.Msg)
	}
	log.Printf("%s\n", larkcore.Prettify(ans))
}

func TestListOpenIDs(t *testing.T) {
	for _, m := range larkinfra.FetchAllGroupMembers(context.Background(), "oc_b0ef9fad3ee6828ecfebd4c8f0cfa249", utils.OpenIDType) {
		fmt.Println(larkcore.Prettify(m))
	}
}

func TestGetDoc(t *testing.T) {
	ctx := context.Background()
	getMessageDetails(ctx, 500)
}

func TestPick(t *testing.T) {
	for i := 0; i < 10; i += 1 {
		rand.Seed(time.Now().Unix())
		choiceNumber := rand.Intn(3) + 1
		picker := newPickInstance(choiceNumber)
		members := larkinfra.FetchAllGroupMembers(context.Background(), utils.GetWorkshopTopicGroupChatID(), utils.OpenIDType)
		t.Log(len(members))
		picker.PickOpenIDs(members, 10)
		t.Log(picker.members)
	}
}

func TestCreateGroup(t *testing.T) {
	//larkinfra.CreateGroup(context.Background(), "ou_844198543d625613bb5f9a5e4f2366c2", "测试")
}

func TestUpdateVis(t *testing.T) {
	larkinfra.UpdateVisibility(context.Background(), utils.Get1V1MentorIDs())
}

func TestWelcome1(t *testing.T) {
	postInfo := proto.NewJoinGroupPost(utils.PengGe, "鹏哥")
	im := larkinfra.NewLarkInstance(utils.GetJoinMainGroupChatID())
	_ = im.SendLarkGroup(context.Background(), postInfo, larkinfra.PostMsg)
}

func TestWelcome2(t *testing.T) {
	postInfo := proto.NewOneVOneGroupPost(utils.PengGe, "鹏哥")
	im := larkinfra.NewLarkInstance(utils.GetJoinMainGroupChatID())
	_ = im.SendLarkGroup(context.Background(), postInfo, larkinfra.PostMsg)
}

func TestRemindChat(t *testing.T) {
	infra.StartCronjob()
	RemindTalk(context.Background())
	time.Sleep(time.Minute * 2)
}

func TestRemindPost(t *testing.T) {
	ctx := context.Background()
	msg := proto.NewRemindTalkPost(dayOff)
	im := larkinfra.NewLarkInstance(utils.GetRemindPostChatID())
	_ = im.SendLarkGroup(ctx, msg, larkinfra.PostMsg)
}

func TestCreateWiki(t *testing.T) {
	wiki := NewLarkWiki(context.Background(), utils.PengGe, "测试狼", utils.MentorOpenDepartmentID)
	_, err := wiki.CreateSpaceForUser()
	if err != nil {
		logrus.Errorf("create space err:%+v", err)
	}
}

func GetFirstDayOfMonth(d time.Time) time.Time {
	d = d.AddDate(0, 0, -d.Day()+1)
	return GetZeroTime(d)
}
func GetZeroTime(d time.Time) time.Time {
	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, d.Location())
}
func TestTime(t *testing.T) {

	fmt.Println(GetFirstDayOfMonth(time.Now().AddDate(0, -4, 0)))
}
