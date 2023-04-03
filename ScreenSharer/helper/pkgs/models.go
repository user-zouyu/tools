package helper

import "fmt"

type MessageLog struct {
	ID        uint   `json:"id" gorm:"primaryKey"`
	GroupName string `json:"group" gorm:"type:varchar(16); index"`
	Username  string `json:"username" gorm:"type:varchar(32)"`
	Type      string `json:"type" gorm:"type:varchar(32)"`
	Data      string `json:"url" gorm:"type:varchar(2048)"`
}

func (m *MessageLog) String() string {
	return fmt.Sprintf(
		"{ \"id\": %d, \"username\": \"%s\", \"type\": \"%s\", \"url\":\"%s\"",
		m.ID,
		m.Username,
		m.Type,
		m.Data,
	)
}

type R struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

type BatchImageResponse struct {
	CurrentID int          `json:"currentID"` // 当前展示的消息索引
	ShowMinID int          `json:"showMinID"` // 显示的最小索引
	ReadMaxID int          `json:"readMaxID"` // 已读消息的最大索引
	List      []MessageLog `json:"list"`
}
