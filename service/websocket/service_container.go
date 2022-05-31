package websocket

import (
	"sync"
)

//客服容器
type CustomerServiceContainer struct {
	WebsocketCustomerServiceMap map[int]*CustomerServiceClient
	lock                        *sync.RWMutex
	CustomerWebsocketCount      int
}

//容器加载
func serviceStart() *CustomerServiceContainer {
	customerService := make(map[int]*CustomerServiceClient, 0)
	CustomerServiceContainerHandle := &CustomerServiceContainer{
		WebsocketCustomerServiceMap: customerService,
		lock:                        &sync.RWMutex{},
		CustomerWebsocketCount:      0,
	}
	return CustomerServiceContainerHandle
}

//客服连入初始化
func (Cont *CustomerServiceContainer) NewClient(customerServiceClient *CustomerServiceClient) error {
	return Cont.append(customerServiceClient)
}

func (Cont CustomerServiceContainer) GetConnCount() int {
	return Cont.CustomerWebsocketCount
}

//主动删除
func (Cont *CustomerServiceContainer) Remove(customerServiceClient *CustomerServiceClient) error {
	if _, ok := Cont.WebsocketCustomerServiceMap[customerServiceClient.AdminId]; !ok {
		return ClientNotFoundErr
	}
	err := Cont.WebsocketCustomerServiceMap[customerServiceClient.AdminId].Close()
	return err
}

func (Cont *CustomerServiceContainer) append(customerServiceClient *CustomerServiceClient) error {
	Cont.lock.Lock()
	defer Cont.lock.Unlock()
	if _, ok := Cont.WebsocketCustomerServiceMap[customerServiceClient.AdminId]; !ok {
		Cont.WebsocketCustomerServiceMap[customerServiceClient.AdminId] = customerServiceClient
		Cont.CustomerWebsocketCount++
		return nil
	}
	return ClientAlreadyInContainer
}

func (Cont *CustomerServiceContainer) remove(customerServiceClient *CustomerServiceClient) error {
	Cont.lock.Lock()
	defer Cont.lock.Unlock()
	if _, ok := Cont.WebsocketCustomerServiceMap[customerServiceClient.AdminId]; !ok {
		return ClientNotFoundErr
	}
	delete(Cont.WebsocketCustomerServiceMap, customerServiceClient.AdminId)
	Cont.CustomerWebsocketCount--
	return nil
}

func getCustomerService() (*CustomerServiceClient, error) {
	CustomerServiceContainerHandle.lock.Lock()
	defer CustomerServiceContainerHandle.lock.Unlock()
	for _, serviceClient := range CustomerServiceContainerHandle.WebsocketCustomerServiceMap {
		delete(CustomerServiceContainerHandle.WebsocketCustomerServiceMap, serviceClient.AdminId)
		CustomerServiceContainerHandle.CustomerWebsocketCount--
		return serviceClient, nil
	}
	return nil, CustomerServiceNotFoundErr
}
