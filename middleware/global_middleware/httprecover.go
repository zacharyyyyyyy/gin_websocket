package global_middleware

import (
	"gin_websocket/lib/logger"
	"github.com/gin-gonic/gin"
)

func HttpRecover(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			logger.Runtime.Error(err.(error).Error())
			return
		}
	}()
	c.Next()
}
