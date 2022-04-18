package dao

import "gin_websocket/model"

const (
	_adminTable = "admin"
)

func SelectOneByUsername(username string) (res *model.Admin, err error) {
	db := model.DbConn.Table(_adminTable)
	if err := db.Where("username = ?", username).Limit(1).Find(&res).Error; err != nil {
		return nil, err
	}
	return
}
