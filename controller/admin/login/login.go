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

	err := admin.Login(username, password)
	if err != nil {
		controller.PanicResponse(c, err, http.StatusInternalServerError)
	}
	controller.QuickSuccessResponse(c)
}

func Logout(c *gin.Context) {
	_ = admin.Logout(c)
	controller.QuickSuccessResponse(c)
}
