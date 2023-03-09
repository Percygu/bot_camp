package proto

import (
	larktask "github.com/larksuite/oapi-sdk-go/v3/service/task/v1"
	"strconv"
)

const taskDemo = `
{
  "rich_summary": "系统设计workshop回答任务，[按照执行规范完成问题](https://open.feishu.cn)",
  "rich_description": "系统设计workshop回答任务创建，[在规定时间内完成作答哦](https://open.feishu.cn)",
  "due": {
    "time": "1623124318",
    "timezone": "Asia/Shanghai",
    "is_all_day": false
  },
  "origin": {
    "platform_i18n_name": "{\"zh_cn\": \"训练营机器人\", \"en_us\": \"Training Bot\"}",
    "href": {
      "url": "https://support.feishu.com/internal/foo-bar",
      "title": "mysql问题"
    }
  },
  "can_edit": true,
  "follower_ids": [
    "ou_6e75a4323ad0eef42ae349d71203969e"
  ],
  "collaborator_ids": [
    "ou_6e75a4323ad0eef42ae349d71203969e"
  ],
  "repeat_rule": "FREQ=DAILY;INTERVAL=1;BYDAY=MO,TU,WE,TH,FR"
}
`

type TaskInfo struct {
	TaskTitle      string
	TaskDocLink    string
	FinishTime     int64
	WorkerOpenID   string
	WatcherOpenIDs []string
}

func BuildTaskTemplate(taskInfo *TaskInfo) *larktask.CreateTaskReq {
	// 创建请求对象
	timeStr := strconv.FormatInt(taskInfo.FinishTime, 10)
	return larktask.NewCreateTaskReqBuilder().
		UserIdType(`open_id`).
		Task(larktask.NewTaskBuilder().
			Due(larktask.NewDueBuilder().
				Time(timeStr).
				Timezone(`Asia/Shanghai`).
				IsAllDay(true).
				Build()).
			Origin(larktask.NewOriginBuilder().
				PlatformI18nName(`{"zh_cn": "训练营机器人", "en_us": "Training Bot"}`).
				Href(larktask.NewHrefBuilder().
					Url(taskInfo.TaskDocLink).
					Title(taskInfo.TaskTitle).
					Build()).
				Build()).
			CanEdit(true).
			CollaboratorIds([]string{taskInfo.WorkerOpenID}).
			FollowerIds(taskInfo.WatcherOpenIDs).
			RichSummary(`系统设计workshop回答任务创建提醒，参考[场景问题专项突破执行指南](https://ls8sck0zrg.feishu.cn/wiki/wikcnokuXzZeK0YWdvtpNmTe4zo)作答哦`).
			RichDescription(`系统设计workshop回答任务创建，参考[场景问题专项突破执行指南](https://ls8sck0zrg.feishu.cn/wiki/wikcnokuXzZeK0YWdvtpNmTe4zo)作答哦`).
			Build()).
		Build()
}
