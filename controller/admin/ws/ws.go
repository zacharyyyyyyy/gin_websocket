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

func Register(c *gin.Context) {
	adminStruct, err := dao.GetAdminCurrent(c.Request)
	if err != nil {
		controller.PanicResponse(c, err, http.StatusInternalServerError, "未登录")
		return
	}
	ws.RegisterService(c.ClientIP(), adminStruct.Id)
	controller.QuickSuccessResponse(c)
}

func Cancel(c *gin.Context) {
	adminStruct, err := dao.GetAdminCurrent(c.Request)
	if err != nil {
		controller.PanicResponse(c, err, http.StatusInternalServerError, "未登录")
		return
	}
	_ = ws.CustomerServiceContainerHandle.Remove(adminStruct.Id)
	controller.QuickSuccessResponse(c)
}

func ServiceLink(c *gin.Context) {
	param := new(struct {
		WsKey string `form:"ws_key" binding:"required,min=1"`
	})
	if err := c.BindQuery(param); err != nil {
		errMsg := validator.GetValidMsg(err)
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
	adminStruct, err := dao.GetAdminCurrent(c.Request)
	resp := make(map[string]interface{}, 0)
	userClients := make([]map[string]interface{}, 0)
	if err != nil {
		controller.PanicResponse(c, err, http.StatusInternalServerError, "未登录")
		return
	}
	customerService, err := ws.CustomerServiceContainerHandle.GetCustomerService(adminStruct.Id)
	if err == nil {
		for wsKey, userClient := range customerService.GetAllBindUser() {
			userClients = append(userClients, map[string]interface{}{
				"ip":     userClient.Ip,
				"ws_key": wsKey,
			})
		}
	}
	resp["user"] = userClients
	resp["user_count"] = len(userClients)
	respStruct := controller.ResponseStruct{
		C:    c,
		Data: resp,
	}
	respStruct.JsonResponse()
}

func Info(c *gin.Context) {
	resp := make(map[string]interface{}, 0)
	userClientMapSlice, serviceClientMapSlice := make([]map[string]interface{}, 0), make([]map[string]interface{}, 0)
	userCount := ws.WsContainerHandle.GetConnCount()
	serviceCount := ws.CustomerServiceContainerHandle.GetConnCount()
	for wsKey, userClient := range ws.WsContainerHandle.WebSocketClientMap {
		var userClientId interface{}
		if userClient.GetBindService() != nil {
			userClientId = userClient.GetBindService().Id
		}
		userClientMapSlice = append(userClientMapSlice, map[string]interface{}{
			"ws_key":       wsKey,
			"bind_service": userClientId,
		})
	}
	for adminId, serviceClient := range ws.CustomerServiceContainerHandle.WebsocketCustomerServiceMap {
		admin, _ := dao.SelectOneById(adminId)
		bindUserMap := make([]string, 0)
		var serviceClientId interface{}
		for _, userClient := range serviceClient.GetAllBindUser() {
			bindUserMap = append(bindUserMap, string(userClient.Id))
		}
		serviceClientMapSlice = append(serviceClientMapSlice, map[string]interface{}{
			"ws_key":     serviceClient.Id,
			"admin_id":   adminId,
			"admin_name": admin.Name,
			"bind_user":  serviceClientId,
		})
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
