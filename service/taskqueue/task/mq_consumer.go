package task

import (
	"gin_websocket/lib/mq"
	"gin_websocket/model"
)

type mqhandler struct {
}

func init() {
	registerTask(model.TypeMq, func() Handler {
		return mqhandler{}
	})
}

func (handler mqhandler) Exec(param map[string]interface{}) error {
	err := mq.MqServer.TaskSingleSend(mq.SendMap(param), mq.QueueKeySms)
	return err
}
