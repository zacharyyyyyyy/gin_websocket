package dao

import (
	"errors"
	"net/http"
	"time"

	"gin_websocket/lib/session"
	"gin_websocket/model"
)

const (
	_adminTable = "admin"
)

func SelectOneByUsername(username string) (res *model.Admin, err error) {
	db := model.DbConn.GetSlaveDb().Table(_adminTable)
	if err := db.Where("username = ?", username).Limit(1).Find(&res).Error; err != nil {
		return nil, err
	}
	return
}

func GetAllAdminByLimitAndOffset(limit, offset int) (res []*model.AdminWithRole, err error) {
	db := model.DbConn.GetSlaveDb().Table(_adminTable)
	db.Select("admin.id, admin.username, admin.name, admin_role.name as role_name, admin_role.describe, admin.create_time")
	db.Joins("join admin_role on admin_role.id = admin.role")
	if err := db.Limit(limit).Offset(offset).Order("admin.id DESC").Find(&res).Error; err != nil {
		return nil, err
	}
	return
}

func GetAdminCount() (count int64, err error) {
	db := model.DbConn.GetSlaveDb().Table(_adminTable)
	db.Joins("join admin_role on admin_role.id = admin.role")
	if err := db.Count(&count).Error; err != nil {
		return 0, err
	}
	return
}

func AddAdmin(username, name, password string, role int) error {
	db := model.DbConn.GetMasterDb().Table(_adminTable)
	saveAdmin := model.Admin{
		Username:   username,
		Password:   password,
		Name:       name,
		Role:       role,
		CreateTime: int(time.Now().Unix()),
	}
	return db.Create(saveAdmin).Error
}

func EditAdmin(username, name, password string, role, id int) error {
	db := model.DbConn.GetMasterDb().Table(_adminTable)
	saveAdmin := model.Admin{
		Id:       id,
		Username: username,
		Password: password,
		Name:     name,
		Role:     role,
	}
	if err := db.Save(saveAdmin).Error; err != nil {
		return err
	}
	return nil
}

func DelAdmin(id int) error {
	db := model.DbConn.GetMasterDb().Table(_adminTable)
	if err := db.Where("id = ?", id).Delete(model.Admin{}).Error; err != nil {
		return err
	}
	return nil
}

func GetCurrent(cRequest *http.Request) (res *model.Admin, err error) {
	sessionStruct, err := session.GetCurrent(cRequest)
	if err != nil {
		return nil, errors.New("未登录")
	}
	adminId, err := sessionStruct.GetString("admin")
	if err != nil {
		return nil, errors.New("未登录")
	}
	db := model.DbConn.GetSlaveDb().Table(_adminTable)
	if err = db.Where("id = ?", adminId).Limit(1).Find(res).Error; err != nil {
		return nil, err
	}
	return
}
