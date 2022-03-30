package websocket

import (
	"context"
	"gin_websocket/lib/logger"
	"github.com/gin-gonic/gin"
	"sync"
)

//用户容器
type WsContainer struct {
	WebSocketClientMap   map[WsKey]*UserClient
	lock                 *sync.RWMutex
	ClientWebSocketCount int
}

//容器加载
func UserStart() *WsContainer {
	client := make(map[WsKey]*UserClient, 1)
	WsContainerHandle := &WsContainer{
		WebSocketClientMap:   client,
		lock:                 &sync.RWMutex{},
		ClientWebSocketCount: 0,
	}
	return WsContainerHandle
}

//用户连入初始化
func (Cont *WsContainer) NewClient(ctx context.Context, c *gin.Context) error {
	var userClient *UserClient
	userClient, err := newUser(ctx, c)
	if err != nil {
		logger.Service.Error(err.Error())
		return err
	}
	Cont.append(userClient)
	return nil
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
	err := Cont.WebSocketClientMap[userClient.Id].close()
	if err != nil {
		return err
	}
	delete(Cont.WebSocketClientMap, userClient.Id)
	Cont.ClientWebSocketCount--
	return nil
}

func (Cont *WsContainer) append(userClient *UserClient) {
	Cont.lock.Lock()
	defer Cont.lock.Unlock()
	Cont.WebSocketClientMap[userClient.Id] = userClient
	Cont.ClientWebSocketCount++
}
