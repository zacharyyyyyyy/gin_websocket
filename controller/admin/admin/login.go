package admin

import (
	"gin_websocket/controller"
	"gin_websocket/service/admin"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	err := admin.Login(username, password, c.Request, c.Writer)
	if err != nil {
		respStruct := controller.ResponseStruct{
			C:    c,
			Data: nil,
			Code: http.StatusInternalServerError,
		}
		respStruct.SetMessage(err.Error())
		respStruct.JsonResponse()
		return
	}
	controller.QuickSuccessResponse(c)
}

func Logout(c *gin.Context) {
	_ = admin.Logout(c.Request, c.Writer)
	controller.QuickSuccessResponse(c)
}
