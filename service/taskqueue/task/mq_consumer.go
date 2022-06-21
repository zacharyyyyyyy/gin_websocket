package task

import (
	"gin_websocket/lib/mq"
	"gin_websocket/model"
)

type mqHandler struct {
}

func init() {
	registerTask(model.TypeMq, func() Handler {
		return mqHandler{}
	})
}

func (handler mqHandler) Exec(param map[string]interface{}) error {
	err := mq.MqServer.TaskSingleSend(param["data"].(map[string]interface{}), param["qKey"].(string))
	return err
}
