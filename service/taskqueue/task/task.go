package task

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gin_websocket/dao"
	"gin_websocket/lib/logger"
	"gin_websocket/model"

	jsoniter "github.com/json-iterator/go"
)

type Handler interface {
	Exec(param map[string]interface{}) error
}

type Task struct {
	taskId      int
	TaskHandler Handler
	param       map[string]interface{}
}

var (
	taskHandlerNotFoundErr = errors.New("消费者未注册")
	timeOutErr             = errors.New("任务超时")
)

var (
	taskEachCount  = 10
	timeDuration   = 3 * time.Second
	eachTaskTime   = 10 * time.Second
	taskHandlerMap = make(map[string]func() Handler)
)

func Start() {
	go start()
}

func registerTask(taskName string, f func() Handler) {
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
				logger.TaskQueue.Error(err.Error())
				continue
			}
			if len(taskStructArray) <= 0 {
				continue
			}
			for _, taskStruct := range taskStructArray {
				handler, err := getTask(*taskStruct)
				if err != nil {
					wrapErr := fmt.Errorf("%w(id:%d,type:%s)", err, taskStruct.Id, taskStruct.Type)
					logger.TaskQueue.Error(wrapErr.Error())
				} else {
					runtimeErr := handler.run()
					if runtimeErr != nil {
						wrapErr := fmt.Errorf("%w(id:%d)", runtimeErr, taskStruct.Id)
						logger.TaskQueue.Error(wrapErr.Error())
					}
				}
			}
		}
	}
}

func getTask(taskRow model.Taskqueue) (Task, error) {
	if f, ok := taskHandlerMap[taskRow.Type]; ok {
		var param = make(map[string]interface{})
		_ = jsoniter.Unmarshal([]byte(taskRow.Param), &param)
		return Task{
			taskId:      taskRow.Id,
			TaskHandler: f(),
			param:       param,
		}, nil
	}
	return Task{}, taskHandlerNotFoundErr
}

func (task Task) run() error {
	var runtimeErr error
	ctx, cancel := context.WithTimeout(context.Background(), eachTaskTime)
	done := make(chan struct{}, 1)
	go func() {
		err := task.TaskHandler.Exec(task.param)
		if err != nil {
			runtimeErr = err
			task.delayTask()
		} else {
			task.delTask()
		}
		done <- struct{}{}
	}()
	select {
	case <-ctx.Done():
		return timeOutErr
	case <-done:
		cancel()
		return runtimeErr
	}
}

func (task Task) delTask() {
	_ = dao.DelTask(task.taskId)
}

func (task Task) delayTask() {
	_ = dao.DelayTask(task.taskId)
}

func (task Task) runningTask() {
	_ = dao.UpdateStatus(task.taskId)
}
