package dao

import "gin_websocket/model"

const (
	_adminAuthMapTable = "admin_auth_map"
)

func GetRoleByAuth(auth int) (res []*model.AdminAuthMap, err error) {
	db := model.DbConn.Table(_adminAuthMapTable)
	if err = db.Where("auth = ?", auth).Find(&res).Error; err != nil {
		return nil, err
	}
	return
}
