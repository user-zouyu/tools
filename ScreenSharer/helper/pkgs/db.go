package helper

import (
	"errors"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func GetDB() *gorm.DB {
	return db
}

func InitDB() error {
	open, err := gorm.Open(sqlite.Open("./data/db.db"), &gorm.Config{})
	if err != nil {
		return errors.New("数据库连接错误")
	}
	db = open
	err = db.AutoMigrate(&MessageLog{})
	if err != nil {
		return errors.New("数据库初始化错误")
	}
	return nil
}
