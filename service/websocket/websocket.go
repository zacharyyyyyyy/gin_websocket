package websocket

import (
	"context"
	"errors"
	"fmt"
	"gin_websocket/lib/config"
	"gin_websocket/lib/logger"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

type WsKey string

var (
	WsContainerHandle              = UserStart()
	CustomerServiceContainerHandle = ServiceStart()
	upGrader                       = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

var (
	ClientBuildFailErr    = errors.New("websocket创建失败")
	ClientNotFoundErr     = errors.New("客户端链接不存在,或已关闭")
	WrongConnErr          = errors.New("该请求非websocket")
	SendMsgErr            = errors.New("发送消息失败")
	ClientAlreadyBoundErr = errors.New("客户端已被绑定")
)

var WsConf config.WebsocketConf = config.BaseConf.GetWsConf()

func Start() {
	ctx, _ := context.WithCancel(context.Background())
	limitTime := time.Duration(WsConf.CleanLimitTimeSec) * time.Second
	go cleanClient(ctx, WsContainerHandle, limitTime)
}

//定时释放webcoket
func cleanClient(ctx context.Context, Cont *WsContainer, timeDuration time.Duration) {
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
