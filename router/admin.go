package router

import (
	"gin_websocket/controller/admin/admin"
	"gin_websocket/controller/admin/ws"
	"github.com/gin-gonic/gin"
)

func initAdminRoute(r *gin.Engine) {
	adminRoute := r.Group("/admin")
	adminRoute.POST("/login", admin.Login)
	//adminRoute.Use(router_middleware.AdminAuthentication())
	{
		adminRoute.POST("/all_admin_user", admin.GetAllAdmin)
		adminRoute.GET("/logout", admin.Logout)
		adminRoute.POST("/info", ws.Info)
		adminRoute.GET("/service_link", ws.ServiceLink)
		adminRoute.POST("/admin_auth", admin.GetAllAdminAuth)
	}
}
