package Router

import (
	"gin_websocket/Controller/admin"
	"github.com/gin-gonic/gin"
)

func initAdminRoute(r *gin.Engine) {
	adminRoute := r.Group("/admin")
	{
		adminRoute.GET("/link", admin.Link)
		adminRoute.GET("/ping", admin.Ping)
	}
}
