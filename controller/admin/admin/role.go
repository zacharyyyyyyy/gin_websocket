package admin

import (
	"net/http"
	"time"

	"gin_websocket/controller"
	"gin_websocket/dao"
	"gin_websocket/lib/validator"

	"github.com/gin-gonic/gin"
)

func GetAllAdminAuth(c *gin.Context) {
	param := new(struct {
		Pn int `form:"pn" binding:"required,min=1" msg:"pn为整型且最小值为1"`
		Pc int `form:"pc" binding:"required,min=1" msg:"pc为整型且最小值为1"`
	})
	if err := c.Bind(param); err != nil {
		errMsg := validator.GetValidMsg(err, param)
		controller.PanicResponse(c, err, http.StatusInternalServerError, errMsg)
		return
	}
	authData := make([]interface{}, 0)
	data := make(map[string]interface{}, 0)
	page := param.Pn
	limit := param.Pc
	offset := (page - 1) * limit
	result, err := dao.GetAuthByLimitAndOffset(limit, offset)
	if err != nil {
		controller.PanicResponse(c, err, http.StatusInternalServerError, "")
		return
	}
	count, err := dao.GetAllAuthCountByEnable()
	if err != nil {
		controller.PanicResponse(c, err, http.StatusInternalServerError, "")
		return
	}
	for _, auth := range result {
		authData = append(authData, map[string]interface{}{
			"id":          auth.Id,
			"name":        auth.Name,
			"method":      auth.Method,
			"path":        auth.Path,
			"enable":      auth.Enable,
			"create_time": time.Unix(int64(auth.CreateTime), 0).Format("2006-01-02 15:04:05"),
		})
	}
	data["data"] = authData
	data["count"] = int(count)
	data["pn"] = param.Pn
	data["pc"] = param.Pc
	ctl := controller.ResponseStruct{
		C:    c,
		Data: data,
		Code: http.StatusOK,
	}
	ctl.JsonResponse()
}
