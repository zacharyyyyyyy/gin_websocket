package router

import (
	"gin_websocket/controller/admin/admin"
	"gin_websocket/controller/admin/ws"
	"gin_websocket/controller/perf"
	//"gin_websocket/middleware/router_middleware"

	"github.com/gin-gonic/gin"
)

func initAdminRoute(r *gin.Engine) {
	adminRoute := r.Group("/admin")
	adminRoute.POST("/login", admin.Login)
	//adminRoute.Use(router_middleware.AdminAuthentication())
	{
		adminRoute.POST("/all_admin_user", admin.GetAllAdmin)
		adminRoute.POST("/user/add", admin.AddAdmin)
		adminRoute.GET("/logout", admin.Logout)
		adminRoute.POST("/info", ws.Info)
		adminRoute.GET("/service_link", ws.ServiceLink)
		adminRoute.POST("/admin_auth", admin.GetAllAdminAuth)
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
	}
}
