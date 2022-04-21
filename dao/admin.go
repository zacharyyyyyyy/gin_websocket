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

func GetAllAdminByLimitAndOffset(limit, offset int) (res []*model.AdminWithRole, err error) {
	db := model.DbConn.Table(_adminTable)
	db.Select("admin.id, admin.username, admin.name, admin_role.name as role_name, admin_role.describe, admin.create_time")
	db.Joins("join admin_role on admin_role.id = admin.role")
	if err := db.Limit(limit).Offset(offset).Order("admin.id DESC").Find(&res).Error; err != nil {
		return nil, err
	}
	return
}

func GetAdminCount() (count int64, err error) {
	db := model.DbConn.Table(_adminTable)
	db.Joins("join admin_role on admin_role.id = admin.role")
	if err := db.Count(&count).Error; err != nil {
		return 0, err
	}
	return
}
