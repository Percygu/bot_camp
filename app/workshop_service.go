package app

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github/rotatebot/dao/task"
	"github/rotatebot/infra/larkinfra"
	"github/rotatebot/model"
	"github/rotatebot/proto"
	"github/rotatebot/utils"
	"math/rand"
	"strconv"
	"strings"
	"text/template"
	"time"

	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	log "github.com/sirupsen/logrus"
)

const (
	taskSpaceID = "7174087082065805314" // 系统设计问题知识空间id

)

const (
	defaultFetchTime = 600 // 600s 内消息转达
)

type pickInstance struct {
	members    []string
	callNumber int
}

func newPickInstance(count int) *pickInstance {
	return &pickInstance{members: make([]string, 0), callNumber: count}
}

func (p *pickInstance) GetMembers() []string {
	return p.members
}

func (p *pickInstance) GetMembersText() string {
	ans := ""
	for _, member := range p.members {
		item := fmt.Sprintf(proto.AtPersonText, member)
		ans += item
	}
	return ans
}

func (p *pickInstance) PickOpenIDs(members []*larkim.ListMember, times int) {
	if p.callNumber > len(members) {
		p.callNumber = len(members)
	}
	if len(members) == 0 || times <= 0 || len(p.members) > p.callNumber {
		log.Infof("Finish Pick,Pick Source Length=%d,Total Picked Member Length=%d,Max Number=%d", len(members), len(p.members), p.callNumber)
		return
	}
	pickedMember := *members[rand.Intn(len(members))]
	id := *pickedMember.MemberId
	whiteIDs := utils.GetWhiteOpenIDs()
	// 跳过白名单和去重的
	if utils.Contains(whiteIDs, id) || utils.Contains(p.members, id) {
		p.PickOpenIDs(members, times-1)
		return
	}
	log.Infof("pick memeber id = %s,name is %s", id, *pickedMember.Name)
	p.members = append(p.members, id)
	if len(p.members) <= p.callNumber {
		p.PickOpenIDs(members, times)
	}
}

// getMessageDetails message format: @bot https://...
func getMessageDetails(ctx context.Context, before int64) ([]*proto.MessageDetail, bool) {
	links := make([]*proto.MessageDetail, 0)
	start := strconv.FormatInt(time.Now().Unix()-before, 10)
	end := strconv.FormatInt(time.Now().Unix(), 10)
	req := larkim.NewListMessageReqBuilder().ContainerIdType("chat").
		ContainerId(utils.GetWorkshopProduceChatID()).
		StartTime(start).
		EndTime(end).
		PageSize(20).
		Build()
	msgs, err := larkinfra.Client.Im.Message.List(ctx, req)
	if err != nil {
		log.Errorf("%+v", err)
		return nil, false
	}
	if msgs.Data == nil {
		return nil, false
	}
	for _, msg := range msgs.Data.Items {
		//vb, _ := json.Marshal(msg)
		//log.Printf("%s", vb)
		if len(msg.Mentions) != 1 {
			continue
		}
		// mention message
		mention := msg.Mentions[0]
		if *mention.Id != utils.BotOpenID {
			log.Println("not calling me ")
			continue
		}
		// doc link
		docLink, ok := extractRawText(msg)
		if !ok {
			continue
		}
		docLink.MessageID = *msg.MessageId
		links = append(links, docLink)
		//log.Infof("got message=%v", docLink)
	}
	if len(links) > 0 {
		return links, true
	}
	return nil, false
}

func extractRawText(msg *larkim.Message) (*proto.MessageDetail, bool) {
	rawText := &proto.RawText{}
	err := json.Unmarshal([]byte(*msg.Body.Content), rawText)
	if err != nil {
		panic(err)
	}
	mentionMe := msg.Mentions[0]
	//log.Infof("raw text is %s", rawText.Text)
	rawText.Text = strings.Trim(rawText.Text, "<p>")
	rawText.Text = strings.Trim(rawText.Text, "</p>")
	messages := strings.TrimSpace(strings.TrimLeft(rawText.Text, *mentionMe.Key))
	if !strings.Contains(messages, "https") {
		//log.Printf("wrong message,no https,msg=%s", messages)
		return nil, false
	}

	// xxx https://
	chunks := strings.Fields(messages)
	if len(chunks) < 2 {
		log.Printf("wrong format,chunks=%s", chunks)
		return nil, false
	}

	detail := &proto.MessageDetail{
		Topic:   strings.TrimSpace(chunks[0]),
		DocLink: strings.TrimSpace(chunks[1]),
	}
	return detail, true
}

func prepareTaskCard(info *proto.PlaceholderInfo) (string, error) {
	// todo： build sub doc and call the writer
	// send larkinfra message
	tmp, err := template.New("larkcard").Parse(proto.WorkTemplate)
	if err != nil {
		log.Errorf("%+v", err)
		return "", err
	}
	var output []byte
	buf := bytes.NewBuffer(output)
	err = tmp.Execute(buf, info)
	if err != nil {
		log.Errorf("%+v", err)
		return "", err
	}
	holder := &proto.SysDesignCardModel{}
	err = json.Unmarshal(buf.Bytes(), holder)
	if err != nil {
		log.Errorf("%+v", err)
		return "", err
	}
	hbytes, err := json.Marshal(holder)
	if err != nil {
		log.Errorf("%+v", err)
		return "", err
	}
	return string(hbytes), nil
}

type taskEngine struct {
	taskTitle   string
	taskDocLink string
	finishTime  int64
	workerID    string
	watcherIDs  []string
}

func NewTaskEngine(
	title string,
	docLink string,
	workerID string,
	watcherIDs []string,
) *taskEngine {
	endTime := time.Now().AddDate(0, 0, 7).Unix()
	return &taskEngine{
		taskTitle:   title,
		taskDocLink: docLink,
		finishTime:  endTime,
		workerID:    workerID,
		watcherIDs:  watcherIDs,
	}
}

func (t *taskEngine) createWorkTask(ctx context.Context) error {
	taskInfo := &proto.TaskInfo{
		TaskTitle:      t.taskTitle,
		TaskDocLink:    t.taskDocLink,
		FinishTime:     t.finishTime,
		WorkerOpenID:   t.workerID,
		WatcherOpenIDs: t.watcherIDs,
	}
	req := proto.BuildTaskTemplate(taskInfo)
	_, err := larkinfra.CreateTask(ctx, req)
	if err != nil {
		log.Error(err)
	}
	log.Infof("Create Work")
	return err
}

func saveToDB(ctx context.Context) {
	for {
		// 获取前段时间的全部消息
		details, has := getMessageDetails(ctx, defaultFetchTime)
		if !has {
			time.Sleep(time.Second)
			continue
		}
		for _, detail := range details {
			record := &model.NoticeModel{
				Topic:     detail.Topic,
				DocLink:   detail.DocLink,
				MessageID: detail.MessageID,
				Stage:     model.StageInit,
			}
			err := task.CreateTask(ctx, record)
			if err != nil {
				log.Errorf("saveDB%+v", err)
				continue
			}
			//log.Infof("save record success=%+v", record)
		}
	}
}

func executeFromDB(ctx context.Context) {
	sendInstance := larkinfra.NewLarkInstance(utils.GetWorkshopTopicGroupChatID())
	var lastProcessID uint = 0
	for {
		time.Sleep(time.Second)
		// todo:不保证必然触达
		//log.Infof("execute on ID=%d", lastProcessID)
		records, err := task.GetUnsentTasks(ctx, lastProcessID, 10)
		if err != nil || len(records) == 0 {
			if err != nil {
				log.Warnf("ERROR=%+v", err)
			}
			continue
		}
		for _, record := range records {
			// 随机选6个人
			choiceNumber := 6
			picker := newPickInstance(choiceNumber)
			picker.PickOpenIDs(larkinfra.FetchAllGroupMembers(ctx, utils.GetWorkshopTopicGroupChatID(), utils.OpenIDType), 10)
			if len(picker.GetMembers()) == 0 {
				log.Errorf("NO ANSWER PEOPLES!")
			}
			// 准备@的文案
			text := picker.GetMembersText()
			info := &proto.PlaceholderInfo{
				Topic:       record.Topic,
				DocLink:     record.DocLink,
				OpenIDsText: text,
			}
			// 准备lark的card消息
			flow, iErr := prepareTaskCard(info)
			if iErr != nil {
				log.Errorf("build task %v", iErr)
				continue
			}
			// 优先更新记录，避免重复触达
			iErr = task.UpdateTasks(ctx, record.ID, &model.NoticeModel{
				Stage: model.StageSent,
			})
			if iErr != nil {
				log.Errorf("send task %v", iErr)
				continue
			}
			// lark消息触达卡片消息
			err = sendInstance.SendLarkGroup(ctx, flow, larkinfra.InteractiveMsg)
			if err != nil {
				log.Errorf("send lark group err:%+v", err)
				continue
			}
			// lark任务创建
			for _, worker := range picker.GetMembers() {
				engine := NewTaskEngine(record.Topic, record.DocLink, worker, utils.GetTaskFocusIDs())
				_ = engine.createWorkTask(ctx)
			}
			//拉预讨论群
			createTopicPrepareGroup(ctx, record, picker.GetMembers())
		}
		lastProcessID = records[len(records)-1].ID
	}
}

func createTopicPrepareGroup(ctx context.Context, notice *model.NoticeModel, memberOpenIDs []string) {
	memberIDs := utils.GetWorkshopPrepareChatMentors()
	memberIDs = append(memberIDs, memberOpenIDs...)
	// 创建请求对象
	req := larkim.NewCreateChatReqBuilder().
		UserIdType(`open_id`).
		SetBotManager(true).
		Body(larkim.NewCreateChatReqBodyBuilder().
			Name(fmt.Sprintf("预讨论-%s", notice.Topic)).
			Description(fmt.Sprintf("问题链接：%s", notice.DocLink)).
			UserIdList(memberIDs).
			ChatMode(`group`).
			ChatType(`private`).
			External(false).
			JoinMessageVisibility(`all_members`).
			LeaveMessageVisibility(`all_members`).
			MembershipApproval(`no_approval_required`).
			Build()).
		Build()
	_, err := larkinfra.CreateGroup(ctx, req)
	if err != nil {
		log.Errorf("create group:%+v", err)
	}

}

func TriggerWorkshop(ctx context.Context) {
	go saveToDB(ctx)
	go executeFromDB(ctx)
}
