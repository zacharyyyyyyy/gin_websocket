package consumer

import "gin_websocket/service/taskqueue"

type mqhandler struct {
}

var TypeMq = "mq"

func init() {
	taskqueue.RegisterTask(TypeMq, func() taskqueue.TaskHandler {
		return mqhandler{}
	})
}

func (handler mqhandler) Exec(param map[string]interface{}) error {
	//todo
}
