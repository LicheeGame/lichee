package router

import (
	"strings"
	"web/auth"
	"web/controller"
	"web/logger"

	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	r := gin.New()
	r.Use(logger.GinLogger(logger.Logger), logger.GinRecovery(logger.Logger, true))

	user := r.Group("minigame/api/user")
	{
		userController := controller.UserController{}
		user.GET("/login/:appid/:code", userController.Login)
		user.POST("/update", JwtAuth(), userController.UpdateUser)
		user.GET("/ranklist", JwtAuth(), userController.GetRankUser)
	}

	//r.GET("/code2Session", code2Session)
	return r
}

func JwtAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		//Header的authorization中获取tokenString
		tokenHeader := c.Request.Header.Get("Authorization")
		kv := strings.Split(tokenHeader, " ")
		if tokenHeader == "" || len(kv) != 2 || kv[0] != "Bearer" {
			controller.RetErr(c, 400, "token error")
			c.Abort()
			return
		}
		tokenString := kv[1]

		// Parse token
		customClaims := auth.JWT.ParseJWT(tokenString)
		if customClaims == nil || customClaims.UID == "" || customClaims.Appid == "" {
			controller.RetErr(c, 400, "token error")
			c.Abort()
			return
		}

		//将uid写入请求参数
		c.Set("uid", customClaims.UID)
		c.Set("appid", customClaims.Appid)
		c.Next()
	}
}

/*
中间件：https://blog.csdn.net/Shoulen/article/details/136141292
參數：https://www.jianshu.com/p/916ce255de83

HTTP上传参数3个部分：Header 、URL、 HTTP Body
Header：键值对集合 Content-Type Accept

GET请求
URL：请求路径， http://localhost:8080/user/add/1/name
	参数获取：ctx.Param("id")
	//http://localhost:8080/user/100
	engine.GET("/user/:id/:name", func(ctx *gin.Context) {
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

  POST请求
  HTTP Body中的参数
    //http://localhost:8080/user/add
	//body中添加name= gender= habits= worls=
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

//通过绑定可以接受json的body数据
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
//绑定时可以传入结构体或者param := make(map[string]interface{})
    engine.POST("/user/add", func(ctx *gin.Context) {
        var u User
        if err := ctx.Bind(&u); err != nil {
            ctx.JSON(http.StatusBadRequest, err.Error())
            return
        }
        fmt.Fprintf(ctx.Writer, "你输入的用户名为：%s,邮箱为：%s\n", u.Name, u.Email)
    })
*/
