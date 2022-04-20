package admin

import (
	"net/http"

	"gin_websocket/controller"
	"gin_websocket/lib/validator"

	"github.com/gin-gonic/gin"
)

func GetAllAdmin(c *gin.Context) {
	param := new(struct {
		Pn int `form:"pn" binding:"required,min=1" msg:"pn为整型且最小值为1"`
		Pc int `form:"pc" binding:"required,min=1" msg:"pc为整型且最小值为1"`
	})
	if err := c.Bind(param); err != nil {
		errMsg := validator.GetValidMsg(err, param)
		controller.PanicResponse(c, err, http.StatusNotImplemented, errMsg)
		return
	}
}
