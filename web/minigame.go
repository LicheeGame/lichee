package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type Config struct {
	Port   int
	Appid  string
	Secret string
}

var Conf = new(Config)

// //GET https://api.weixin.qq.com/sns/jscode2session?appid=APPID&secret=SECRET&js_code=JSCODE&grant_type=authorization_code
const (
	code2sessionURL = "https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code"
)

func main() {
	viperConf := viper.New()
	viperConf.SetConfigFile("./conf/config.yaml")
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

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/code2Session", code2Session)

	if err := r.Run(fmt.Sprintf(":%d", Conf.Port)); err != nil {
		panic(err)
	}

	//

	//r.RunTLS(":8080", "./testdata/server.pem", "./testdata/server.key")
}

// /http://localhost:8375/code2Session?appid=wxb00370e58ccf0603&code=1111
func code2Session(c *gin.Context) {
	appid := c.Query("appid")
	code := c.Query("code")

	if appid != Conf.Appid {
		return
	}

	url := fmt.Sprintf(code2sessionURL, Conf.Appid, Conf.Secret, code)
	resp, err := http.Get(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch session key and openId"})
		return
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var result map[string]interface{}
	json.Unmarshal(body, &result)
	if result["errcode"] != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": result["errmsg"].(string)})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"openId":  result["openid"],
	})
}
