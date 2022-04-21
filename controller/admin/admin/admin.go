package admin

import (
	"gin_websocket/dao"
	"net/http"
	"time"

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
	adminData := make([]interface{}, 0)
	data := make(map[string]interface{}, 0)
	page := param.Pn
	limit := param.Pc
	offset := (page - 1) * limit
	result, err := dao.GetAllAdminByLimitAndOffset(limit, offset)
	if err != nil {
		controller.PanicResponse(c, err, http.StatusInternalServerError, "")
		return
	}
	count, err := dao.GetAdminCount()
	if err != nil {
		controller.PanicResponse(c, err, http.StatusInternalServerError, "")
		return
	}
	for _, admin := range result {
		adminData = append(adminData, map[string]interface{}{
			"id":          admin.Id,
			"name":        admin.Name,
			"user_name":   admin.Username,
			"role_name":   admin.RoleName,
			"describe":    admin.Describe,
			"create_time": time.Unix(int64(admin.CreateTime), 0).Format("2006-01-02 15:04:05"),
		})
	}
	data["data"] = adminData
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
