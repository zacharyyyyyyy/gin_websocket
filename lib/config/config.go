package config

import (
	"errors"
	"fmt"
	"sync"

	"gin_websocket/lib/logger"

	"github.com/go-ini/ini"
)

type confHandle interface {
	getPath() string
	getSectionName() string
	register()
}

type baseConf struct {
	lock sync.RWMutex

	wsConf    WebsocketConf
	redisConf RedisConf
	dbConf    DbConf
	mqConf    MqConf
	kafkaConf KafkaConf
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
		kafka  = &KafkaConf{}
	)
	register(wsConf, rdConf, dbConf, mqConf, kafka)
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

func (ConfHandle *baseConf) GetKafkaConf() KafkaConf {
	ConfHandle.lock.RLock()
	defer ConfHandle.lock.RUnlock()
	return ConfHandle.kafkaConf
}

func match(confMap confHandle) (*ini.Section, error) {
	iniPath := confMap.getPath()
	iniSection := confMap.getSectionName()
	cfg, err := ini.Load(iniPath)
	if err != nil {
		recordError(IniFileNotFoundErr, confMap)
		return nil, IniFileNotFoundErr
	}
	return cfg.Section(iniSection), nil
}

func register(confMapArray ...confHandle) {
	for _, confMap := range confMapArray {
		confMap.register()
	}
}

func recordError(error error, confMap confHandle) {
	errString := fmt.Sprintf("%s:%s", error, confMap.getSectionName())
	logger.Runtime.Error(errString)
}
