package websocket

import (
	"context"
	"sync"

	"gin_websocket/lib/logger"
	"github.com/gin-gonic/gin"
)

//客服容器
type CustomerServiceContainer struct {
	WebsocketCustomerServiceMap map[WsKey]*CustomerServiceClient
	lock                        *sync.RWMutex
	CustomerWebsocketCount      int
}

//容器加载
func ServiceStart() *CustomerServiceContainer {
	customerService := make(map[WsKey]*CustomerServiceClient, 1)
	CustomerServiceContainerHandle := &CustomerServiceContainer{
		WebsocketCustomerServiceMap: customerService,
		lock:                        &sync.RWMutex{},
		CustomerWebsocketCount:      0,
	}
	return CustomerServiceContainerHandle
}

//客服连入初始化
func (Cont *CustomerServiceContainer) NewClient(ctx context.Context, c *gin.Context) error {
	var customerServiceClient *CustomerServiceClient
	customerServiceClient, err := newCustomerService(ctx, c)
	if err != nil {
		logger.Service.Error(err.Error())
		return err
	}
	Cont.append(customerServiceClient)
	return nil
}

func (Cont CustomerServiceContainer) GetConnCount() int {
	return Cont.CustomerWebsocketCount
}

func (Cont *CustomerServiceContainer) Remove(customerServiceClient *CustomerServiceClient) error {
	Cont.lock.Lock()
	defer Cont.lock.Unlock()
	if _, ok := Cont.WebsocketCustomerServiceMap[customerServiceClient.Id]; !ok {
		return ClientNotFoundErr
	}
	//先释放链接
	err := Cont.WebsocketCustomerServiceMap[customerServiceClient.Id].close()
	if err != nil {
		return err
	}
	delete(Cont.WebsocketCustomerServiceMap, customerServiceClient.Id)
	Cont.CustomerWebsocketCount--
	return nil
}

func (Cont *CustomerServiceContainer) append(customerServiceClient *CustomerServiceClient) {
	Cont.lock.Lock()
	defer Cont.lock.Unlock()
	Cont.WebsocketCustomerServiceMap[customerServiceClient.Id] = customerServiceClient
	Cont.CustomerWebsocketCount++
}
