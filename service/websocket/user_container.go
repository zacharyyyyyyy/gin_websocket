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

func (Cont *WsContainer) Remove(userClient *UserClient) error {
	Cont.lock.Lock()
	defer Cont.lock.Unlock()
	if _, ok := Cont.WebSocketClientMap[userClient.Id]; !ok {
		return ClientNotFoundErr
	}
	//先释放链接
	err := Cont.WebSocketClientMap[userClient.Id].Close()
	if err != nil {
		return err
	}
	delete(Cont.WebSocketClientMap, userClient.Id)
	Cont.ClientWebSocketCount--
	return nil
}

func (Cont *WsContainer) append(userClient *UserClient) error {
	Cont.lock.Lock()
	defer Cont.lock.Unlock()
	if Cont.ClientWebSocketCount > wsConf.MaxConnection {
		return TooManyConnectionErr
	}
	Cont.WebSocketClientMap[userClient.Id] = userClient
	Cont.ClientWebSocketCount++
	return nil
}
