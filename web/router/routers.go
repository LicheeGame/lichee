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
		// /user/info
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

/*
中间件：https://blog.csdn.net/Shoulen/article/details/136141292
參數：https://www.jianshu.com/p/916ce255de83

HTTP上传参数3个部分：Header 、URL、 HTTP Body
Header：键值对集合 Content-Type Accept

URL：请求路径， http://localhost:8080/user/add/1
	参数获取：ctx.Param("id")
	//http://localhost:8080/user/100
	engine.GET("/user/:id", func(ctx *gin.Context) {
    	id := ctx.Param("id")
	}


HTTP Body：请求体所携带的参数， Content-Type：application/json时 body是一个json串
    获取URL Query中的参数
	  //http://localhost:8080/user/list?name=test&gender=xxxx'&habits=1,2,3,4,5&map["name"]=ji&map["age"]=18
	  engine.GET("/user/list", func(ctx *gin.Context) {
		//获取单个值
		name := ctx.Query("name")
		//带默认值
		gender := ctx.DefaultQuery("gender", "男")
		//数组
		habits := ctx.QueryArray("habits")
		//map
		works := ctx.QueryMap("works")
		fmt.Printf("%s, %s, %s, %s\n", name, gender, habits, works)
  })

  HTTP Body中的参数
    engine.POST("/user/add", func(ctx *gin.Context) {
		//获取单个值
		name := ctx.PostForm("name")
		//带默认值
		gender := ctx.DefaultPostForm("gender", "男")
		//数组
		habits := ctx.PostFormArray("habits")
		//map
		works := ctx.PostFormMap("works")
		fmt.Printf("%s,%s,%s,%s\n", name, gender, habits, works)
  })

//绑定请求参数 BindUri()或者ShouldBindUri()
type User struct {
  Name  string `uri:"name"`
  Email string `uri:"email"`
}
engine.GET("/user/list/:id/:name", func(ctx *gin.Context) {
	var u User
	if err := ctx.BindUri(&u);err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	fmt.Fprintf(ctx.Writer, "你输入的用户名为：%s,邮箱为：%s\n", u.Name, u.Email)
})

//绑定URL Query参数
    engine.GET("/user/list", func(ctx *gin.Context) {
        var u User
        if err := ctx.BindQuery(&u);err != nil {
            ctx.JSON(http.StatusBadRequest, err)
            return
        }
        fmt.Fprintf(ctx.Writer, "你输入的用户名为：%s,邮箱为：%s\n", u.Name, u.Email)
    })
//绑定HTTP Body参数,POST请求时才会进行绑定  Bind BindJSON()
    engine.POST("/user/add", func(ctx *gin.Context) {
        var u User
        if err := ctx.Bind(&u); err != nil {
            ctx.JSON(http.StatusBadRequest, err.Error())
            return
        }
        fmt.Fprintf(ctx.Writer, "你输入的用户名为：%s,邮箱为：%s\n", u.Name, u.Email)
    })
*/
