package controller

import (
	"web/auth"
	"web/config"
	"web/model"

	"github.com/gin-gonic/gin"
)

var (
	UseTest = true
)

type UserController struct {
}

// /:appid/:code
func (u UserController) Login(ctx *gin.Context) {
	//获取appid,code
	appid := ctx.Query("appid") //ctx.Param("appid")
	code := ctx.Query("code")   //ctx.Param("code")
	if code == "" || appid == "" || config.GetWechatInfo(appid) == nil {
		RetErr(ctx, 400, "code appid error")
		return
	}

	openid := code

	if !config.Conf.TestDev {
		//根据code获取openid
		openid, err := model.Code2Session(appid, code)
		if openid == "" || err != nil {
			RetErr(ctx, 400, "openid error")
			return
		}
	}

	//根据openid获取用户信息
	user, err := model.GetUserByOpenid(appid, openid)
	if err == nil && user.Openid == openid && !user.ID.IsZero() {
		//登录成功
		user.Token = auth.JWT.GenerateJWT(appid, user.ID.Hex(), 2)
		RetSuc(ctx, 0, "success", user, 1)
		return
	}

	//新用户则存储用户信息
	user, err = model.AddUser(appid, openid)
	if err == nil && !user.ID.IsZero() {
		//注册成功
		user.Token = auth.JWT.GenerateJWT(appid, user.ID.Hex(), 2)
		RetSuc(ctx, 0, "success", user, 1)
		return
	}

	RetErr(ctx, 400, "login error")
}

func (u UserController) GetRankUser(ctx *gin.Context) {
	//token获得的uid和appid
	uid := ctx.GetString("uid")
	appid := ctx.GetString("appid")
	if uid == "" || appid == "" || config.GetWechatInfo(appid) == nil {
		RetErr(ctx, 400, "code appid error")
		return
	}

	results, err := model.GetRankUser(appid)
	if len(results) == 0 || err != nil {
		RetErr(ctx, 400, "no rank")
	} else {
		RetSuc(ctx, 0, "success", results, 1)
	}
}

type UserUpdate struct {
	NickName  string `json:"nickName,omitempty"`
	AvatarUrl string `json:"avatarUrl,omitempty"`
	Province  int    `json:"province,omitempty"`
	Score     int    `json:"score,omitempty"`
}

func (u UserController) UpdateUser(ctx *gin.Context) {
	//token获得的uid和appid
	uid := ctx.GetString("uid")
	appid := ctx.GetString("appid")
	if uid == "" || appid == "" || config.GetWechatInfo(appid) == nil {
		RetErr(ctx, 400, "uid appi error")
		return
	}

	var user UserUpdate
	if err := ctx.BindJSON(&user); err != nil {
		RetErr(ctx, 400, "json body error")
		return
	}
	ret := model.UpdateUser(appid, uid, user.NickName, user.AvatarUrl, user.Province, user.Score)

	if ret {
		RetSuc(ctx, 0, "success", "update user", 1)
	} else {
		RetErr(ctx, 400, "not update")
	}
}
