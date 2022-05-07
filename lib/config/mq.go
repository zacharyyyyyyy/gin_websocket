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
