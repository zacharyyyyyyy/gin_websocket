package websocket

import (
	"context"
	"errors"
	"fmt"
	"gin_websocket/Lib/Logger"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type WsContainer struct {
	WebSocketClientMap map[WsKey]UserClient
	lock               *sync.RWMutex
	webSocketCont      int
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

var (
	pingLastTimeSec = 60
	chatLastTimeSec = 180
)

func Start(ctx context.Context) *WsContainer {
	client := make(map[WsKey]UserClient, 0)
	WsContainerHandle = &WsContainer{
		WebSocketClientMap: client,
		lock:               &sync.RWMutex{},
		webSocketCont:      0,
	}
	return WsContainerHandle
}

func (Cont WsContainer) GetConnCount() int {
	return Cont.webSocketCont
}

func (Cont *WsContainer) Append(userClient UserClient) {
	Cont.lock.Lock()
	defer Cont.lock.Unlock()
	Cont.WebSocketClientMap[userClient.Id] = userClient
	Cont.webSocketCont++
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
	Cont.webSocketCont--
	return nil
}

//定时释放webcoket
func (Cont *WsContainer) CleanClient(ctx context.Context, timeDuration time.Duration) {

	timeTicker := time.NewTicker(timeDuration)
	defer timeTicker.Stop()
	for {
		select {
		case <-timeTicker.C:
			for _, userClient := range WsContainerHandle.WebSocketClientMap {
				err := userClient.timeout()
				errString := fmt.Sprintf("websocket timeout func err:%s", err)
				Logger.Service.Error(errString)
			}
			timeTicker.Reset(timeDuration)
		case <-ctx.Done():
			Logger.Service.Info("websocket Cleanclient Func close")
			return
		}
	}
}
