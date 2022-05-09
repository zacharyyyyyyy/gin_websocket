package model

import (
	"fmt"
	"gin_websocket/lib/config"
	"gin_websocket/lib/logger"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gLog "gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
)

var DbConn *gorm.DB

func init() {
	dbConf := config.BaseConf.GetDbConf()
	dbString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True", dbConf.UsernameMaster, dbConf.PasswordMaster, dbConf.HostMaster, dbConf.PortMaster, dbConf.DatabaseMaster, "utf8")
	slaveDbString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True", dbConf.UsernameSlave, dbConf.PasswordSlave, dbConf.HostSlave, dbConf.PortSlave, dbConf.DatabaseSlave, "utf8")
	var err error
	ormLog := gLog.New(
		log.New(logger.Model, "SQL log:", log.Lshortfile),
		gLog.Config{
			SlowThreshold: time.Second * 5, // 慢 SQL 阈值
			LogLevel:      gLog.Info,       // Log level
			Colorful:      false,           // 禁用彩色打印
		})

	DbConn, err = gorm.Open(mysql.Open(dbString), &gorm.Config{
		Logger: ormLog,
	})
	if err != nil {
		panic(err)
	}
	err = DbConn.Use(dbresolver.Register(
		dbresolver.Config{
			Sources: []gorm.Dialector{
				mysql.Open(slaveDbString),
			},
			Replicas: nil,
			Policy:   dbresolver.RandomPolicy{},
		}),
	)
	sqlDB, _ := DbConn.DB()
	// SetMaxIdleConns 用于设置连接池中空闲连接的最大数量。
	sqlDB.SetMaxIdleConns(dbConf.MaxIdleConns)
	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(dbConf.MaxOpenConns)
	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(time.Duration(dbConf.ConnMaxLifetimeMinus) * time.Minute)
	if err != nil {
		logger.Runtime.Error(err.Error())
	}
}
