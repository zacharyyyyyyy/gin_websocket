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
	r.Use(gzip.Gzip(gzip.DefaultCompression), global_middleware.Cors, global_middleware.HttpTrace, global_middleware.HttpRecover, global_middleware.Timeout)
	r.NoRoute(global_middleware.NoRouteHandle)
	initAdminRoute(r)
	initApiRoute(r)
	return r

}
