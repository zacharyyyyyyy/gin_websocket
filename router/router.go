package router

import (
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.New()
	if gin.IsDebugging() {
		r.Use(gin.Logger())
	}
	r.Use(gzip.Gzip(gzip.DefaultCompression))
	initAdminRoute(r)
	return r

}
