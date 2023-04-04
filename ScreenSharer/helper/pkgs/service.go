package helper

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var Host = "127.0.0.1"
var Port = "8080"

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
	gin.SetMode(gin.DebugMode)
	r := gin.New()
	// 基本配置
	{
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
			api.POST("/image/upload", UploadImage())
			api.POST("/text/upload", UploadText())
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

		group, exists := c.GetQuery("group")
		if !exists {
			ResponseUtils(c, http.StatusBadRequest, "缺少参数 group")
			return
		}

		filename := strconv.Itoa(int(uuid.New().ID())) + file.Filename[strings.LastIndex(file.Filename, "."):]

		msg := MessageLog{
			Username:  username,
			GroupName: group,
			Type:      LogTypeImage,
			Data:      fmt.Sprintf("http://%s:%s/image/%s", Host, Port, filename),
		}

		err = c.SaveUploadedFile(file, "./image/"+filename)
		if err != nil {
			ResponseUtils(c, http.StatusInternalServerError, "照片保存错误")
			return
		}

		err = CreateMessageLog(&msg)
		if err != nil {
			ResponseUtils(c, http.StatusInternalServerError, "数据库错误")
			return
		}

		g, ok := namespace.GetGroup(group)
		if ok {
			g.Broadcast(username, &R{
				Code: ImageCode,
				Msg:  fmt.Sprintf("( %s ) 发送了照片", username),
				Data: []MessageLog{msg},
			})
		}
		c.JSON(http.StatusOK, &R{
			Code: http.StatusOK,
			Msg:  "图片上传成功",
			Data: []MessageLog{msg},
		})
	}
}

func UploadText() gin.HandlerFunc {
	return func(c *gin.Context) {
		username, exists := c.GetQuery("username")
		if !exists {
			ResponseUtils(c, http.StatusBadRequest, "缺少参数 username")
			return
		}

		group, exists := c.GetQuery("group")
		if !exists {
			ResponseUtils(c, http.StatusBadRequest, "缺少参数 group")
			return
		}

		language, exists := c.GetQuery("language")
		if !exists {
			ResponseUtils(c, http.StatusBadRequest, "缺少参数 language")
			return
		}

		bytes, err := io.ReadAll(c.Request.Body)
		marshal, err := json.Marshal(map[string]string{"language": language, "text": string(bytes)})
		if err != nil {
			ResponseUtils(c, http.StatusInternalServerError, "文本解析错误")
			return
		}
		msg := MessageLog{
			Username:  username,
			GroupName: group,
			Type:      LogTypeText,
			Data:      string(marshal),
		}

		err = CreateMessageLog(&msg)
		if err != nil {
			ResponseUtils(c, http.StatusInternalServerError, "数据库错误")
			return
		}

		g, ok := namespace.GetGroup(group)
		if ok {
			g.Broadcast(username, &R{
				Code: ImageCode,
				Msg:  fmt.Sprintf("( %s ) 发送了文本", username),
				Data: []MessageLog{msg},
			})
		}
		c.JSON(http.StatusOK, &R{
			Code: http.StatusOK,
			Msg:  "文本上传成功",
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

		group, exists := c.GetQuery("group")
		if !exists {
			ResponseUtils(c, http.StatusBadRequest, "缺少参数 group")
			return
		}

		command, exists := c.GetQuery("command")
		if !exists {
			ResponseUtils(c, http.StatusBadRequest, "缺少参数 command")
			return
		}

		client, ok := namespace.GetClient(group, username)
		if !ok {
			ResponseUtils(c, http.StatusBadRequest, "用户客户端未上线")
			return
		}

		f := httpCmd[command]
		if f == nil {
			ResponseUtils(c, http.StatusBadRequest, fmt.Sprintf("不支持命令: %s", command))
			return
		}
		f(client, c)
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

		groupName, exists := c.GetQuery("group")
		if !exists {
			ResponseUtils(c, http.StatusBadRequest, "需要参数 username")
			return
		}

		ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			ResponseUtils(c, http.StatusInternalServerError, "协议升级错误")
			return

		}

		client := &Client{
			Username:  username,
			Conn:      ws,
			GroupName: groupName,
			CurrentID: -1,
			ReadMaxID: -1,
		}

		namespace.AddClient(client)

		go client.Listener()
	}
}

func ResponseUtils(c *gin.Context, code int, msg string) {
	c.JSON(code, R{
		Code: code,
		Msg:  msg,
	})
}
