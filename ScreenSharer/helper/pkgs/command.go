package helper

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Command func(client *Client, ctx *gin.Context)

type CR struct {
	Command string `json:"command"`
	Data    any    `json:"data"`
}

const (
	PrevCommand = "prev"
	NextCommand = "next"
	CopyCommand = "copy"
)

var cmd map[string]Command

func init() {
	cmd = make(map[string]Command, 10)
	cmd[NextCommand] = NextImageCommand()
	cmd[PrevCommand] = PrevImageCommand()
	cmd[CopyCommand] = CopyTextCommand()
}

func NextImageCommand() Command {
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

func PrevImageCommand() Command {
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

func CopyTextCommand() Command {
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
