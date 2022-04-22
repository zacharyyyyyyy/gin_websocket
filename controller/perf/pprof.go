package perf

import (
	"net/http/pprof"

	"github.com/gin-gonic/gin"
)

func IndexPprof(c *gin.Context) {
	pprof.Index(c.Writer, c.Request)
}

func CmdLinePprof(c *gin.Context) {
	pprof.Cmdline(c.Writer, c.Request)
}

func ProfilePprof(c *gin.Context) {
	pprof.Profile(c.Writer, c.Request)
}

func SymbolPprof(c *gin.Context) {
	pprof.Symbol(c.Writer, c.Request)
}

func TracePprof(c *gin.Context) {
	pprof.Trace(c.Writer, c.Request)
}

func AllocsPprof(c *gin.Context) {
	pprof.Handler("allocs").ServeHTTP(c.Writer, c.Request)
}

func BlockPprof(c *gin.Context) {
	pprof.Handler("block").ServeHTTP(c.Writer, c.Request)
}

func GoroutinePprof(c *gin.Context) {
	pprof.Handler("goroutine").ServeHTTP(c.Writer, c.Request)
}

func HeapPprof(c *gin.Context) {
	pprof.Handler("heap").ServeHTTP(c.Writer, c.Request)
}

func MutexPprof(c *gin.Context) {
	pprof.Handler("mutex").ServeHTTP(c.Writer, c.Request)
}

func ThreadCreatePprof(c *gin.Context) {
	pprof.Handler("threadcreate").ServeHTTP(c.Writer, c.Request)
}
