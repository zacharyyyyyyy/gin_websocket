package global_middleware

import (
	"errors"
	"net/http"

	"gin_websocket/controller"

	"github.com/gin-gonic/gin"
)

func NoRouteHandle(c *gin.Context) {
	controller.PanicResponse(c, errors.New("请求地址不存在"), http.StatusNotFound, "请求地址不存在")
}
