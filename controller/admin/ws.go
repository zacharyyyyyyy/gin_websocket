package admin

import (
	"fmt"
	"gin_websocket/controller"
	"gin_websocket/lib/logger"
	"gin_websocket/lib/redis"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
)

func Link(c *gin.Context) {
	if !websocket.IsWebSocketUpgrade(c.Request) {
		c.JSON(http.StatusOK, gin.H{"message": "非websocket", "status": 0})
	}
	var upGrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	//升级get请求为webSocket协议
	fmt.Println(c.Request.Header.Get("Sec-Websocket-Key"))
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	fmt.Println("connect success")
	defer ws.Close()
	fmt.Println()
	for {
		//读取ws中的数据
		mt, message, err := ws.ReadMessage()
		if err != nil {
			fmt.Println("connect close")
			break
		}

		//写入ws数据

		err = ws.WriteMessage(mt, message)
		if err != nil {
			break
		}
	}

}

func Ping(c *gin.Context) {
	var code int
	var str []string
	str, err := redis.RedisDb.SMembers("test")
	if err != nil {
		logger.Api.Error(err.Error())
		code = http.StatusInternalServerError
	} else {
		code = http.StatusOK
	}
	baseController := controller.ResponseStruct{Data: str, Code: code, C: c}
	baseController.JsonResponse()

}
