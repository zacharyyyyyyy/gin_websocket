package dao

import (
	"time"

	"gin_websocket/model"
)

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

func AddAuth(role, auth int) error {
	db := model.DbConn.Table(_adminAuthMapTable)
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
	db := model.DbConn.Table(_adminAuthMapTable)
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
	db := model.DbConn.Table(_adminAuthMapTable)
	authMap := model.AdminAuthMap{
		Role: role,
		Auth: auth,
	}
	if err := db.Delete(authMap).Error; err != nil {
		return err
	}
	return nil
}
