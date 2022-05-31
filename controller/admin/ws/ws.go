package ws

import (
	"context"
	"errors"
	"net/http"

	"gin_websocket/controller"
	"gin_websocket/dao"
	"gin_websocket/lib/validator"
	ws "gin_websocket/service/websocket"

	"github.com/gin-gonic/gin"
)

func ServiceLink(c *gin.Context) {
	param := new(struct {
		WsKey string `form:"ws_key" binding:"required,min=1" msg:"ws_key为字符串且不能为空"`
	})
	if err := c.ShouldBind(param); err != nil {
		errMsg := validator.GetValidMsg(err, param)
		controller.PanicResponse(c, err, http.StatusInternalServerError, errMsg)
		return
	}
	ctx, _ := context.WithCancel(context.Background())
	adminStruct, err := dao.GetAdminCurrent(c.Request)
	if err != nil {
		controller.PanicResponse(c, err, http.StatusInternalServerError, "未登录")
		return
	}

	serviceClient, err := ws.NewCustomerService(ctx, c.Request, c.Writer, c.ClientIP(), adminStruct.Id)
	if err != nil {
		controller.PanicResponse(c, err, http.StatusInternalServerError, "")
		return
	}
	for {
		err := serviceClient.Receive(ws.WsKey(param.WsKey))
		if errors.Is(err, ws.CloseErr) {
			break
		}
	}
}

func GetLinkUser(c *gin.Context) {

}

func Info(c *gin.Context) {
	resp := make(map[string]interface{}, 0)
	userClientMap, serviceClientMap := make(map[string]interface{}, 0), make(map[string]interface{}, 0)
	userClientMapSlice, serviceClientMapSlice := make([]map[string]interface{}, 0), make([]map[string]interface{}, 0)
	userCount := ws.WsContainerHandle.GetConnCount()
	serviceCount := ws.CustomerServiceContainerHandle.GetConnCount()

	for wsKey, userClient := range ws.WsContainerHandle.WebSocketClientMap {
		userClientMap["ws_key"] = wsKey
		var userClientId interface{}
		if userClient.GetBindService() != nil {
			userClientId = userClient.GetBindService().Id
		}
		userClientMap["bind_service"] = userClientId
		userClientMapSlice = append(userClientMapSlice, userClientMap)

	}
	for adminId, serviceClient := range ws.CustomerServiceContainerHandle.WebsocketCustomerServiceMap {
		serviceClientMap["ws_key"] = serviceClient.Id
		serviceClientMap["admin_id"] = adminId
		admin, _ := dao.SelectOneById(adminId)
		serviceClientMap["admin_name"] = admin.Name
		bindUserMap := make([]string, 0)
		var serviceClientId interface{}
		for _, userClient := range serviceClient.GetAllBindUser() {
			bindUserMap = append(bindUserMap, string(userClient.Id))
		}
		serviceClientMap["bind_user"] = serviceClientId
		serviceClientMapSlice = append(serviceClientMapSlice, serviceClientMap)
	}
	resp["user_count"] = userCount
	resp["service_count"] = serviceCount
	resp["user_clients"] = userClientMapSlice
	resp["service_clients"] = serviceClientMapSlice
	respStruct := controller.ResponseStruct{
		C:    c,
		Data: resp,
	}
	respStruct.JsonResponse()
}
