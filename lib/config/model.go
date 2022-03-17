package config

type DbConf struct {
	HostMaster     string
	PortMaster     string
	UsernameMaster string
	PasswordMaster string
	DatabaseMaster string
	HostSlave      string
	PortSlave      string
	UsernameSlave  string
	PasswordSlave  string
	DatabaseSlave  string

	//用于设置连接池中空闲连接的最大数量。
	MaxIdleConns int
	//设置打开数据库连接的最大数量。
	MaxOpenConns int
	//设置了连接可复用的最大时间。
	ConnMaxLifetimeMinus int
}

func (DbConf DbConf) getPath() string {
	return "conf/config.ini"
}

func (DbConf DbConf) getSectionName() string {
	return "Db"
}
