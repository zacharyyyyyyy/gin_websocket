package router_middleware

import (
	"errors"
	"net/http"
	"strconv"

	"gin_websocket/controller"
	"gin_websocket/lib/redis"

	"github.com/gin-gonic/gin"
)

func LoginLimit(c *gin.Context) {
	ip := c.ClientIP()
	if ip == "::1" {
		ip = "127.0.0.1"
	}
	count, err := redis.RedisDb.Get("login_limit_ip_" + ip)
	ipCount, _ := strconv.Atoi(count)
	if err == nil && ipCount >= 5 {
		controller.PanicResponse(c, errors.New("错误次数超限"), http.StatusInternalServerError, "今日错误次数超限，请明日重试")
		c.Abort()
		return
	}
	c.Next()
}
