package dao

import "gin_websocket/model"

const (
	_adminRoleTable = "admin_role"
)

func ExistsRole(roleId int) bool {
	db := model.DbConn.Table(_adminRoleTable)
	var count int64
	if err := db.Where("id = ?", roleId).Count(&count).Error; err != nil {
		return false
	}
	if count > 0 {
		return true
	}
	return false
}
