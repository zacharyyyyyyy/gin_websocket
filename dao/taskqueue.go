package dao

import (
	"gorm.io/gorm"
	"time"

	"gin_websocket/model"

	jsoniter "github.com/json-iterator/go"
)

const (
	_taskqueueTable = "taskqueue"
)

func SelectMultiByStatusAndLimitAndOffset(status, limit, offset int) (res []*model.Taskqueue, err error) {
	db := model.DbConn.GetMasterDb().Table(_taskqueueTable)
	if err = db.Where("status = ? AND begin_time < ?", status, time.Now().Unix()).Order("begin_time ASC").Limit(limit).Offset(offset).Find(&res).Error; err != nil {
		return nil, err
	}
	return
}

func AddTask(typeString string, param map[string]interface{}, beginTime int) error {
	db := model.DbConn.GetMasterDb().Table(_taskqueueTable)
	ParamString, _ := jsoniter.Marshal(param)
	saveTask := model.Taskqueue{
		Type:       typeString,
		Param:      string(ParamString),
		CreateTime: int(time.Now().Unix()),
		BeginTime:  beginTime,
	}
	return db.Create(&saveTask).Error
}

func DelTask(id int) error {
	db := model.DbConn.GetMasterDb().Table(_taskqueueTable)
	if err := db.Where("id = ?", id).Delete(&model.Taskqueue{}).Error; err != nil {
		return err
	}
	return nil
}

func DelayTask(id int, time time.Time, err error) error {
	db := model.DbConn.GetMasterDb().Table(_taskqueueTable)
	if err := db.Where("id = ?", id).Updates(map[string]interface{}{"status": model.StatusNotBegin, "begin_time": time.Unix(), "retry_times": gorm.Expr("retry_times + ?", 1), "fail_msg": err.Error()}).Error; err != nil {
		return err
	}
	return nil
}
func UpdateStatusToRunning(id int) error {
	db := model.DbConn.GetMasterDb().Table(_taskqueueTable)
	if err := db.Where("id = ?", id).Update("status", model.StatusRunning).Error; err != nil {
		return err
	}
	return nil
}
