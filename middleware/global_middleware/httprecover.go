package global_middleware

import (
	"gin_websocket/controller"
	"gin_websocket/lib/logger"
	"github.com/gin-gonic/gin"
	"net/http"
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
