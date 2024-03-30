package task

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/sync/semaphore"
	"sync"
	"time"

	"gin_websocket/dao"
	"gin_websocket/lib/logger"
	"gin_websocket/lib/tools"
	"gin_websocket/model"

	jsoniter "github.com/json-iterator/go"
)

type Handler interface {
	Exec(param map[string]interface{}) error
}

type (
	Task struct {
		taskId      int
		TaskHandler Handler
		param       map[string]interface{}
	}
	TaskQueueStruct struct {
	}
)

var (
	taskQueueErr           = errors.New("taskqueue err")
	taskHandlerNotFoundErr = fmt.Errorf("%w:消费者未注册", taskQueueErr)
	timeOutErr             = fmt.Errorf("%w:任务超时", taskQueueErr)
	semaphoreFullErr       = fmt.Errorf("%w:taskqueue限制的gouroutine数已满", taskQueueErr)
	taskRunningErr         = fmt.Errorf("%w:任务正在运行", taskQueueErr)
	taskFailErr            = fmt.Errorf("%w:任务发生panic", taskQueueErr)
)

var (
	taskEachCount  = 10
	timeDuration   = 3 * time.Second
	eachTaskTime   = 30 * time.Second
	taskTimeDelay  = 180 * time.Second
	taskHandlerMap = make(map[string]func() Handler)
	//限制任务goroutine数
	taskGoroutineMaxCount int64 = 100
	taskGoroutineEach     int64 = 1
	sema                        = semaphore.NewWeighted(taskGoroutineMaxCount)
	taskRegisterLock      sync.Mutex
	TaskStopChan          = make(chan struct{}, 1)
)

func (ts TaskQueueStruct) Start(ctx context.Context) {
	go start(ctx)
}

func (ts TaskQueueStruct) Stop() <-chan struct{} {
	return TaskStopChan
}

func registerTask(taskName string, f func() Handler) {
	taskRegisterLock.Lock()
	defer taskRegisterLock.Unlock()
	taskHandlerMap[taskName] = f
}

func start(ctx context.Context) {
	logger.Service.Info("taskqueue Func start")
	timeTicker := time.NewTicker(timeDuration)
	defer timeTicker.Stop()
	defer func() {
		logger.Service.Info("taskqueue Func stop")
	}()
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
					if !sema.TryAcquire(taskGoroutineEach) {
						logger.TaskQueue.Error(semaphoreFullErr.Error())
						continue
					}
					//获取到信号则判断task是否运行
					err = handler.runningTask()
					if err != nil {
						wrapErr := fmt.Errorf("%w(id:%d,type:%s)", err, taskStruct.Id, taskStruct.Type)
						logger.TaskQueue.Error(wrapErr.Error())
						//任务已运行 释放信号
						sema.Release(taskGoroutineEach)
						continue
					}
					go func(taskHandler Task, taskId int) {
						runtimeErr := taskHandler.run()
						if runtimeErr != nil {
							wrapErr := fmt.Errorf("%w(id:%d)", runtimeErr, taskId)
							logger.TaskQueue.Error(wrapErr.Error())
						}
					}(handler, taskStruct.Id)
				}
			}
		case <-ctx.Done():
			timeTicker.Stop()
			now := time.Now()
			for {
				//任务超时 或 无任务时 退出
				if sema.TryAcquire(taskGoroutineMaxCount) || time.Now().Unix()-now.Unix() >= int64(eachTaskTime) {
					break
				}

				time.Sleep(1 * time.Second)
			}
			TaskStopChan <- struct{}{}
			return
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
		//发生panic 默认任务失败，有事务需回滚，延迟一次任务
		defer func() {
			tools.RecoverFunc()
			task.delayTask(time.Now().Add(taskTimeDelay), taskFailErr)
		}()
		err := task.TaskHandler.Exec(task.param)
		if err != nil {
			runtimeErr = err
			task.delayTask(time.Now().Add(taskTimeDelay), err)
		} else {
			task.delTask()
		}
		//释放信号量
		sema.Release(taskGoroutineEach)
		done <- struct{}{}
	}()
	select {
	case <-ctx.Done():
		cancel()
		return timeOutErr
	case <-done:
		cancel()
		return runtimeErr
	}
}

func (task Task) delTask() {
	_ = dao.DelTask(task.taskId)
}

func (task Task) delayTask(time time.Time, err error) {
	_ = dao.DelayTask(task.taskId, time, err)
}

func (task Task) runningTask() error {
	err := dao.UpdateStatusToRunning(task.taskId)
	if err != nil {
		return taskRunningErr
	}
	return nil
}
