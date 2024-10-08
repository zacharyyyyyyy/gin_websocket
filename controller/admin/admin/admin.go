package admin

import (
	"html"
	"net/http"
	"time"

	"gin_websocket/controller"
	"gin_websocket/dao"
	"gin_websocket/lib/validator"
	"gin_websocket/service/admin"

	"github.com/gin-gonic/gin"
)

func GetAllAdmin(c *gin.Context) {
	param := new(struct {
		Pn int `form:"pn" binding:"required,min=1"`
		Pc int `form:"pc" binding:"required,min=1"`
	})
	if err := c.ShouldBind(param); err != nil {
		errMsg := validator.GetValidMsg(err)
		controller.PanicResponse(c, err, http.StatusInternalServerError, errMsg)
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
	for _, adminResult := range result {
		adminData = append(adminData, map[string]interface{}{
			"id":          adminResult.Id,
			"name":        adminResult.Name,
			"user_name":   adminResult.Username,
			"role_name":   adminResult.RoleName,
			"describe":    adminResult.Describe,
			"create_time": time.Unix(int64(adminResult.CreateTime), 0).Format("2006-01-02 15:04:05"),
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

func AddAdmin(c *gin.Context) {
	param := new(struct {
		Username string `form:"username" binding:"required,min=2"`
		Password string `form:"password" binding:"required"`
		Name     string `form:"name" binding:"required"`
		Role     int    `form:"role" binding:"required,existsAdminRole"`
	})
	if err := c.ShouldBind(param); err != nil {
		errMsg := validator.GetValidMsg(err)
		controller.PanicResponse(c, err, http.StatusInternalServerError, errMsg)
		return
	}
	param.Username = html.EscapeString(param.Username)
	param.Name = html.EscapeString(param.Name)
	param.Password = admin.ChangePassword(param.Password)
	_ = dao.AddAdmin(param.Username, param.Name, param.Password, param.Role)
	controller.QuickSuccessResponse(c)
}

func EditAdmin(c *gin.Context) {
	param := new(struct {
		Username string `form:"username" binding:"required,min=2"`
		Password string `form:"password" binding:"required"`
		Name     string `form:"name" binding:"required"`
		Role     int    `form:"role" binding:"required,existsAdminRole"`
		Id       int    `form:"id" binding:"required"`
	})
	if err := c.ShouldBind(param); err != nil {
		errMsg := validator.GetValidMsg(err)
		controller.PanicResponse(c, err, http.StatusInternalServerError, errMsg)
		return
	}
	param.Username = html.EscapeString(param.Username)
	param.Name = html.EscapeString(param.Name)
	param.Password = admin.ChangePassword(param.Password)
	_ = dao.EditAdmin(param.Username, param.Name, param.Password, param.Role, param.Id)
	controller.QuickSuccessResponse(c)
}

func DelAdmin(c *gin.Context) {
	param := new(struct {
		Id int `form:"id" binding:"required"`
	})
	if err := c.ShouldBind(param); err != nil {
		errMsg := validator.GetValidMsg(err)
		controller.PanicResponse(c, err, http.StatusInternalServerError, errMsg)
		return
	}
	_ = dao.DelAdmin(param.Id)
	controller.QuickSuccessResponse(c)
}
