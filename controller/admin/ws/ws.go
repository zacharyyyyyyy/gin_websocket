package ws

import (
	"context"
	"errors"
	"gin_websocket/controller"
	ws "gin_websocket/service/websocket"
	"github.com/gin-gonic/gin"
)

func ServiceLink(c *gin.Context) {
	ctx, _ := context.WithCancel(context.Background())
	serviceClient, err := ws.NewCustomerService(ctx, c)
	if err != nil {
		controller.PanicResponse(c, err)
		return
	}
	for {
		err := serviceClient.Receive()
		if errors.Is(err, ws.CloseErr) {
			break
		}
	}

}

func Info(c *gin.Context) {
	resp := make(map[string]interface{}, 2)
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
	for wsKey, serviceClient := range ws.CustomerServiceContainerHandle.WebsocketCustomerServiceMap {
		serviceClientMap["ws_key"] = wsKey
		var serviceClientId interface{}
		if serviceClient.GetBindUser() != nil {
			serviceClientId = serviceClient.GetBindUser().Id
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
