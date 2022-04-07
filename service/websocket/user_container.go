package websocket

import (
	"sync"
)

//用户容器
type WsContainer struct {
	WebSocketClientMap   map[WsKey]*UserClient
	lock                 *sync.RWMutex
	ClientWebSocketCount int
}

//容器加载
func userStart() *WsContainer {
	client := make(map[WsKey]*UserClient, 1)
	WsContainerHandle := &WsContainer{
		WebSocketClientMap:   client,
		lock:                 &sync.RWMutex{},
		ClientWebSocketCount: 0,
	}
	return WsContainerHandle
}

//用户连入初始化
func (Cont *WsContainer) NewClient(userClient *UserClient) error {
	return Cont.append(userClient)
}

func (Cont WsContainer) GetConnCount() int {
	return Cont.ClientWebSocketCount
}

//主动删除
func (Cont *WsContainer) Remove(userClient *UserClient) error {
	if _, ok := Cont.WebSocketClientMap[userClient.Id]; !ok {
		return ClientNotFoundErr
	}
	err := Cont.WebSocketClientMap[userClient.Id].Close()
	return err
}

func (Cont *WsContainer) append(userClient *UserClient) error {
	Cont.lock.Lock()
	defer Cont.lock.Unlock()
	if Cont.ClientWebSocketCount > wsConf.MaxConnection {
		return TooManyConnectionErr
	}
	if _, ok := Cont.WebSocketClientMap[userClient.Id]; !ok {
		Cont.WebSocketClientMap[userClient.Id] = userClient
		Cont.ClientWebSocketCount++
		return nil
	}
	return ClientAlreadyInContainer
}

//链接断时删除
func (Cont *WsContainer) remove(userClient *UserClient) error {
	Cont.lock.Lock()
	defer Cont.lock.Unlock()
	if _, ok := Cont.WebSocketClientMap[userClient.Id]; !ok {
		return ClientNotFoundErr
	}
	delete(Cont.WebSocketClientMap, userClient.Id)
	Cont.ClientWebSocketCount--
	return nil
}
