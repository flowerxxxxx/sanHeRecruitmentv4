package config

import "time"

const month = 60 * 60 * 24 * 30 //一个月30天

//项目初始化时的管理员信息-----------------------------------

const AdminUsername = "sanheRec@admin"
const AdminPassword = "sanHeRecAdmin"

//------------------------------------------https运行证书配置

// TLSConfig tls配置
var TLSConfig = &TLSConf{
	Addr:     ":9090", //开启端口
	CertFile: "./ssl/server.pem",
	KeyFile:  "./ssl/server.key",
}

//------------------------------------------mysql配置

// MysqlConfig mysql配置
var MysqlConfig = &MysqlConf{
	Dsn:          "root:020804@(127.0.0.1:3306)/programcom?charset=utf8mb4&parseTime=True&loc=Local",
	Host:         "127.0.0.1",
	Port:         "3306",
	User:         "root",
	Password:     "020804",
	DataBaseName: "programcom",
}

const MysqlConnMaxLivingTime = 300 * time.Second //根据服务器数据库的存活时间配置

//------------------------------------------redis配置

// RedisConfig redis配置
var RedisConfig = &RedisConf{
	Addr:     "127.0.0.1:6379",
	Password: "020804",
	DB:       0,
}

//------------------------------------------消费队列配置

// NsqConfig Nsq配置
var NsqConfig = &NsqConf{
	ProducerAddr:    "127.0.0.1:4150",
	ProducerTopic:   "websocket",
	ConsumerAddr:    "127.0.0.1:4150",
	ConsumerTopic:   "websocket",
	ConsumerChannel: "websocketChannel",
}

//------------------------------------------数据存储配置

// PicSaverPath 图片存储地址
// const PicSaverPath = "./uploadPic"
const PicSaverPath = "D:\\uploadPicSaver"

// MsgExpiredTime 会话消息存储时间
const MsgExpiredTime = month //30天（即时生效，只在生成后的时间生效，之前不生效。）

// BackUpConfig 备份存储地址
var BackUpConfig = &BackupConf{
	SavePath: "./test",
}

// BackerExpireTime 备份存储时间（单位：月）
const BackerExpireTime = 2

//------------------------------------------系统运行配置

// GoroutineNum 主监听线程数
const GoroutineNum = 10

// FirstUnreadMsgNum 对方未回消息用户首次能发的最大消息数-1
const FirstUnreadMsgNum = 10

// ErrorLogAddr 错误及系统日志地址
const ErrorLogAddr = "./logs/systemLogOut.txt"

//-------------------------------------------限流器配置

// ConcurrentPeak 并发峰值
const ConcurrentPeak = 2000

// CurrentLimiterQuantum 每秒添加的令牌数
const CurrentLimiterQuantum = 1500

//------------------------------------------pprof监视配置

//url : https://yanmingyu.free.svipss.top/producer/sanHeRec_pprof/
//router : /producer/sanHeRec_pprof/

const ProducerUsername = "yanmingyu55@gmail.com"
const ProducerPassword = "sanHeRecProducer"
