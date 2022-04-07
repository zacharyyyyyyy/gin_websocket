package router

import (
	apiWs "gin_websocket/controller/api/ws"
	"github.com/gin-gonic/gin"
)

func initApiRoute(r *gin.Engine) {
	apiRoute := r.Group("/api")
	{
		apiRoute.GET("/link", apiWs.Link)
	}
}
