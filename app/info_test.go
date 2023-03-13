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

// TestListChatID 列出机器人在的聊天群
func TestListChatID1(t *testing.T) {

	larkinfra.InitClient("cli_a31ac04fbef8900b", "83BUHQBSPLm136PftJXhcBBrPUGyqv1I")

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
	larkinfra.InitClient("cli_a4880a71a638d00e", "Je8ygYnhXu2y4BlIUCPMbbojRdWtyyaX")
	for _, m := range larkinfra.FetchAllGroupMembers(context.Background(), "oc_a054cf305b885f62b3fb76ba3fecbef0", utils.OpenIDType) {
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
	larkinfra.InitClient("cli_a4880a71a638d00e", "Je8ygYnhXu2y4BlIUCPMbbojRdWtyyaX")
	//larkinfra.CreateGroup(context.Background(), "ou_844198543d625613bb5f9a5e4f2366c2", "测试")
	userName := "鹏哥"
	memberIDs := []string{
		"ou_0f065b034d6ea3a4158f389a4fae3b1e",
		"ou_697abf63fa74f346b67f68f93b2bc81b",
		"ou_eb016ce9f4bccdc5bb02413912ce6ff3",
		"ou_89455588b44f5210f698b36b7c7b6141",
		"ou_6a11169c46840d664140c8d697e552d2",
		"ou_af5d943271517308c3c5d0f7f688276f",
		"ou_acc578572eed6500f8f30dd6da2a6d88",
	}
	req := larkim.NewCreateChatReqBuilder().
		UserIdType(`open_id`).
		SetBotManager(true).
		Body(larkim.NewCreateChatReqBodyBuilder().
			Name(fmt.Sprintf("%s 1v1", userName)).
			Description(fmt.Sprintf("%s的专属群", userName)).
			UserIdList(memberIDs).
			ChatMode(`group`).
			OwnerId(utils.PengGe).
			ChatType(`private`).
			External(false).
			JoinMessageVisibility(`all_members`).
			LeaveMessageVisibility(`all_members`).
			MembershipApproval(`no_approval_required`).
			Build()).
		Build()
	resp, err := larkinfra.CreateGroup(context.Background(), req)
	if err != nil {
		logrus.Errorf("create group:%+v", err)
	}
	t.Logf("resp === %v", resp)
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
