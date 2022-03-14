package config

type WebsocketConf struct {
	PingLastTimeSec int
	ChatLastTimeSec int
}

func (WsConf *ConfHandle) MapTo() {

}

func (WsConf ConfHandle) getPath() string {
	return "conf/config.ini"
}
