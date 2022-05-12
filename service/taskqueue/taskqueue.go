package taskqueue

import (
	"context"
	"errors"
	"gin_websocket/dao"
	"gin_websocket/lib/logger"
	"gin_websocket/model"
	"time"
)

type TaskHandler interface {
	Exec(param map[string]interface{}) error
}

var (
	taskHandlerNotFoundErr = errors.New("消费者未注册")
	timeOutErr             = errors.New("任务超时")
)

var (
	taskEachCount  = 10
	timeDuration   = 2 * time.Second
	taskHandlerMap = make(map[string]func() TaskHandler)
)

func Start() {
	go start()
}

func RegisterTask(taskName string, f func() TaskHandler) {
	taskHandlerMap[taskName] = f
}

func start() {
	timeTicker := time.NewTicker(timeDuration)
	defer timeTicker.Stop()
	for {
		select {
		case <-timeTicker.C:
			taskStructArray, err := dao.SelectMultiByStatusAndLimitAndOffset(model.StatusNotBegin, taskEachCount, 0)
			if err != nil {
				logger.Service.Error(err.Error())
				continue
			}
			if len(taskStructArray) > 0 {
				for _, taskStruct := range taskStructArray {
					handler, err := getTask(taskStruct.Type)
					if err != nil {
						//todo
					}
				}

			}
		}
	}
}

func getTask(taskName string) (TaskHandler, error) {
	if f, ok := taskHandlerMap[taskName]; ok {
		return f(), nil
	}
	return nil, taskHandlerNotFoundErr
}

func run(ctx context.Context, handler TaskHandler, param map[string]interface{}) error {
	done := make(chan struct{}, 1)
	go func() {
		err := handler.Exec(param)
		if err != nil {
			logger.Service.Error(err.Error())
		}
		done <- struct{}{}
	}()
	select {
	case <-ctx.Done():
		return timeOutErr
	case <-done:
		return nil
	}
}

func AddTask(typeString string, param map[string]interface{}, beginTime int) {
	err := dao.AddTask(typeString, param, beginTime)
	if err != nil {
		logger.Service.Error(err.Error())
		return
	}
}

func delTask(id int) {
	_ = dao.DelTask(id)
}

func delayTask(id int) {
	_ = dao.DelayTask(id)
}
