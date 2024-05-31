package router

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"web/auth"
	"web/config"
	"web/controllers"

	"github.com/gin-gonic/gin"
)

var mySigningKey = []byte("LicheeGameServer")
var jwtInstance = auth.InitJwt(mySigningKey)

const (
	code2sessionURL = "https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code"
)

func Router() *gin.Engine {
	r := gin.Default()

	r.GET("/hello", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "hello world")
	})

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	user := r.Group("user")
	{
		userController := controllers.UserController{}
		user.GET("/info", userController.GetUserInfo)
		user.POST("/list", userController.GetList)
		user.PUT("/add", userController.AddUser)
		user.DELETE("/delete", userController.DeleteUser)
	}

	r.GET("/code2Session", code2Session)
	return r
}

/*
http://localhost:8375/code2Session?appid=wxb00370e58ccf0603&code=1111
GET https://api.weixin.qq.com/sns/jscode2session?appid=APPID&secret=SECRET&js_code=JSCODE&grant_type=authorization_code
{
	"openid":"xxxxxx",
	"session_key":"xxxxx",
	"unionid":"xxxxx",
	"errcode":0,
	"errmsg":"xxxxx"
}
*/

func code2Session(c *gin.Context) {
	log.Println("code2Session")

	appid := c.Query("appid")
	code := c.Query("code")

	if appid != config.Conf.Appid {
		return
	}

	token := jwtInstance.GenerateJWT(appid, 2)
	log.Println(token)
	claims := jwtInstance.ParseJWT(token)
	log.Println(claims)

	url := fmt.Sprintf(code2sessionURL, config.Conf.Appid, config.Conf.Secret, code)
	resp, err := http.Get(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch session key and openId"})
		return
	}

	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
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
