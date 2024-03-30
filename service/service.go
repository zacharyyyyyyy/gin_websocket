package service

import (
	"context"
	"fmt"

	"gin_websocket/lib/logger"
	"gin_websocket/service/taskqueue/task"
	"gin_websocket/service/websocket"
)

type serviceHandle interface {
	Start(serviceCtx context.Context)
	Stop() <-chan struct{}
}

var (
	serviceCtx, serviceCancelFunc = context.WithCancel(context.Background())
	serviceHandlerSlice           = make([]serviceHandle, 0)
)

func init() {
	serviceHandlerSlice = append(serviceHandlerSlice, &websocket.WsStruct{}, &task.TaskQueueStruct{})
	fmt.Println(serviceHandlerSlice)
}

func Setup() {
	for _, handler := range serviceHandlerSlice {
		go handler.Start(serviceCtx)
	}
}

func Stop(ctx context.Context) {
	stopChan := make(chan struct{}, 1)
	go func() {
		serviceCancelFunc()
		serviceMapLen := len(serviceHandlerSlice)
		if serviceMapLen == 0 {
			return
		}
		for _, handler := range serviceHandlerSlice {
			<-handler.Stop()
		}
		stopChan <- struct{}{}
	}()
	select {
	case <-stopChan:
		logger.Service.Info("service stop by stopChan")
	case <-ctx.Done():
		logger.Service.Info("service stop timeout")
	}
}
