package admin

import (
	"html"
	"net/http"
	"time"

	"gin_websocket/controller"
	"gin_websocket/dao"
	"gin_websocket/lib/validator"

	"github.com/gin-gonic/gin"
)

func GetAllRole(c *gin.Context) {
	result, err := dao.GetAllRole()
	if err != nil {
		controller.PanicResponse(c, err, http.StatusInternalServerError, "")
		return
	}
	roleData := make([]interface{}, 0)
	data := make(map[string]interface{}, 0)

	for _, role := range result {
		roleData = append(roleData, map[string]interface{}{
			"id":          role.Id,
			"name":        role.Name,
			"describe":    role.Describe,
			"create_time": time.Unix(int64(role.CreateTime), 0).Format("2006-01-02 15:04:05"),
		})
	}
	data["data"] = roleData
	ctl := controller.ResponseStruct{
		C:    c,
		Data: data,
		Code: http.StatusOK,
	}
	ctl.JsonResponse()
}

func AddRole(c *gin.Context) {
	param := new(struct {
		Name     string `form:"name" binding:"required,min=1" msg:"name为字符串且不能为空"`
		Describe string `form:"describe" binding:"min=1" msg:"describe为字符串"`
	})
	if err := c.ShouldBind(param); err != nil {
		errMsg := validator.GetValidMsg(err, param)
		controller.PanicResponse(c, err, http.StatusInternalServerError, errMsg)
		return
	}
	param.Name = html.EscapeString(param.Name)
	if len(param.Describe) > 0 {
		param.Describe = html.EscapeString(param.Describe)
	}
	_ = dao.AddRole(param.Name, param.Describe)
	controller.QuickSuccessResponse(c)
}

func EditRole(c *gin.Context) {
	param := new(struct {
		Name     string `form:"name" binding:"required,min=1" msg:"name为字符串且不能为空"`
		Describe string `form:"describe" binding:"min=1" msg:"describe为字符串"`
		Id       int    `form:"id" binding:"required,min=1" msg:"id为整型且不能为空"`
	})
	if err := c.ShouldBind(param); err != nil {
		errMsg := validator.GetValidMsg(err, param)
		controller.PanicResponse(c, err, http.StatusInternalServerError, errMsg)
		return
	}
	param.Name = html.EscapeString(param.Name)
	if len(param.Describe) > 0 {
		param.Describe = html.EscapeString(param.Describe)
	}
	if err := dao.EditRole(param.Name, param.Describe, param.Id); err != nil {
		controller.PanicResponse(c, err, http.StatusInternalServerError, "")
		return
	}
	controller.QuickSuccessResponse(c)
	return
}

func DelRole(c *gin.Context) {
	param := new(struct {
		Id int `form:"id" binding:"required,min=1" msg:"id为整型且不能为空"`
	})
	if err := c.ShouldBind(param); err != nil {
		errMsg := validator.GetValidMsg(err, param)
		controller.PanicResponse(c, err, http.StatusInternalServerError, errMsg)
		return
	}
	_ = dao.DelRole(param.Id)
	controller.QuickSuccessResponse(c)
}

func GetAllRoleAuth(c *gin.Context) {
	param := new(struct {
		Pn   int `form:"pn" binding:"required,min=1" msg:"pn为整型且最小值为1"`
		Pc   int `form:"pc" binding:"required,min=1" msg:"pc为整型且最小值为1"`
		Role int `form:"role" binding:"required,min=1" msg:"role为整型且最小值为1"`
	})
	if err := c.ShouldBind(param); err != nil {
		errMsg := validator.GetValidMsg(err, param)
		controller.PanicResponse(c, err, http.StatusInternalServerError, errMsg)
		return
	}
	authMapData := make([]interface{}, 0)
	data := make(map[string]interface{}, 0)
	page := param.Pn
	limit := param.Pc
	offset := (page - 1) * limit

	result, err := dao.GetAllAuthMapByRole(limit, offset, param.Role)
	if err != nil {
		controller.PanicResponse(c, err, http.StatusInternalServerError, "")
		return
	}
	count, err := dao.GetAuthMapCount()
	if err != nil {
		controller.PanicResponse(c, err, http.StatusInternalServerError, "")
		return
	}
	for _, authMap := range result {
		authMapData = append(authMapData, map[string]interface{}{
			"role":          authMap.Role,
			"auth":          authMap.Auth,
			"role_name":     authMap.RoleName,
			"role_describe": authMap.RoleDescribe,
			"auth_name":     authMap.AuthName,
		})
	}
	data["data"] = authMapData
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

func AddAuthMap(c *gin.Context) {
	param := new(struct {
		Role int `form:"role" binding:"required,min=1" msg:"role为整型且不能为空"`
		Auth int `form:"auth" binding:"min=1" msg:"auth为整型且不能为空"`
	})
	if err := c.ShouldBind(param); err != nil {
		errMsg := validator.GetValidMsg(err, param)
		controller.PanicResponse(c, err, http.StatusInternalServerError, errMsg)
		return
	}

	_ = dao.AddAuth(param.Role, param.Auth)
	controller.QuickSuccessResponse(c)
}

func EditAuthMap(c *gin.Context) {
	param := new(struct {
		Role int `form:"role" binding:"required,min=1" msg:"role为整型且不能为空"`
		Auth int `form:"auth" binding:"min=1" msg:"auth为整型且不能为空"`
	})
	if err := c.ShouldBind(param); err != nil {
		errMsg := validator.GetValidMsg(err, param)
		controller.PanicResponse(c, err, http.StatusInternalServerError, errMsg)
		return
	}
	_ = dao.EditAuth(param.Role, param.Auth)
	controller.QuickSuccessResponse(c)
}

func DelAuthMap(c *gin.Context) {
	param := new(struct {
		Role int `form:"role" binding:"required,min=1" msg:"role为整型且不能为空"`
		Auth int `form:"auth" binding:"min=1" msg:"auth为整型且不能为空"`
	})
	if err := c.ShouldBind(param); err != nil {
		errMsg := validator.GetValidMsg(err, param)
		controller.PanicResponse(c, err, http.StatusInternalServerError, errMsg)
		return
	}
	_ = dao.DelAuth(param.Role, param.Auth)
	controller.QuickSuccessResponse(c)
}
