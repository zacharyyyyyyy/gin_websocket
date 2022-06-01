package config

type WebsocketConf struct {
	PingLastTimeSec   uint
	ChatLastTimeSec   uint
	CleanLimitTimeSec uint
	MaxConnection     uint
}

func (WsConf WebsocketConf) getPath() string {
	return "conf/config.ini"
}

func (WsConf WebsocketConf) getSectionName() string {
	return "Websocket"
}
