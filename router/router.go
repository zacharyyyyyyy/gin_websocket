package router

import (
	"gin_websocket/middleware/global_middleware"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.New()
	if gin.IsDebugging() {
		r.Use(gin.Logger())
	}
	r.Use(global_middleware.Cors)
	r.Use(gzip.Gzip(gzip.DefaultCompression))
	initAdminRoute(r)
	initApiRoute(r)
	return r

}
