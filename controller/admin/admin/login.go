package admin

import (
	"gin_websocket/controller"
	"gin_websocket/lib/validator"
	"gin_websocket/service/admin"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	param := new(struct {
		Username string `form:"username" binding:"required,min=1"`
		Password string `form:"password" binding:"required,min=1"`
	})
	if err := c.ShouldBind(param); err != nil {
		errMsg := validator.GetValidMsg(err)
		controller.PanicResponse(c, err, http.StatusInternalServerError, errMsg)
		return
	}
	err := admin.Login(param.Username, param.Password, c.Request, c.Writer, c.ClientIP())
	if err != nil {
		controller.PanicResponse(c, err, http.StatusInternalServerError, err.Error())
		return
	}
	controller.QuickSuccessResponse(c)
}

func Logout(c *gin.Context) {
	_ = admin.Logout(c.Request, c.Writer)
	controller.QuickSuccessResponse(c)
}
