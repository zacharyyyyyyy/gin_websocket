package perf

import (
	"gin_websocket/controller"
	"net/http"
	"net/http/pprof"
	"os"
	rpprof "runtime/pprof"
	"time"

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

func WritePprof(c *gin.Context) {
	path := "pprof"
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.Mkdir("pprof", os.ModePerm)
			if err != nil {
				controller.PanicResponse(c, err, http.StatusInternalServerError, "文件夹创建失败")
				return
			}
		}
	}
	cpuProfile, err := os.Create("pprof/cpu_profile_" + time.Now().Format("2006-01-02_15-04-05"))
	if err != nil {
		controller.PanicResponse(c, err, http.StatusInternalServerError, "cpu文件创建失败")
		return
	}
	defer cpuProfile.Close()
	memProfile, err := os.Create("pprof/mem_profile_" + time.Now().Format("2006-01-02_15-04-05"))
	if err != nil {
		controller.PanicResponse(c, err, http.StatusInternalServerError, "mem文件创建失败")
		return
	}
	defer memProfile.Close()
	//采集CPU信息
	rpprof.StartCPUProfile(cpuProfile)
	defer rpprof.StopCPUProfile()
	//采集内存信息
	rpprof.WriteHeapProfile(memProfile)
	controller.QuickSuccessResponse(c)
}
