package config

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

/*
{
	wechat:{
		"appId": "123456789",
		"secret": "maple123456",
	},
	"host": {
	  "address": "localhost",
	  "port": 5799
	},
	"mongodb":{

	},
	"redis":{

	}
}
*/

type Config struct {
	Port   int    `json:"port"`
	Appid  string `json:"appid"`
	Secret string `json:"app_secret"`

	Log LogConfig
}

type LogConfig struct {
	Level      string `json:"level"`       // 日志等级
	Filename   string `json:"filename"`    // 基准日志文件名
	MaxSize    int    `json:"maxsize"`     // 单个日志文件最大内容，单位：MB
	MaxAge     int    `json:"max_age"`     // 日志文件保存时间，单位：天
	MaxBackups int    `json:"max_backups"` // 最多保存几个日志文件
}

var Conf = new(Config)

func init() {
	viperConf := viper.New()
	viperConf.SetConfigFile("./conf/config.json")
	err := viperConf.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	if err := viperConf.Unmarshal(Conf); err != nil {
		panic(fmt.Errorf("unmarshal conf failed, err:%s \n", err))
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		if err := viper.Unmarshal(Conf); err != nil {
			panic(fmt.Errorf("unmarshal conf failed, err:%s \n", err))
		}
	})

}
