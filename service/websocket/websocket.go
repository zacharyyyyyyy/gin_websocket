package websocket

import (
	"context"
	"errors"
	"fmt"
	"gin_websocket/lib/logger"
	"github.com/gin-gonic/gin"
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
	ClientBuildFailErr = errors.New("websocket创建失败")
	ClientNotFoundErr  = errors.New("客户端链接不存在,或已关闭")
	WrongConnErr       = errors.New("该请求非websocket")
	SendMsgErr         = errors.New("发送消息失败")
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
	return Cont.webSocketCont
}

func (Cont *WsContainer) Remove(userClient UserClient) error {
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
	Cont.webSocketCont--
	return nil
}

func (Cont *WsContainer) Send(message Message) error {
	err := Cont.WebSocketClientMap[message.WebsocketKey].send(message)
	if err != nil {
		logger.Service.Error(err.Error())
		return err
	}
	return nil
}
func (Cont *WsContainer) Receive(message []byte) error {
	//TODO
}

//定时释放webcoket
func (Cont *WsContainer) CleanClient(ctx context.Context, timeDuration time.Duration) {
	timer := time.NewTimer(timeDuration)
	defer timer.Stop()
	for {
		select {
		case <-timer.C:
			for _, userClient := range WsContainerHandle.WebSocketClientMap {
				err := userClient.timeout()
				errString := fmt.Sprintf("websocket timeout func err:%s", err)
				logger.Service.Error(errString)
				err = WsContainerHandle.Remove(userClient)
				errString = fmt.Sprintf("websocket remove func err:%s", err)
				logger.Service.Error(errString)
			}
			timer.Reset(timeDuration)
		case <-ctx.Done():
			logger.Service.Info("websocket Cleanclient Func close")
			return
		}
	}
}
func (Cont *WsContainer) append(userClient *UserClient) {
	Cont.lock.Lock()
	defer Cont.lock.Unlock()
	Cont.WebSocketClientMap[userClient.Id] = *userClient
	Cont.webSocketCont++
}
