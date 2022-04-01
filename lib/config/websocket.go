package config

type WebsocketConf struct {
	PingLastTimeSec   int
	ChatLastTimeSec   int
	CleanLimitTimeSec int
	MaxConnection     int
}

func (WsConf WebsocketConf) getPath() string {
	return "conf/config.ini"
}

func (WsConf WebsocketConf) getSectionName() string {
	return "Websocket"
}
