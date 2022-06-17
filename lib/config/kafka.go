package config

type KafkaConf struct {
	User string
	Pwd  string
	Host string
	Port string
}

func (KafkaConf KafkaConf) getPath() string {
	return "conf/config.ini"
}

func (KafkaConf KafkaConf) getSectionName() string {
	return "Kafka"
}

func (KafkaConf KafkaConf) register() {
	if cfg, err := match(&KafkaConf); err == nil {
		err := cfg.MapTo(&KafkaConf)
		if err != nil {
			recordError(IniSectionNotFoundErr, &KafkaConf)
		} else {
			BaseConf.kafkaConf = KafkaConf
		}
	}
}
