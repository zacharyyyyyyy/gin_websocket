package dao

import (
	"time"

	"gin_websocket/model"
)

const (
	_adminRoleTable = "admin_role"
)

func ExistsRole(roleId int) bool {
	db := model.DbConn.GetSlaveDb().Table(_adminRoleTable)
	var count int64
	if err := db.Where("id = ?", roleId).Count(&count).Error; err != nil {
		return false
	}
	if count > 0 {
		return true
	}
	return false
}

func GetAllRole() (res []*model.AdminRole, err error) {
	db := model.DbConn.GetSlaveDb().Table(_adminRoleTable)
	if err := db.Order("id ASC").Find(res).Error; err != nil {
		return nil, err
	}
	return
}

func AddRole(name, describe string) error {
	db := model.DbConn.GetMasterDb().Table(_adminRoleTable)
	role := model.AdminRole{}
	role.Name = name
	role.CreateTime = int(time.Now().Unix())
	if describe != "" {
		role.Describe = describe
	}
	if err := db.Create(role).Error; err != nil {
		return err
	}
	return nil
}

func EditRole(name, describe string, id int) error {
	db := model.DbConn.GetMasterDb().Table(_adminRoleTable)
	role := model.AdminRole{
		Id:         id,
		Name:       name,
		Describe:   describe,
		CreateTime: int(time.Now().Unix()),
	}
	if err := db.Save(role).Error; err != nil {
		return err
	}
	return nil
}

func DelRole(id int) error {
	db := model.DbConn.GetMasterDb().Table(_adminRoleTable)
	if err := db.Where("id = ?", id).Delete(model.AdminRole{}).Error; err != nil {
		return err
	}
	return nil
}
