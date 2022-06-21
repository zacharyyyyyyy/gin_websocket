package task

import (
	"gin_websocket/lib/kafka"
	"gin_websocket/model"
)

type kafkaHandler struct {
}

func init() {
	registerTask(model.TypeKafka, func() Handler {
		return kafkaHandler{}
	})
}

func (handler kafkaHandler) Exec(param map[string]interface{}) error {
	_, _, err := kafka.KafkaServer.TaskSingleSend(param["topic"].(string), param["key"].(string), param["data"].(map[string]interface{}))
	return err
}
