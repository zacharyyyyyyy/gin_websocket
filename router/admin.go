package router

import (
	"gin_websocket/controller/admin/admin"
	"gin_websocket/controller/admin/ws"
	"gin_websocket/controller/perf"

	"gin_websocket/middleware/router_middleware"

	"github.com/gin-gonic/gin"
)

func initAdminRoute(r *gin.Engine) {
	adminRoute := r.Group("/admin")
	adminRoute.POST("/login", router_middleware.LoginLimit, admin.Login)
	adminRoute.Use(router_middleware.AdminAuthentication())
	{
		adminRoute.GET("/logout", admin.Logout)
		adminRoute.POST("/user", admin.GetAllAdmin)
		adminRoute.POST("/user/add", admin.AddAdmin)
		adminRoute.POST("/user/edit", admin.EditAdmin)
		adminRoute.POST("/user/del", admin.DelAdmin)

		adminRoute.POST("/auth", admin.GetAllAdminAuth)
		adminRoute.POST("/role", admin.GetAllRole)
		adminRoute.POST("/role/add", admin.AddRole)
		adminRoute.POST("/role/edit", admin.EditRole)
		adminRoute.POST("/role/del", admin.DelRole)

		adminRoute.POST("/auth_map", admin.GetAllRoleAuth)
		adminRoute.POST("/auth_map/add", admin.AddAuthMap)
		adminRoute.POST("/auth_map/edit", admin.EditAuthMap)
		adminRoute.POST("/auth_map/del", admin.DelAuthMap)

		adminRoute.POST("/ws", ws.Register)
		adminRoute.POST("/ws/info", ws.Info)
		adminRoute.GET("/ws/link", ws.ServiceLink)
		adminRoute.POST("/ws/link_user", ws.GetLinkUser)
		adminRoute.POST("/ws/close", ws.Cancel)

	}
	//pprof采集
	{
		adminRoute.GET("/perf/pprof", perf.IndexPprof)
		adminRoute.GET("/perf/cmdline", perf.CmdLinePprof)
		adminRoute.GET("/perf/profile", perf.ProfilePprof)
		adminRoute.GET("/perf/symbol", perf.SymbolPprof)
		adminRoute.GET("/perf/trace", perf.TracePprof)
		adminRoute.GET("/perf/allocs", perf.AllocsPprof)
		adminRoute.GET("/perf/block", perf.BlockPprof)
		adminRoute.GET("/perf/goroutine", perf.GoroutinePprof)
		adminRoute.GET("/perf/heap", perf.HeapPprof)
		adminRoute.GET("/perf/mutex", perf.MutexPprof)
		adminRoute.GET("/perf/threadcreate", perf.ThreadCreatePprof)
		adminRoute.POST("/perf/write", perf.WritePprof)
	}
}
