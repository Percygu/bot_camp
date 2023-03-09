package task

import (
	"context"
	"github/rotatebot/dao"
	"github/rotatebot/model"
	"gorm.io/gorm/clause"
)

func GetUnsentTasks(ctx context.Context, lastProcessID uint, limit int) ([]*model.NoticeModel, error) {
	var dest = make([]*model.NoticeModel, 0)
	err := dao.DB.Model(&model.NoticeModel{}).Where("stage = ?", model.StageInit).Where("id > ? ", lastProcessID).
		Limit(limit).Find(&dest).Error
	if err != nil {
		return nil, err
	}
	return dest, nil
}

func UpdateTasks(ctx context.Context, id uint, info *model.NoticeModel) error {
	return dao.DB.Model(&model.NoticeModel{}).Where("id = ?", id).Updates(info).Error
}

func CreateTask(ctx context.Context, info *model.NoticeModel) error {
	err := dao.DB.Model(&model.NoticeModel{}).Clauses(
		clause.OnConflict{
			Columns:   []clause.Column{{Name: "message_id"}},
			DoNothing: true},
	).Create(info).Error
	if err != nil {
		return err
	}
	return nil
}
