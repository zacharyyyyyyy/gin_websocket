package admin

import (
	"context"
	"fmt"
	"net/http"

	"gin_websocket/controller"
	"gin_websocket/lib/logger"
	"gin_websocket/lib/redis"
	ws "gin_websocket/service/websocket"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func Link(c *gin.Context) {
	ctx, _ := context.WithCancel(context.Background())
	userClient, err := ws.NewUserClient(ctx, c)
	if err != nil {
		controller.PanicResponse(c, err)
	}
	fmt.Println("connect success")
	for {
		err := userClient.Receive()
		if err != nil {
			fmt.Println("connect close", err.Error())
			break
		}

		fmt.Println("loop")
	}
	fmt.Println("connect close")
}

func ServiceLink(c *gin.Context) {
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
	wsH, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	fmt.Println("connect success", wsH)
	defer wsH.Close()
	fmt.Println()
	for {
		//读取ws中的数据
		mt, message, err := wsH.ReadMessage()
		if err != nil {
			fmt.Println("connect close", mt, message, err.Error())
			break
		}

		//写入ws数据

		err = wsH.WriteMessage(mt, message)
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
