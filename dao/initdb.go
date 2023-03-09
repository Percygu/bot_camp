package dao

import (
	"github/rotatebot/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	var err error
	DB, err = gorm.Open(sqlite.Open("practice.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	//DB = DB.Debug()

	err = DB.AutoMigrate(&model.NoticeModel{})

	if err != nil {
		panic(err)
	}
}
