package helper

import "fmt"

type MessageLog struct {
	ID       uint   `json:"id" gorm:"primarykey"`
	Username string `json:"username" gorm:"type:varchar(32)"`
	Type     string `json:"type" gorm:"type:varchar(32)"`
	Url      string `json:"url" gorm:"type:varchar(128)"`
}

func (m *MessageLog) String() string {
	return fmt.Sprintf(
		"{ \"id\": %d, \"username\": \"%s\", \"type\": \"%s\", \"url\":\"%s\"",
		m.ID,
		m.Username,
		m.Type,
		m.Url,
	)
}

type R struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}
