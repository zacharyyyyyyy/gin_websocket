package config

import (
	"errors"
	"fmt"
	"gin_websocket/lib/logger"
	"github.com/go-ini/ini"
	"sync"
)

type ConfMap interface {
	getPath() string
	getSectionName() string
}

type baseConf struct {
	lock sync.RWMutex

	wsConf    WebsocketConf
	redisConf RedisConf
	dbConf    DbConf
	mqConf    MqConf
}

var BaseConf = &baseConf{}

var (
	IniFileNotFoundErr    = errors.New("配置文件未找到")
	IniSectionNotFoundErr = errors.New("配置项未找到")
)

func init() {
	BaseConf.Load()
}

func (ConfHandle *baseConf) Load() {
	BaseConf.lock.Lock()
	defer BaseConf.lock.Unlock()
	var (
		wsConf = &WebsocketConf{}
		rdConf = &RedisConf{}
		dbConf = &DbConf{}
		mqConf = &MqConf{}
	)
	if cfg, err := match(wsConf); err == nil {
		err := cfg.MapTo(wsConf)
		if err != nil {
			recordError(IniSectionNotFoundErr, wsConf)
		} else {
			BaseConf.wsConf = *wsConf
		}
	}
	if cfg, err := match(rdConf); err == nil {
		err := cfg.MapTo(rdConf)
		if err != nil {
			recordError(IniSectionNotFoundErr, rdConf)
		} else {
			BaseConf.redisConf = *rdConf
		}
	}
	if cfg, err := match(dbConf); err == nil {
		err := cfg.MapTo(dbConf)
		if err != nil {
			recordError(IniSectionNotFoundErr, dbConf)
		} else {
			BaseConf.dbConf = *dbConf
		}
	}
	if cfg, err := match(mqConf); err == nil {
		err := cfg.MapTo(mqConf)
		if err != nil {
			recordError(IniSectionNotFoundErr, mqConf)
		} else {
			BaseConf.mqConf = *mqConf
		}
	}
}

func (ConfHandle *baseConf) GetWsConf() WebsocketConf {
	ConfHandle.lock.RLock()
	defer ConfHandle.lock.RUnlock()
	return ConfHandle.wsConf
}
func (ConfHandle *baseConf) GetRedisConf() RedisConf {
	ConfHandle.lock.RLock()
	defer ConfHandle.lock.RUnlock()
	return ConfHandle.redisConf
}

func (ConfHandle *baseConf) GetDbConf() DbConf {
	ConfHandle.lock.RLock()
	defer ConfHandle.lock.RUnlock()
	return ConfHandle.dbConf
}

func (ConfHandle *baseConf) GetMqConf() MqConf {
	ConfHandle.lock.RLock()
	defer ConfHandle.lock.RUnlock()
	return ConfHandle.mqConf
}

func match(confMap ConfMap) (*ini.Section, error) {
	iniPath := confMap.getPath()
	iniSection := confMap.getSectionName()
	cfg, err := ini.Load(iniPath)
	if err != nil {
		recordError(IniFileNotFoundErr, confMap)
		return nil, IniFileNotFoundErr
	}
	return cfg.Section(iniSection), nil
}

func recordError(error error, confMap ConfMap) {
	errString := fmt.Sprintf("%s:%s", error, confMap.getSectionName())
	logger.Runtime.Error(errString)
}
