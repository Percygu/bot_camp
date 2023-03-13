package proto

import (
	"encoding/json"
	"github/rotatebot/utils"
	"github/rotatebot/utils/tools"
)

const WorkTemplate = `{
  "config": {
    "wide_screen_mode": true
  },
  "header": {
    "template": "red",
    "title": {
      "content": "「系统设计workshop」题目上新，本轮问题：{{.Topic}}",
      "tag": "plain_text"
    }
  },
  "i18n_elements": {
    "zh_cn": [
      {
        "alt": {
          "content": "",
          "tag": "plain_text"
        },
        "img_key": "img_v2_3db7cbbc-2072-4640-ae4a-a13dccdae0bg",
        "tag": "img"
      },
      {
        "tag": "div",
        "text": {
          "content": "**问题链接：** [{{.Topic}}]({{.DocLink}})",
          "tag": "lark_md"
        }
      },
      {
        "tag": "div",
        "text": {
          "content": "**本轮主要答题人：**{{.OpenIDsText}}，除了以上同学必须回答，其他同学也十分建议参与回答，增强面试能力",
          "tag": "lark_md"
        }
      },
      {
        "tag": "div",
        "text": {
          "content": "**系统设计workshop**不仅面试中真正考察能力的地方，而且有很多工作中会遇到的场景，解题前参考 [系统设计问题专项突破执行指南](https://ls8sck0zrg.feishu.cn/wiki/wikcnokuXzZeK0YWdvtpNmTe4zo) 完成答题",
          "tag": "lark_md"
        }
      }
    ]
  }
}`

const AtPersonText = `<at id=%s></at>`

type MessageDetail struct {
	MessageID string
	Topic     string
	DocLink   string
}

type PlaceholderInfo struct {
	Topic       string
	DocLink     string
	OpenIDsText string
}

// SysDesignCardModel 系统设计workshop，卡片消息
type SysDesignCardModel struct {
	Config struct {
		WideScreenMode bool `json:"wide_screen_mode"`
	} `json:"config"`
	Header struct {
		Template string `json:"template"`
		Title    struct {
			Content string `json:"content"`
			Tag     string `json:"tag"`
		} `json:"title"`
	} `json:"header"`
	I18NElements struct {
		ZhCn []struct {
			Alt struct {
				Content string `json:"content"`
				Tag     string `json:"tag"`
			} `json:"alt,omitempty"`
			ImgKey string `json:"img_key,omitempty"`
			Tag    string `json:"tag"`
			Text   struct {
				Content string `json:"content"`
				Tag     string `json:"tag"`
			} `json:"text,omitempty"`
		} `json:"zh_cn"`
	} `json:"i18n_elements"`
}

type RawText struct {
	Text string `json:"text"`
}

// PostStruct 富文本类型
type PostStruct struct {
	ZhCn *PostItem `json:"zh_cn"`
}

func NewPostStruct(title string, lines []*onePostLine) string {
	post := &PostStruct{
		ZhCn: &PostItem{
			Title:   title,
			Content: nil,
		},
	}
	post.ZhCn.Content = make([][]*GeneralPostLine, 0)
	for _, content := range lines {
		post.ZhCn.Content = append(post.ZhCn.Content, content.line)
	}
	pBytes, err := json.Marshal(post)
	if err != nil {
		return ""
	}
	return string(pBytes)
}

type PostItem struct {
	Title   string               `json:"title"`
	Content [][]*GeneralPostLine `json:"content"`
}

type GeneralPostLine struct {
	Tag       string `json:"tag"`
	Text      string `json:"text,omitempty"`
	Href      string `json:"href,omitempty"`
	UserId    string `json:"user_id,omitempty"`
	UserName  string `json:"user_name,omitempty"`
	ImageKey  string `json:"image_key,omitempty"`
	FileKey   string `json:"file_key,omitempty"`
	EmojiType string `json:"emoji_type,omitempty"`
}

type onePostLine struct {
	line []*GeneralPostLine
}

func (o *onePostLine) WriteText(text string) {
	if len(o.line) == 0 {
		o.line = make([]*GeneralPostLine, 0)
	}
	o.line = append(o.line, &GeneralPostLine{
		Tag:  "text",
		Text: text,
	})
}

func (o *onePostLine) WriteALink(href, text string) {
	if len(o.line) == 0 {
		o.line = make([]*GeneralPostLine, 0)
	}
	o.line = append(o.line, &GeneralPostLine{
		Tag:  "a",
		Text: text,
		Href: href,
	})
}

func (o *onePostLine) WriteAt(userID, userName string) {
	if len(o.line) == 0 {
		o.line = make([]*GeneralPostLine, 0)
	}
	o.line = append(o.line, &GeneralPostLine{
		Tag:      "at",
		UserId:   userID,
		UserName: userName,
	})
}

func (o *onePostLine) WriteImage(image string) {
	if len(o.line) == 0 {
		o.line = make([]*GeneralPostLine, 0)
	}
	o.line = append(o.line, &GeneralPostLine{
		Tag:      "img",
		ImageKey: image,
	})
}

// JoinGroupPost 进群消息
type JoinGroupPost struct {
	welcomeNewUserID   string
	welcomeNewUserName string
}

func NewJoinGroupPost(userID, userName string) string {
	join := &JoinGroupPost{
		welcomeNewUserID:   userID,
		welcomeNewUserName: userName,
	}
	contents := []*onePostLine{
		join.buildFirstLine(),
	}
	post := NewPostStruct("欢迎新同学加入训练营", contents)
	tools.Prettier(post)
	return post
}

func (j *JoinGroupPost) buildFirstLine() *onePostLine {
	line := &onePostLine{}
	line.WriteText("欢迎")
	line.WriteAt(j.welcomeNewUserID, j.welcomeNewUserName)
	line.WriteText("来到训练营！")
	return line
}

// OneVOneGroupPost 加入1v1专属群消息
type OneVOneGroupPost struct {
	welcomeNewUserID   string
	welcomeNewUserName string
	lines              []*onePostLine
}

func NewOneVOneGroupPost(userID, userName string) string {
	post := &OneVOneGroupPost{
		welcomeNewUserID:   userID,
		welcomeNewUserName: userName,
		lines:              make([]*onePostLine, 0),
	}
	post.buildFirstLine()
	post.buildPatchFirstLine()
	post.buildRemindLine()
	//post.buildSecondLine()
	//post.buildSpace()
	//post.buildThirdLine()
	//post.build4thLine()
	//post.build5thLine()
	postStr := NewPostStruct("专属群通知", post.lines)
	tools.Prettier(postStr)
	return postStr
}

func (o *OneVOneGroupPost) buildFirstLine() {
	line := &onePostLine{}
	line.WriteText("欢迎")
	line.WriteAt(o.welcomeNewUserID, o.welcomeNewUserName)
	line.WriteText("来到训练营！")
	o.lines = append(o.lines, line)
}

func (o *OneVOneGroupPost) buildPatchFirstLine() {
	line := &onePostLine{}
	line.WriteText("这里是你的专属1v1的群，我们将在这里为你定制计划，解答疑问，专项突破！\n")
	o.lines = append(o.lines, line)
}

func (o *OneVOneGroupPost) buildSpace() {
	line := &onePostLine{}
	line.WriteText("  ")
	o.lines = append(o.lines, line)
}

func (o *OneVOneGroupPost) buildSecondLine() {
	line := &onePostLine{}
	line.WriteText("这里有一些你需要注意的事情：")
	o.lines = append(o.lines, line)
}

func (o *OneVOneGroupPost) buildRemindLine() {
	line := &onePostLine{}
	line.WriteText("稍后会发送基本信息表，请等待填写")
	o.lines = append(o.lines, line)
}

func (o *OneVOneGroupPost) buildThirdLine() {
	line := &onePostLine{}
	line.WriteText("1. 为了跟进每周的学习进度，需要参考")
	line.WriteALink("https://ls8sck0zrg.feishu.cn/wiki/wikcnL59dpzGeCbALbG56jS699d", "周报指引")
	line.WriteText("完成周报，要注意参考里面的详细执行规范")
	o.lines = append(o.lines, line)
}

func (o *OneVOneGroupPost) build4thLine() {
	line := &onePostLine{}
	line.WriteText("2.参考")
	line.WriteALink("https://ls8sck0zrg.feishu.cn/wiki/wikcnRCnp5DQMvy7bED5M4n90NS", "个人信息模板")
	line.WriteText("复制模版，在自己知识空间下参写个人的简单信息介绍")
	o.lines = append(o.lines, line)
}

func (o *OneVOneGroupPost) build5thLine() {
	line := &onePostLine{}
	line.WriteText("3. 如果对于飞书不熟悉，那么必读")
	line.WriteALink("https://ls8sck0zrg.feishu.cn/wiki/wikcnlw1gwkzRBdlRtdfnm7JKmh", "常见Q&A")
	o.lines = append(o.lines, line)
}

// 发言提醒消息

type RemindTalkPost struct {
	day int
}

func NewRemindTalkPost(dayOff int) string {
	r := &RemindTalkPost{
		day: dayOff,
	}
	contents := []*onePostLine{
		r.friendlyRemind(),
	}
	post := NewPostStruct("群消息安静提醒", contents)
	tools.Prettier(post)
	return post
}

func (r *RemindTalkPost) friendlyRemind() *onePostLine {
	line := &onePostLine{}
	line.WriteText("gmgm, good day, it's time to discuss something!")
	line.WriteAt(utils.All, "所有人")
	return line
}

type KnowledgeSpaceFallbackPost struct {
}

func NewKnowledgeSpaceFallbackPost() string {
	r := &KnowledgeSpaceFallbackPost{}
	contents := []*onePostLine{
		r.callHelp(),
	}
	post := NewPostStruct("知识空间创建提醒", contents)
	tools.Prettier(post)
	return post
}

func (r *KnowledgeSpaceFallbackPost) callHelp() *onePostLine {
	line := &onePostLine{}
	line.WriteText("知识空间尚未创建，请")
	line.WriteAt(utils.GetBotCampConf().KnowledgeSpaceCreatorID, utils.GetBotCampConf().KnowledgeSpaceCreatorName)
	line.WriteText("协助创建知识空间哦")
	return line
}
