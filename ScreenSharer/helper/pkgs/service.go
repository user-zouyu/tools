package helper

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
)

var Host = "127.0.0.1"
var Port = "8080"

var wsmap sync.Map

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func Server() *gin.Engine {
	if h := os.Getenv("HOST"); h != "" {
		Host = h
	}
	if p := os.Getenv("PORT"); p != "" {
		Port = p
	}
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	// 基本配置
	{
		gin.SetMode(gin.DebugMode)
		r.Use(gin.Logger())
		r.Use(gin.CustomRecovery(ExceptionHandler()))
		r.NoRoute(NoRouteHandler())
		r.NoMethod(NoMethodHandler())
		r.GET("/ping", Ping())
	}

	// 业务配置
	{
		r.Static("/home", "./html")
		r.GET("/image/:filename", GetImage())

		api := r.Group("/api")
		{
			api.POST("/upload", UploadImage())
			api.GET("/connect", ConnectSession())
			api.GET("/command", CommandService())
		}

	}

	return r
}

func Ping() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, R{
			Code: 200,
			Msg:  "Ping Success",
		})
	}
}

func NoRouteHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		c.JSON(404, &R{
			Code: 404,
			Msg:  "路径不存在",
			Data: path,
		})
	}
}

func NoMethodHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		c.JSON(404, &R{
			Code: 404,
			Msg:  "方法不存在",
			Data: method,
		})
	}
}

func ExceptionHandler() gin.RecoveryFunc {
	return func(c *gin.Context, err any) {
		c.JSON(http.StatusInternalServerError, &R{
			Code: 500,
			Msg:  "服务器处理错误",
			Data: err,
		})
	}
}

func GetImage() gin.HandlerFunc {
	return func(c *gin.Context) {
		filename := c.Param("filename")
		c.Header("content-type", "image/png")

		if _, err := os.Stat("./image/" + filename); err != nil {
			c.File("./image/notfound.png")
			return
		}
		c.File("./image/" + filename)
	}
}

func UploadImage() gin.HandlerFunc {
	return func(c *gin.Context) {
		file, err := c.FormFile("file")
		if err != nil {
			ResponseUtils(c, http.StatusBadRequest, "缺少参数 file")
			return
		}
		username, exists := c.GetQuery("username")
		if !exists {
			ResponseUtils(c, http.StatusBadRequest, "缺少参数 username")
			return
		}

		filename := strconv.Itoa(int(uuid.New().ID())) + file.Filename[strings.LastIndex(file.Filename, "."):]

		msg := MessageLog{
			Username: username,
			Type:     "image",
			Url:      fmt.Sprintf("http://%s:%s/image/%s", Host, Port, filename),
		}

		err = GetDB().Create(&msg).Error
		if err != nil {
			ResponseUtils(c, http.StatusInternalServerError, "数据库错误")
			return
		}

		err = c.SaveUploadedFile(file, "./image/"+filename)
		if err != nil {
			ResponseUtils(c, http.StatusInternalServerError, "照片保存错误")
			return
		}

		broadcast(&R{
			Code: ImageCode,
			Msg:  fmt.Sprintf("( %s ) 发送了照片", username),
			Data: []MessageLog{msg},
		})

		c.JSON(http.StatusOK, &R{
			Code: http.StatusOK,
			Msg:  "图片上传成功",
			Data: []MessageLog{msg},
		})
	}
}

func CommandService() gin.HandlerFunc {
	return func(c *gin.Context) {
		username, exists := c.GetQuery("username")
		if !exists {
			ResponseUtils(c, http.StatusBadRequest, "缺少参数 username")
			return
		}

		command, exists := c.GetQuery("command")
		if !exists {
			ResponseUtils(c, http.StatusBadRequest, "缺少参数 command")
			return
		}

		send(username, &R{
			Code: CommandCode,
			Msg:  "执行命令",
			Data: command,
		})

		c.JSON(http.StatusOK, &R{
			Code: http.StatusOK,
			Msg:  "命令发送成功",
			Data: command,
		})
	}
}

// ws://127.0.0.1:8080/api/connect?username=zouyu
func ConnectSession() gin.HandlerFunc {
	return func(c *gin.Context) {
		username, exists := c.GetQuery("username")
		if !exists {
			ResponseUtils(c, http.StatusBadRequest, "需要参数 username")
			return
		}

		ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			ResponseUtils(c, http.StatusInternalServerError, "协议升级错误")
			return
		}
		wsmap.Store(username, ws)

		go listener(username, ws)
	}
}

func listener(username string, ws *websocket.Conn) {
	{
		broadcast(&R{
			Code: MessageCode,
			Msg:  fmt.Sprintf("( %s ) 上线了", username),
		})

		var list []MessageLog
		err := GetDB().Model(&MessageLog{}).Find(&list).Error
		if err == nil {
			_ = ws.WriteJSON(&R{
				Code: HistoryImageCode,
				Msg:  "接受历史数据",
				Data: list,
			})
		} else {
			log.Printf("( %s ) 历史输出查询错误\n", username)
		}
	}

	for {
		messageType, bytes, err := ws.ReadMessage()
		if err != nil {
			log.Printf("( %s ) 接收消息错误: %v\n", username, err)
			wsmap.Delete(username)
			broadcast(&R{
				Code: MessageCode,
				Msg:  fmt.Sprintf("( %s ) 下线了", username),
			})
			_ = ws.Close()
			break
		}
		log.Printf("( %s ) code: %d, msg: %s", username, messageType, string(bytes))
	}
}

func broadcast(data any) {
	log.Printf("广播消息: %v", data)
	wsmap.Range(func(u, value any) bool {
		ws := value.(*websocket.Conn)
		err := ws.WriteJSON(data)
		if err != nil {
			log.Printf("发送错误(username: %s): %v\n", u, err)
		}
		return true
	})
}

func send(username string, data any) {
	if value, ok := wsmap.Load(username); ok {
		ws, _ := value.(*websocket.Conn)
		err := ws.WriteJSON(data)
		if err != nil {
			log.Printf("( %s ) 发送数据错误, data: %v", username, data)
		}
	}
}

func ResponseUtils(c *gin.Context, code int, msg string) {
	c.JSON(code, R{
		Code: code,
		Msg:  msg,
	})
}
