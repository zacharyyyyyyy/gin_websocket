package dao

import (
	"gin_websocket/model"
	jsoniter "github.com/json-iterator/go"
	"time"
)

const (
	_taskqueueTable = "taskqueue"
)

func SelectMultiByStatusAndLimitAndOffset(status, limit, offset int) (res []*model.Taskqueue, err error) {
	db := model.DbConn.Table(_taskqueueTable)
	if err := db.Where("status = ? AND begin_time < ?", status, time.Now().Unix()).Order("begin_time DESC").Limit(limit).Offset(offset).Find(res).Error; err != nil {
		return nil, err
	}
	return
}

func AddTask(typeString string, param map[string]interface{}, beginTime int) error {
	db := model.DbConn.Table(_taskqueueTable)
	ParamString, _ := jsoniter.Marshal(param)
	saveTask := model.Taskqueue{
		Type:       typeString,
		Param:      string(ParamString),
		CreateTime: int(time.Now().Unix()),
		BeginTime:  beginTime,
	}
	return db.Create(saveTask).Error
}
