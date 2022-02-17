package Router

import "github.com/gin-gonic/gin"

func InitRouter(r *gin.Engine) {
	initAdminRoute(r)
}
