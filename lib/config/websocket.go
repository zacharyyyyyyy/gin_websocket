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

func (WsConf WebsocketConf) register() {
	if cfg, err := match(&WsConf); err == nil {
		err := cfg.MapTo(&WsConf)
		if err != nil {
			recordError(IniSectionNotFoundErr, &WsConf)
		} else {
			BaseConf.wsConf = WsConf
		}
	}
}
