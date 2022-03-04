package websocket

import (
	"errors"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
)

type WsContainer struct {
	WebSocketClientMap map[WsKey]UserClient
	lock               *sync.Mutex
	WebSocketCont      int
}

var (
	WsContainerHandle *WsContainer
	upGrader          = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

var (
	ClientNotFoundErr = errors.New("客户端链接不存在")
	WrongConnErr      = errors.New("该请求非websocket")
	SendMsgErr        = errors.New("发送消息失败")
)

func init() {
	client := make(map[WsKey]UserClient, 0)
	WsContainerHandle = &WsContainer{
		WebSocketClientMap: client,
		lock:               &sync.Mutex{},
		WebSocketCont:      0,
	}
}

func (Cont *WsContainer) Append(userClient UserClient) {
	Cont.lock.Lock()
	defer Cont.lock.Unlock()
	Cont.WebSocketClientMap[userClient.Id] = userClient
	Cont.WebSocketCont++
}
func (Cont *WsContainer) Remove(userClient UserClient) error {
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
	Cont.WebSocketCont--
	return nil
}

//定时释放webcoket
func CleanClient() {
	//TODO
}
