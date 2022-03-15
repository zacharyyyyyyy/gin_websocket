package config

type RedisConf struct {
	Host     string
	Port     int
	Password string
	Db       int
}

func (RdConf RedisConf) getPath() string {
	return "conf/config.ini"
}

func (RdConf RedisConf) getSectionName() string {
	return "Redis"
}
