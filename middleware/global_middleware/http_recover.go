package global_middleware

import (
	"net/http"

	"gin_websocket/controller"
	"gin_websocket/lib/logger"

	"github.com/gin-gonic/gin"
)

func HttpRecover(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			logger.Runtime.Error(err.(error).Error())
			controller.PanicResponse(c, err.(error), http.StatusInternalServerError, "")
			c.Abort()
		}
	}()
	c.Next()
}
