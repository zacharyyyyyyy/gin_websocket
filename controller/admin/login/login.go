package login

import (
	"net/http"

	"gin_websocket/controller"
	"gin_websocket/service/admin"

	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	err := admin.Login(username, password, c.Request, c.Writer)
	if err != nil {
		controller.PanicResponse(c, err, http.StatusInternalServerError)
		return
	}
	controller.QuickSuccessResponse(c)
}

func Logout(c *gin.Context) {
	_ = admin.Logout(c.Request, c.Writer)
	controller.QuickSuccessResponse(c)
}
