package taskqueue

import (
	"context"
	"errors"
	"fmt"
	"gin_websocket/dao"
	"gin_websocket/lib/logger"
	"gin_websocket/model"
	"time"
)

type TaskHandler interface {
	Exec(param map[string]interface{}) error
}

type Task struct {
	taskId      int
	TaskHandler TaskHandler
	param       map[string]interface{}
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

func RegisterTask(taskName string, f func() TaskHandler) {
	taskHandlerMap[taskName] = f
}

func Start() {
	go start()
}

func start() {
	timeTicker := time.NewTicker(timeDuration)
	defer timeTicker.Stop()
	for {
		select {
		case <-timeTicker.C:
			taskStructArray, err := dao.SelectMultiByStatusAndLimitAndOffset(model.StatusNotBegin, taskEachCount, 0)
			if err != nil {
				logger.TaskQueue.Error(err.Error())
				continue
			}
			if len(taskStructArray) > 0 {
				for _, taskStruct := range taskStructArray {
					handler, err := getTask(taskStruct.Type)
					if err != nil {
						//task := Task{
						//	taskId:      taskStruct.Id,
						//	TaskHandler: nil,
						//	param:       ,
						//}
					} else {
						wrapErr := fmt.Errorf("%w:%s", err, taskStruct.Type)
						logger.TaskQueue.Error(wrapErr.Error())
					}
				}

			}
		}
	}
}

func getTask(taskName string) (Task, error) {
	if f, ok := taskHandlerMap[taskName]; ok {
		return Task{}, nil
	}
	return Task{}, taskHandlerNotFoundErr
}

func (task Task) run(ctx context.Context) error {
	done := make(chan struct{}, 1)
	go func() {
		err := task.TaskHandler.Exec(task.param)
		if err != nil {
			logger.TaskQueue.Error(err.Error())
			task.delayTask()
		} else {
			task.delayTask()
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

func (task Task) delTask() {
	_ = dao.DelTask(task.taskId)
}

func (task Task) delayTask() {
	_ = dao.DelayTask(task.taskId)
}

func AddTask(typeString string, param map[string]interface{}, beginTime int) {
	err := dao.AddTask(typeString, param, beginTime)
	if err != nil {
		logger.Service.Error(err.Error())
		return
	}
}
