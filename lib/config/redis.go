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

func (RdConf RedisConf) register() {
	if cfg, err := match(&RdConf); err == nil {
		err := cfg.MapTo(&RdConf)
		if err != nil {
			recordError(IniSectionNotFoundErr, &RdConf)
		} else {
			BaseConf.redisConf = RdConf
		}
	}
}
