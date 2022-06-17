package model

import (
	"errors"
	"fmt"
	"gin_websocket/lib/config"
	"gin_websocket/lib/logger"
	"log"
	"math/rand"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gLog "gorm.io/gorm/logger"
	//"gorm.io/plugin/dbresolver"
)

type dbConn struct {
	masterDb *gorm.DB
	slaveDb  []*gorm.DB
}

var DbConn dbConn

var (
	DbNotFoundErr = errors.New("数据库未配置")
)

func init() {
	dbConf := config.BaseConf.GetDbConf()
	dbString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True", dbConf.UsernameMaster, dbConf.PasswordMaster, dbConf.HostMaster, dbConf.PortMaster, dbConf.DatabaseMaster, "utf8")
	slaveDbString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True", dbConf.UsernameSlave, dbConf.PasswordSlave, dbConf.HostSlave, dbConf.PortSlave, dbConf.DatabaseSlave, "utf8")
	DbConn.masterDb = newDb(dbConf, dbString)
	DbConn.slaveDb = append(DbConn.slaveDb, newDb(dbConf, slaveDbString))
}

func (DbConn dbConn) GetMasterDb() *gorm.DB {
	if DbConn.masterDb == nil {
		logger.Runtime.Error(fmt.Errorf("主%s", DbNotFoundErr).Error())
		return nil
	}
	return DbConn.masterDb
}

func (DbConn dbConn) GetSlaveDb() *gorm.DB {
	length := len(DbConn.slaveDb)
	if length == 0 {
		logger.Runtime.Error(fmt.Errorf("从%s", DbNotFoundErr).Error())
		return nil
	}
	rand.Seed(time.Now().UnixNano())
	randKey := rand.Intn(length)
	return DbConn.slaveDb[randKey]
}

func newDb(dbConf config.DbConf, dbString string) *gorm.DB {
	var err error
	ormLog := gLog.New(
		log.New(logger.Model, "SQL log:", log.Lshortfile),
		gLog.Config{
			SlowThreshold: time.Second * 5, // 慢 SQL 阈值
			LogLevel:      gLog.Warn,       // Log level
			Colorful:      false,           // 禁用彩色打印
		})

	dbConn, err := gorm.Open(mysql.Open(dbString), &gorm.Config{
		Logger: ormLog,
	})
	if err != nil {
		panic(err)
	}
	//err = DbConn.Use(dbresolver.Register(
	//	dbresolver.Config{
	//		Sources: []gorm.Dialector{
	//			mysql.Open(dbString),
	//		},
	//		Replicas: []gorm.Dialector{
	//			mysql.Open(slaveDbString),
	//		},
	//		Policy: dbresolver.RandomPolicy{},
	//	}),
	//)
	sqlDB, _ := dbConn.DB()
	// SetMaxIdleConns 用于设置连接池中空闲连接的最大数量。
	sqlDB.SetMaxIdleConns(dbConf.MaxIdleConns)
	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(dbConf.MaxOpenConns)
	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(time.Duration(dbConf.ConnMaxLifetimeMinus) * time.Minute)
	return dbConn
}
