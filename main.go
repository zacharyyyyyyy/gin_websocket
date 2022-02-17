package gin_websocket

import (
	"github.com/fvbock/endless"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"syscall"

	"gin_websocket/Router"
)

func main() {
	var handler http.Handler
	handler = route()
	server := endless.NewServer(":8080", handler)
	server.BeforeBegin = func(addr string) {
		log.Printf("start http server listening %s, pid is %d", addr, syscall.Getpid())
	}
	err := server.ListenAndServe()
	if err != nil {
		log.Printf("server err: %v", err)
	}
}

func route() *gin.Engine {
	r := gin.New()
	if gin.IsDebugging() {
		r.Use(gin.Logger())
	}
	r.Use(gzip.Gzip(gzip.DefaultCompression))
	Router.InitRouter(r)
	return r
}
