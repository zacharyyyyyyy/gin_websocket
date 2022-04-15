package router

import (
	"gin_websocket/controller/admin/login"
	"gin_websocket/controller/admin/role"
	"gin_websocket/controller/admin/ws"
	"gin_websocket/middleware/router_middleware"

	"github.com/gin-gonic/gin"
)

func initAdminRoute(r *gin.Engine) {
	adminRoute := r.Group("/admin")
	adminRoute.POST("/login", login.Login)
	adminRoute.Use(router_middleware.AdminAuthentication())
	{
		adminRoute.GET("/logout", login.Logout)
		adminRoute.POST("/info", ws.Info)
		adminRoute.GET("/service_link", ws.ServiceLink)
		adminRoute.POST("/admin_auth", role.GetAllAdminAuth)
	}
}
