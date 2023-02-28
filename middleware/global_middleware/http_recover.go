package global_middleware

import (
	"net/http"
	"runtime"

	"gin_websocket/controller"
	"gin_websocket/lib/logger"

	"github.com/brahma-adshonor/gohook"
	"github.com/gin-gonic/gin"
)

var localContext *gin.Context

func HttpRecover(c *gin.Context) {
	gohook.Hook(gopanic, hookRecover, hookTrampoline)
	localContext = c
	c.Next()
}

func hookRecover(e interface{}) {
	defer func() {
		if err := recover(); err != nil {
			logger.Runtime.Error(err.(error).Error())
			controller.PanicResponse(localContext, err.(error), http.StatusInternalServerError, "")
			localContext.Abort()
			runtime.Goexit()
		}
	}()
	hookTrampoline(e)
}

func hookTrampoline(e interface{}) {
}

//go:linkname gopanic runtime.gopanic
func gopanic(e interface{})
