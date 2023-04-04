package helper

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type HttpCommand func(client *Client, ctx *gin.Context)
type WebSocketCommand func(client *Client, data map[string]string)

type CR struct {
	Command string `json:"command"`
	Data    any    `json:"data"`
}

const (
	PrevCommand  = "prev"
	NextCommand  = "next"
	CopyCommand  = "copy"
	SetupCommand = "setup"
)

var httpCmd map[string]HttpCommand
var wsCmd map[string]WebSocketCommand

func init() {
	httpCmd = make(map[string]HttpCommand, 10)
	httpCmd[NextCommand] = NextImageCommand()
	httpCmd[PrevCommand] = PrevImageCommand()
	httpCmd[CopyCommand] = CopyTextCommand()

	wsCmd = make(map[string]WebSocketCommand, 10)
	wsCmd[SetupCommand] = SetupImageCommand()
}

func NextImageCommand() HttpCommand {
	return func(c *Client, ctx *gin.Context) {
		nextId := GetNextIDMessageLog(c)

		if nextId != -1 {
			c.CurrentID = nextId
			if c.CurrentID > c.ReadMaxID {
				c.ReadMaxID = c.CurrentID
			}
		}
		_ = c.Conn.WriteJSON(&R{
			Code: CommandCode,
			Msg:  "接受命令数据",
			Data: &CR{
				Command: NextCommand,
				Data:    c.CurrentID,
			},
		})
		ctx.JSON(http.StatusOK, &R{
			Code: http.StatusOK,
			Msg:  "执行成功",
		})
	}
}

func PrevImageCommand() HttpCommand {
	return func(c *Client, ctx *gin.Context) {
		nextId := GetPrevIDMessageLog(c)
		if nextId != -1 {
			c.CurrentID = nextId
			if c.CurrentID > c.ReadMaxID {
				c.ReadMaxID = c.CurrentID
			}
		}

		_ = c.Conn.WriteJSON(&R{
			Code: CommandCode,
			Msg:  "接受命令数据",
			Data: &CR{
				Command: NextCommand,
				Data:    c.CurrentID,
			},
		})
		ctx.JSON(http.StatusOK, &R{
			Code: http.StatusOK,
			Msg:  "执行成功: prev",
		})
	}
}

func CopyTextCommand() HttpCommand {
	return func(c *Client, ctx *gin.Context) {
		log := GetShowMessageLog(c)

		if log.Type == LogTypeImage {
			_ = c.Conn.WriteJSON(&R{
				Code: MessageCode,
				Msg:  "拷贝错误: 不支持图像拷贝",
			})

			ctx.JSON(http.StatusBadRequest, &R{
				Code: http.StatusBadRequest,
				Msg:  "执行错误",
			})
			return
		}

		_ = c.Conn.WriteJSON(&R{
			Code: MessageCode,
			Msg:  fmt.Sprintf("拷贝成功: %s", log.Data[0:10]),
		})

		ctx.JSON(http.StatusOK, &R{
			Code: http.StatusOK,
			Msg:  "执行成功",
			Data: log.Data,
		})
	}
}

func SetupImageCommand() WebSocketCommand {
	return func(client *Client, data map[string]string) {
		idstr, ok := data["id"]
		id, err := strconv.Atoi(idstr)
		if !ok || err != nil {
			client.Send(client.Username, &R{
				Code: MessageCode,
				Msg:  "没有 id 字段",
			})
			return
		}
		log := GetMessageLogByID(uint(id))
		if log != nil {
			client.CurrentID = int(log.ID)
			client.Send(client.Username, &R{
				Code: MessageCode,
				Msg:  fmt.Sprintf("以切换到: %d", log.ID),
			})
			return
		}
		client.Send(client.Username, &R{
			Code: MessageCode,
			Msg:  "没有查询到该记录!",
		})
	}
}
