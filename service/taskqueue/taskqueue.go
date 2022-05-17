package taskqueue

import (
	"gin_websocket/dao"
	"gin_websocket/lib/logger"
)

func AddTask(typeString string, param map[string]interface{}, beginTime int) {
	err := dao.AddTask(typeString, param, beginTime)
	if err != nil {
		logger.Service.Error(err.Error())
		return
	}
}
