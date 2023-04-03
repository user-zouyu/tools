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

func GetMessageLogs(c *Client) []MessageLog {
	var r []MessageLog
	err := db.Where("group_name = ?", c.GroupName).Order("id desc").Find(&r).Error
	if err != nil {
		return []MessageLog{}
	}
	return reverse(r)
}

func CreateMessageLog(log *MessageLog) error {
	return db.Create(log).Error
}

func GetNextIDMessageLog(c *Client) int {
	var r MessageLog
	err := db.Where("group_name = ? and id > ?", c.GroupName, c.CurrentID).First(&r).Error
	if err != nil {
		return -1
	}
	return int(r.ID)
}

func GetPrevIDMessageLog(c *Client) int {
	var r MessageLog
	err := db.Where("group_name = ? and id < ?", c.GroupName, c.CurrentID).Last(&r).Error
	if err != nil {
		return -1
	}
	return int(r.ID)
}

func GetShowMessageLog(c *Client) *MessageLog {
	var r MessageLog
	err := db.Where("group_name = ? and id = ?", c.GroupName, c.CurrentID).Last(&r).Error
	if err != nil {
		return nil
	}
	return &r
}

func reverse(list []MessageLog) []MessageLog {
	l := len(list)
	for i := 0; i < l/2; i++ {
		temp := list[i]
		list[i] = list[l-i-1]
		list[l-i-1] = temp
	}
	return list
}
