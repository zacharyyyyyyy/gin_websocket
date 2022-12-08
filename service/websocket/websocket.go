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
	WsContainerHandle              = userStart()
	CustomerServiceContainerHandle = serviceStart()
	upGrader                       = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

var (
	ClientBuildFailErr         = errors.New("websocket创建失败")
	ClientNotFoundErr          = errors.New("客户端链接不存在,或已关闭")
	ClientAlreadyInContainer   = errors.New("客户端已在容器内")
	WrongConnErr               = errors.New("该请求非websocket")
	SendMsgErr                 = errors.New("发送消息失败")
	ClientAlreadyBoundErr      = errors.New("客户端已被绑定")
	CustomerServiceNotFoundErr = errors.New("暂无客服！请耐心等候")
	TooManyConnectionErr       = errors.New("服务忙碌中，请稍后重试")
	CloseErr                   = errors.New("链接已关闭")
)

var wsConf config.WebsocketConf = config.BaseConf.GetWsConf()

func Start() {
	ctx, _ := context.WithCancel(context.Background())
	limitTime := time.Duration(wsConf.CleanLimitTimeSec) * time.Second
	go cleanClient(ctx, WsContainerHandle, limitTime)
}

//定时释放webcoket
func cleanClient(ctx context.Context, Cont *WsContainer, timeDuration time.Duration) {
	timer := time.NewTicker(timeDuration)
	defer timer.Stop()
	for {
		select {
		case <-timer.C:
			for _, userClient := range Cont.WebSocketClientMap {
				ip := userClient.Ip
				closeStatus, err := userClient.timeout()
				if err != nil {
					errString := fmt.Sprintf("websocket timeout func err:%s", err)
					logger.Service.Error(errString)
					continue
				}
				if closeStatus {
					err = Cont.Remove(userClient)
					if err != nil {
						errString := fmt.Sprintf("websocket remove func err:%s", err)
						logger.Service.Error(errString)
						continue
					}
				}
				logger.Service.Info(fmt.Sprintf("websocket clear, ip:%s", ip))
			}
		case <-ctx.Done():
			logger.Service.Info("websocket Cleanclient Func close")
			return
		}
	}
}
