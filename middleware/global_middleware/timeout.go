package global_middleware

import (
	"fmt"
	"net/http"

	"gin_websocket/controller"
	"gin_websocket/lib/logger"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/gin-gonic/gin"
)

func init() {
	hystrix.ConfigureCommand("default", hystrix.CommandConfig{
		Timeout:                15000, // 单次请求 超时时间
		MaxConcurrentRequests:  500,   // 最大并发量
		SleepWindow:            1000,  // 熔断后多久去尝试服务是否可用
		RequestVolumeThreshold: 1,     // 验证熔断的 请求数量, 10秒内采样
		ErrorPercentThreshold:  1,     // 验证熔断的 错误百分比
	})
}

func Timeout(c *gin.Context) {
	fmt.Println(c.Writer.Status())
	hystrix.Do("timeout_handle", func() error {
		c.Next()
		return nil
	}, func(err error) error {
		if err != nil {
			newErr := fmt.Errorf("%s %w", c.Request.URL, err)
			controller.PanicResponse(c, newErr, http.StatusInternalServerError, "")
			logger.Runtime.Error(newErr.Error())
			c.Abort()
			return nil
		}
		return nil
	})
}
