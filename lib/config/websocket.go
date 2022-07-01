package config

type WebsocketConf struct {
	//ping包过期时间
	PingLastTimeSec uint
	//最后一次信息过期时间 用于超时清理连接
	ChatLastTimeSec uint
	//连接清理定时
	CleanLimitTimeSec uint
	//最大连接数
	MaxConnection uint
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
