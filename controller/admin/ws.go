package admin

import (
	"context"
	"errors"
	"fmt"
	"gin_websocket/controller"
	ws "gin_websocket/service/websocket"
	"github.com/gin-gonic/gin"
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
		if errors.Is(err, ws.CloseErr) {
			fmt.Println(err.Error())
			break
		}
		fmt.Println("loop")
	}
	fmt.Println("connect close")
}

func ServiceLink(c *gin.Context) {

}

func Ping(c *gin.Context) {

}
