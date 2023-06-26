package global_middleware

import (
	"runtime"
	_ "unsafe"

	"gin_websocket/lib/logger"

	"github.com/brahma-adshonor/gohook"
	"github.com/gin-gonic/gin"
)

func HttpRecover(c *gin.Context) {
	gohook.Hook(gopanic, hookRecover, hookTrampoline)
	c.Next()
}

func hookRecover(e interface{}) {
	defer func() {
		if err := recover(); err != nil {
			logger.Runtime.Error(err.(string))
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
