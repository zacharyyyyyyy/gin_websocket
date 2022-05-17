package dao

import (
	"gin_websocket/model"
)

const (
	_adminAuthTable = "admin_auth"
)

func GetAllAuthByEnable() (res []*model.AdminAuth, err error) {
	db := model.DbConn.GetSlaveDb().Table(_adminAuthTable)
	if err = db.Where("enable = ?", 1).Order("id ASC").Find(&res).Error; err != nil {
		return nil, err
	}
	return
}

func GetAuthByLimitAndOffset(limit, offset int) (res []*model.AdminAuth, err error) {
	db := model.DbConn.GetSlaveDb().Table(_adminAuthTable)
	if err = db.Where("enable = ?", 1).Limit(limit).Offset(offset).Order("id ASC").Find(&res).Error; err != nil {
		return nil, err
	}
	return
}

func GetAllAuthCountByEnable() (count int64, err error) {
	db := model.DbConn.GetSlaveDb().Table(_adminAuthTable)
	if err = db.Where("enable = ?", 1).Count(&count).Error; err != nil {
		return 0, err
	}
	return
}
