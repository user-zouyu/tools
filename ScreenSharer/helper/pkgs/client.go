package helper

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
)

type Client struct {
	Username  string
	Conn      *websocket.Conn
	Group     *Group
	NameSpace *NameSpace
	GroupName string
	CurrentID int // 当前展示的消息索引
	ReadMaxID int // 已读消息的最大索引
}

func (c *Client) Broadcast(data any) {
	log.Printf("广播消息(username: %s): %v", c.Username, data)
	c.Group.Broadcast(c.Username, data)
}

func (c *Client) Send(from string, data any) {
	err := c.Conn.WriteJSON(data)
	if err != nil {
		log.Printf("发送错误( %s->%s ), data: %v", from, c.Username, data)
	}
}

func (c *Client) SendTo(to string, data any) {
	client, ok := c.Group.GetClient(to)
	if ok {
		err := client.Conn.WriteJSON(data)
		if err != nil {
			log.Printf("发送错误( %s->%s ), data: %v", c.Username, to, data)
		}
	} else {
		c.Send(c.Username, &R{
			Code: MessageCode,
			Msg:  fmt.Sprintf("( %s ) 没有上线", to),
		})
	}
}

func (c *Client) Listener() {
	{
		c.Broadcast(&R{
			Code: MessageCode,
			Msg:  fmt.Sprintf("( %s ) 上线了", c.Username),
		})

		list := GetMessageLogs(c)
		if len(list) > 0 {
			c.CurrentID = int(list[len(list)-1].ID)
			if c.CurrentID > c.ReadMaxID {
				c.ReadMaxID = c.CurrentID
			}
		}
		_ = c.Conn.WriteJSON(&R{
			Code: BatchImageCode,
			Msg:  "接受历史数据",
			Data: &BatchImageResponse{
				CurrentID: c.CurrentID,
				ReadMaxID: c.ReadMaxID,
				List:      list,
			},
		})
	}
	for {
		_, bytes, err := c.Conn.ReadMessage()
		if err != nil {
			c.NameSpace.DelClient(c)
			return
		}
		var data map[string]string
		err = json.Unmarshal(bytes, &data)
		if err != nil {
			log.Printf("格式解析错误:data: %s", string(bytes))
			c.Send(c.Username, &R{
				Code: MessageCode,
				Msg:  "格式解析错误",
			})
			continue
		}
		cmd, ok := data["command"]
		if !ok {
			c.Send(c.Username, &R{
				Code: MessageCode,
				Msg:  "缺少字段 ( command )",
			})
			continue
		}

		f, ok := wsCmd[cmd]
		if !ok {
			log.Printf("不支持该命令: (%s) cmd: %s, data: %v", c.Username, cmd, data)
			c.Send(c.Username, &R{
				Code: MessageCode,
				Msg:  fmt.Sprintf("不支持该命令：%s", cmd),
			})
			continue
		}

		log.Printf("处理命令: (%s) cmd: %s, data: %v", c.Username, cmd, data)
		f(c, data)
	}
}
