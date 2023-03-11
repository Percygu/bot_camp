package config

import (
	rlog "github.com/lestrrat-go/file-rotatelogs"
	"github.com/spf13/viper"
	"os"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

var (
	config     GlobalConfig // 全局配置文件
	configFile string
	once       sync.Once
)

// GlobalConfig 业务配置结构体
type GlobalConfig struct {
	LogConfig    *LogConfig          `yaml:"log_conf" mapstructure:"log_conf"`         // 本地日志配置
	BotSvrConfig map[string]*BotConf `yaml:"bot_svr_conf" mapstructure:"bot_svr_conf"` // 写日志服务配置
}

type BotConf struct {
	Port            int      `yaml:"port" mapstructure:"port"`
	AppID           string   `yaml:"app_id" mapstructure:"app_id"`
	AppSecret       string   `yaml:"app_secret" mapstructure:"app_secret"`
	WhiteList       []string `yaml:"white_list" mapstructure:"white_list"`
	ChatIDList      []string `yaml:"chat_id_list" mapstructure:"chat_id_list"`
	BotEncryptedKey string   `yaml:"bot_encrypted_key" mapstructure:"bot_encrypted_key"`
	BotVerifyToken  string   `yaml:"bot_verify_token" mapstructure:"bot_verify_token"`
	EventUrl        string   `yaml:"event_url" mapstructure:"event_url"`
	GroupOwner      string   `yaml:"1v1_group_owner" mapstructure:"1v1_group_owner"`
}

// GetGlobalConf 获取全局配置文件
func GetGlobalConf() *GlobalConfig {
	once.Do(readConf)
	return &config
}

func readConf() {
	viper.SetConfigName("bot")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./conf")
	err := viper.ReadInConfig()
	if err != nil {
		panic("read config file err:" + err.Error())
	}
	err = viper.Unmarshal(&config)
	if err != nil {
		panic("config file unmarshal err:" + err.Error())
	}
}

// InitLog 初始化日志
func InitConfig() {
	globalConf := GetGlobalConf()
	// 设置日志级别
	level, err := log.ParseLevel(globalConf.LogConfig.Level)
	if err != nil {
		panic("log level parse err:" + err.Error())
	}
	// 设置日志格式为json格式
	log.SetFormatter(&logFormatter{
		log.TextFormatter{
			DisableColors:   true,
			TimestampFormat: "2006-01-02 15:04:05",
			FullTimestamp:   true,
		}})
	log.SetReportCaller(true) // 打印文件位置，行号
	log.SetLevel(level)
	switch globalConf.LogConfig.LogPattern {
	case "stdout":
		log.SetOutput(os.Stdout)
	case "stderr":
		log.SetOutput(os.Stderr)
	case "file":
		logger, err := rlog.New(
			globalConf.LogConfig.LogPath+".%Y%m%d",
			//rlog.WithLinkName(globalConf.LogConf.LogPath),
			rlog.WithRotationCount(globalConf.LogConfig.SaveDays),
			//rlog.WithMaxAge(time.Minute*3),
			rlog.WithRotationTime(time.Hour*24),
		)
		if err != nil {
			panic("log conf err: " + err.Error())
		}
		log.SetOutput(logger)
	default:
		panic("log conf err, check log_pattern in logsvr.yaml")
	}
}
