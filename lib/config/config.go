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

	wsConf WebsocketConf
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
	)
	if cfg, err := match(wsConf); err == nil {
		err := cfg.MapTo(wsConf)
		if err != nil {
			errString := fmt.Sprintf("%s:%s", IniSectionNotFoundErr, wsConf.getSectionName())
			logger.Service.Error(errString)
		} else {
			BaseConf.wsConf = *wsConf
		}
	}
}

func match(confMap ConfMap) (*ini.Section, error) {
	iniPath := confMap.getPath()
	iniSection := confMap.getSectionName()
	cfg, err := ini.Load(iniPath)
	if err != nil {
		errString := fmt.Sprintf("%s:%s", err, confMap.getSectionName())
		logger.Service.Error(errString)
		return nil, IniFileNotFoundErr
	}
	return cfg.Section(iniSection), nil
}

func (ConfHandle *baseConf) GetWsConf() WebsocketConf {
	ConfHandle.lock.RLock()
	defer ConfHandle.lock.RUnlock()
	return ConfHandle.wsConf
}
