package app

import (
	"context"
	"fmt"
	"github/rotatebot/infra/larkinfra"
	"github/rotatebot/proto"
	"github/rotatebot/utils"

	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	larkwiki "github.com/larksuite/oapi-sdk-go/v3/service/wiki/v2"
	"github.com/sirupsen/logrus"
)

type JoinGroupHandler struct {
	ctx            context.Context
	userJoinCtx    map[string]*afterJoinInfo
	fallbackStages map[string][]fallbackStage
	msg            *larkim.P2ChatMemberUserAddedV1Data
}

type fallbackStage int64

const (
	createKnowledgeFail fallbackStage = 1
)

type afterJoinInfo struct {
	chatID string
}

func NewJoinGroupHandler() *JoinGroupHandler {
	return &JoinGroupHandler{
		userJoinCtx:    make(map[string]*afterJoinInfo),
		fallbackStages: make(map[string][]fallbackStage),
	}
}

func (j *JoinGroupHandler) Handle(ctx context.Context, msg *larkim.P2ChatMemberUserAddedV1Data) error {
	j.ctx = ctx
	j.msg = msg
	logrus.Infof("listen to msg chat name=%s", *msg.Name)
	if *j.msg.ChatId != utils.GetJoinMainGroupChatID() {
		return nil
	}
	utils.ExecAndCountFuncCtx(
		ctx,
		j.sayHello,
		j.buildNewGroup,
		j.createKnowledgeSpace,
	)
	j.fallback()
	return nil
}

// sayHello 为加入群聊用户在当前群发出欢迎语
func (j *JoinGroupHandler) sayHello() error {
	logrus.Infof("Hello!")
	msg := j.msg
	for _, user := range msg.Users {
		postInfo := proto.NewJoinGroupPost(*user.UserId.UserId, *user.Name)
		im := larkinfra.NewLarkInstance(*msg.ChatId)
		err := im.SendLarkGroup(j.ctx, postInfo, larkinfra.PostMsg)
		if err != nil {
			logrus.Errorf("sayHello err:%+v", err)
		}
	}
	return nil
}

// buildNewGroup 为消息内用户创建1v1群聊，发送欢迎语
func (j *JoinGroupHandler) buildNewGroup() error {
	logrus.Infof("Build Group")
	for _, user := range j.msg.Users {
		// 创建1v1群
		memberIDs := utils.Get1V1MentorIDs()
		memberIDs = append(memberIDs, *user.UserId.OpenId)
		// 创建请求对象
		req := larkim.NewCreateChatReqBuilder().
			UserIdType(`open_id`).
			SetBotManager(true).
			Body(larkim.NewCreateChatReqBodyBuilder().
				Name(fmt.Sprintf("%s 1v1", *user.Name)).
				Description(fmt.Sprintf("%s的专属群", *user.Name)).
				UserIdList(memberIDs).
				ChatMode(`group`).
				OwnerId(utils.Fish).
				ChatType(`private`).
				External(false).
				JoinMessageVisibility(`all_members`).
				LeaveMessageVisibility(`all_members`).
				MembershipApproval(`no_approval_required`).
				Build()).
			Build()
		resp, err := larkinfra.CreateGroup(j.ctx, req)
		if err != nil {
			logrus.Errorf("create group:%+v", err)
			continue
		}
		// 机器人触达专属群规则
		// test chat ID：oc_f308379f780158930d04fdb7efec9231
		if resp == nil || resp.Data == nil {
			logrus.Errorf("build group response empty")
			continue
		}
		chatID := *resp.Data.ChatId
		postInfo := proto.NewOneVOneGroupPost(*user.UserId.UserId, *user.Name)
		im := larkinfra.NewLarkInstance(chatID)
		err = im.SendLarkGroup(j.ctx, postInfo, larkinfra.PostMsg)
		if err != nil {
			logrus.Errorf("buildNewGroup err:%+v", err)
		}
		j.userJoinCtx[*user.UserId.OpenId] = &afterJoinInfo{chatID: chatID}
	}
	return nil
}

// createKnowledgeSpace 为消息内用户创建知识空间
func (j *JoinGroupHandler) createKnowledgeSpace() error {
	logrus.Infof("Create Knowledge Space")
	for _, user := range j.msg.Users {
		wiki := NewLarkWiki(j.ctx, *user.UserId.UserId, *user.Name, utils.MentorOpenDepartmentID)
		_, err := wiki.CreateSpaceForUser()
		if err != nil {
			logrus.Errorf("create space err:%+v", err)
			j.fallbackStages[*user.UserId.OpenId] = append(j.fallbackStages[*user.UserId.OpenId], createKnowledgeFail)
		}
	}
	return nil
}

func (j *JoinGroupHandler) fallback() {
	for openID, stages := range j.fallbackStages {
		for _, stage := range stages {
			switch stage {
			case createKnowledgeFail:
				postInfo := proto.NewKnowledgeSpaceFallbackPost()
				im := larkinfra.NewLarkInstance(j.userJoinCtx[openID].chatID)
				err := im.SendLarkGroup(j.ctx, postInfo, larkinfra.PostMsg)
				if err != nil {
					logrus.Errorf("send lark err:%+v", err)
					continue
				}
			}
		}
	}
}

type LarkWiki struct {
	// ctx 输入
	ctx context.Context
	// userID 输入
	userID string
	// userName 输入
	userName string
	// mentorDepartmentID 输入
	mentorDepartmentID string
	// space 创建space后生成
	space *larkwiki.Space
}

func NewLarkWiki(ctx context.Context, userID string,
	userName string,
	mentorDepartmentID string) *LarkWiki {
	return &LarkWiki{
		ctx:                ctx,
		userName:           userName,
		userID:             userID,
		mentorDepartmentID: mentorDepartmentID,
	}
}

func (s *LarkWiki) inviteStudentAsAdmin() error {
	req := larkwiki.NewCreateSpaceMemberReqBuilder().
		SpaceId(*s.space.SpaceId).
		Member(larkwiki.NewMemberBuilder().
			MemberType(`userid`).
			MemberId(s.userID).
			MemberRole(`admin`).
			Build()).
		Build()
	if _, err := larkinfra.CreateSpaceMember(s.ctx, req); err != nil {
		return err
	}
	return nil
}

// @return spaceID
func (s *LarkWiki) CreateSpaceForUser() (string, error) {
	err := s.createLarkSpace()
	if err != nil {
		return "", err
	}
	if err = s.inviteSpaceDepartmentAdmin(); err != nil {
		return "", err
	}

	if err = s.inviteStudentAsAdmin(); err != nil {
		return "", err
	}
	return *s.space.SpaceId, nil
}

func (s *LarkWiki) inviteSpaceDepartmentAdmin() error {
	// 创建请求对象
	req := larkwiki.NewCreateSpaceMemberReqBuilder().
		SpaceId(*s.space.SpaceId).
		Member(larkwiki.NewMemberBuilder().
			MemberType(`opendepartmentid`).
			MemberId(s.mentorDepartmentID).
			MemberRole(`admin`).
			Build()).
		Build()
	if _, err := larkinfra.CreateSpaceMember(s.ctx, req); err != nil {
		return err
	}
	return nil
}

func (s *LarkWiki) createLarkSpace() error {
	req := larkwiki.NewCreateSpaceReqBuilder().
		Space(larkwiki.NewSpaceBuilder().
			Name(s.userName + "的知识空间").
			Description(`这是你的工作学习的私人空间,导师团会陪你共同成长`).
			Build()).
		Build()

	user := larkinfra.NewLarkUser()
	userToken := user.GetMyToken()
	resp, err := larkinfra.CreateSpace(s.ctx, req, userToken)
	if err != nil {
		return err
	}
	s.space = resp.Data.Space
	return nil
}
