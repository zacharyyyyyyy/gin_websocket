package websocket

import (
	"sync"

	"gin_websocket/lib/tools"
)

//客服容器
type CustomerServiceContainer struct {
	WebsocketCustomerServiceMap map[int]*CustomerServiceClient
	lock                        *sync.RWMutex
	customerWebsocketCount      uint
}

//容器加载
func serviceStart() *CustomerServiceContainer {
	customerService := make(map[int]*CustomerServiceClient, 0)
	CustomerServiceContainerHandle := &CustomerServiceContainer{
		WebsocketCustomerServiceMap: customerService,
		lock:                        &sync.RWMutex{},
		customerWebsocketCount:      0,
	}
	return CustomerServiceContainerHandle
}

//客服连入初始化
func (Cont *CustomerServiceContainer) NewClient(customerServiceClient *CustomerServiceClient) {
	Cont.append(customerServiceClient)
}

func (Cont CustomerServiceContainer) GetConnCount() uint {
	return Cont.customerWebsocketCount
}

//主动删除
func (Cont *CustomerServiceContainer) Remove(adminId int) error {
	if _, ok := Cont.WebsocketCustomerServiceMap[adminId]; !ok {
		return ClientNotFoundErr
	}
	err := Cont.WebsocketCustomerServiceMap[adminId].close()
	return err
}

func (Cont *CustomerServiceContainer) GetCustomerService(adminId int) (*CustomerServiceClient, error) {
	customerService, ok := Cont.WebsocketCustomerServiceMap[adminId]
	if !ok {
		return nil, ClientNotFoundErr
	}
	return customerService, nil
}

func (Cont *CustomerServiceContainer) append(customerServiceClient *CustomerServiceClient) {
	Cont.lock.Lock()
	defer Cont.lock.Unlock()
	//再次append则仅做更新
	if _, ok := Cont.WebsocketCustomerServiceMap[customerServiceClient.AdminId]; !ok {
		Cont.customerWebsocketCount++
	}
	Cont.WebsocketCustomerServiceMap[customerServiceClient.AdminId] = customerServiceClient
}

func (Cont *CustomerServiceContainer) remove(customerServiceClient *CustomerServiceClient) error {
	Cont.lock.Lock()
	defer Cont.lock.Unlock()
	if _, ok := Cont.WebsocketCustomerServiceMap[customerServiceClient.AdminId]; !ok {
		return ClientNotFoundErr
	}
	delete(Cont.WebsocketCustomerServiceMap, customerServiceClient.AdminId)
	Cont.customerWebsocketCount--
	return nil
}

func getCustomerService() (*CustomerServiceClient, error) {
	CustomerServiceContainerHandle.lock.Lock()
	defer CustomerServiceContainerHandle.lock.Unlock()
	var customerServer *CustomerServiceClient
	if CustomerServiceContainerHandle.customerWebsocketCount == 0 {
		return nil, CustomerServiceNotFoundErr
	}
	randKey := tools.Rand(0, int(CustomerServiceContainerHandle.customerWebsocketCount))
	index := 0
	//双重随机
	for _, customer := range CustomerServiceContainerHandle.WebsocketCustomerServiceMap {
		if index == randKey {
			customerServer = customer
			break
		}
		index++
	}
	return customerServer, nil
}
