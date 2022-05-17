package dao

import (
	"time"

	"gin_websocket/model"
)

const (
	_adminAuthMapTable = "admin_auth_map"
)

func GetRoleByAuth(auth int) (res []*model.AdminAuthMap, err error) {
	db := model.DbConn.GetSlaveDb().Table(_adminAuthMapTable)
	if err = db.Where("auth = ?", auth).Find(&res).Error; err != nil {
		return nil, err
	}
	return
}

func GetAllAuthMapByRole(limit, offset, role int) (res []*model.AdminAuthMapDetail, err error) {
	db := model.DbConn.GetSlaveDb().Table(_adminAuthMapTable)
	db.Joins("join admin_auth on admin_auth.id = admin_auth_map.auth").Joins("join admin_role on admin_role.id = admin_auth_map.role")
	db.Select("admin_auth_map.role, admin_auth_map.auth, admin_role.name as role_name, admin_role.describe as role_describe, admin_auth.name as auth_name")
	if err = db.Where("role = ?", role).Limit(limit).Offset(offset).Find(res).Error; err != nil {
		return nil, err
	}
	return
}

func GetAuthMapCount() (count int64, err error) {
	db := model.DbConn.GetSlaveDb().Table(_adminAuthMapTable)
	if err = db.Count(&count).Error; err != nil {
		return 0, err
	}
	return
}

func AddAuth(role, auth int) error {
	db := model.DbConn.GetMasterDb().Table(_adminAuthMapTable)
	authMap := model.AdminAuthMap{
		Role:       role,
		Auth:       auth,
		CreateTime: int(time.Now().Unix()),
	}
	if err := db.Create(authMap).Error; err != nil {
		return err
	}
	return nil
}

func EditAuth(role, auth int) error {
	db := model.DbConn.GetMasterDb().Table(_adminAuthMapTable)
	authMap := model.AdminAuthMap{
		Role: role,
		Auth: auth,
	}
	if err := db.Save(authMap).Error; err != nil {
		return err
	}
	return nil
}

func DelAuth(role, auth int) error {
	db := model.DbConn.GetMasterDb().Table(_adminAuthMapTable)
	authMap := model.AdminAuthMap{
		Role: role,
		Auth: auth,
	}
	if err := db.Delete(authMap).Error; err != nil {
		return err
	}
	return nil
}
