package service

import (
	"gin_websocket/service/taskqueue/task"
	"gin_websocket/service/websocket"
)

func Setup() {
	websocket.Start()
	task.Start()
}
