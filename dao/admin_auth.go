package dao

import "gin_websocket/model"

const (
	_adminAuthTable = "admin_auth"
)

func GetAllAuthByEnable() (res []*model.AdminAuth, err error) {
	db := model.DbConn.Table(_adminAuthTable)
	if err = db.Where("enable = ?", 1).Order("id ASC").Find(&res).Error; err != nil {
		return nil, err
	}
	return
}
