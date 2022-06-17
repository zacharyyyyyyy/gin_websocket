package config

type MqConf struct {
	User string
	Pwd  string
	Host string
	Port string
}

func (MqConf MqConf) getPath() string {
	return "conf/config.ini"
}

func (MqConf MqConf) getSectionName() string {
	return "Mq"
}

func (MqConf MqConf) register() {
	if cfg, err := match(&MqConf); err == nil {
		err := cfg.MapTo(&MqConf)
		if err != nil {
			recordError(IniSectionNotFoundErr, &MqConf)
		} else {
			BaseConf.mqConf = MqConf
		}
	}
}
