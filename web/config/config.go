package config

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type LogConfig struct {
	Level      string `json:"level"`      // 日志等级
	Filename   string `json:"filename"`   // 基准日志文件名
	MaxSize    int    `json:"maxsize"`    // 单个日志文件最大内容，单位：MB
	MaxAge     int    `json:"maxage"`     // 日志文件保存时间，单位：天
	MaxBackups int    `json:"maxbackups"` // 最多保存几个日志文件
}

type WechatConf struct {
	Appid  string `json:"appid"`
	Secret string `json:"secret"`
}

type MongoConf struct {
	Dns string `json:"dns"`
	Db  string `json:"db"`
}

type RedisoConf struct {
	Addr     string `json:"addr"`
	Password string `json:"password"`
	Db       int    `json:"db"`
}

type Config struct {
	Port    int          `json:"port"`
	Log     LogConfig    `json:"log"`
	Wechats []WechatConf `json:"wechats"`
	Mongodb MongoConf    `mapstructure:"mongodb"`
	Redis   RedisoConf   `json:"redis"`
}

var (
	Conf      *Config
	ViperConf *viper.Viper
)

func init() {
	ViperConf = viper.New()
	ViperConf.SetConfigType("json")
	ViperConf.SetConfigFile("./conf/config.json")
	err := ViperConf.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	Conf = new(Config)
	if err := ViperConf.Unmarshal(Conf); err != nil {
		panic(fmt.Errorf("unmarshal conf failed, err:%s \n", err))
	}

	ViperConf.WatchConfig()
	ViperConf.OnConfigChange(func(in fsnotify.Event) {
		if err := ViperConf.Unmarshal(Conf); err != nil {
			panic(fmt.Errorf("unmarshal conf failed, err:%s \n", err))
		}
	})
}

func GetWechatInfo(appid string) *WechatConf {
	for _, v := range Conf.Wechats {
		if v.Appid == appid {
			return &v
		}
	}
	return nil
}
