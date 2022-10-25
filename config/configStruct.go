package config

// TLSConf tls配置
type TLSConf struct {
	Addr     string `json:"addr"`
	CertFile string `json:"cert_file"`
	KeyFile  string `json:"key_file"`
}

// MysqlConf mysql配置
type MysqlConf struct {
	Dsn          string `json:"dsn"`
	Host         string `json:"host"`
	Port         string `json:"port"`
	User         string `json:"user"`
	Password     string `json:"password"`
	DataBaseName string `json:"data_base_name"`
}

// RedisConf redis配置
type RedisConf struct {
	Addr     string `json:"addr"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}

// NsqConf nsq配置
type NsqConf struct {
	ProducerAddr    string `json:"producer_addr"`
	ProducerTopic   string `json:"producer_topic"`
	ConsumerAddr    string `json:"consumer_addr"`
	ConsumerTopic   string `json:"consumer_topic"`
	ConsumerChannel string `json:"consumer_channel"`
}

type BackupConf struct {
	SavePath string `json:"save_path"`
}
