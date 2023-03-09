package model

import (
	"gorm.io/gorm"
)

type NoticeStage int

const (
	StageInit NoticeStage = iota
	StageSent
)

type NoticeModel struct {
	gorm.Model
	Topic     string
	DocLink   string
	MessageID string      `gorm:"column:message_id;unique"`
	Stage     NoticeStage `gorm:"column:stage"`
}
