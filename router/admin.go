package router

import (
	"gin_websocket/controller/admin/ws"
	"github.com/gin-gonic/gin"
)

func initAdminRoute(r *gin.Engine) {
	adminRoute := r.Group("/admin")
	//adminRoute.Use(router_middleware.AdminAuthentication())
	{
		adminRoute.POST("/info", ws.Info)
		adminRoute.GET("/service_link", ws.ServiceLink)
	}
}
