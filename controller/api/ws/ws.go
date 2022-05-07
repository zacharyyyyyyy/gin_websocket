package ws

import (
	"context"
	"errors"
	"net/http"

	"gin_websocket/controller"
	ws "gin_websocket/service/websocket"
	"github.com/gin-gonic/gin"
)

func Link(c *gin.Context) {
	ctx, _ := context.WithCancel(context.Background())
	userClient, err := ws.NewUserClient(ctx, c.Request, c.Writer, c.ClientIP())
	if err != nil {
		controller.PanicResponse(c, err, http.StatusInternalServerError, "")
		return
	}
	for {
		err := userClient.Receive()
		if errors.Is(err, ws.CloseErr) {
			break
		}
	}
}
