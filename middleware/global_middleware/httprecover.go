package global_middleware

import (
	"gin_websocket/controller"
	"github.com/gin-gonic/gin"
	"net/http"
)

func HttpRecover(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			controller.PanicResponse(c, err.(error), http.StatusInternalServerError, "")
			return
		}
	}()
	c.Next()
}
