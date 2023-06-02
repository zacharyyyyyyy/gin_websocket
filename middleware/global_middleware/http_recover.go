package global_middleware

import (
	"runtime"
	_ "unsafe"

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
			logger.Runtime.Error(err.(string))
			localContext.Abort()
			runtime.Goexit()
		}
	}()
	hookTrampoline(e)
}

//go:noinline
func hookTrampoline(e interface{}) {
}

//go:linkname gopanic runtime.gopanic
func gopanic(e interface{})
