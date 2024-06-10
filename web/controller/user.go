package controller

import (
	"strconv"
	"web/config"
	"web/model"

	"github.com/gin-gonic/gin"
)

type UserController struct {
}

// /:appid/:code
func (u UserController) Login(ctx *gin.Context) {
	//获取appid,code
	appid := ctx.Param("appid")
	code := ctx.Param("code")
	if code == "" || appid == "" || config.GetWechatInfo(appid) == nil {
		RetErr(ctx, 400, "code appid error")
		return
	}

	//根据code获取openid
	openid := code

	if openid == "" {
		RetErr(ctx, 400, "openid error")
		return
	}

	//根据openid获取用户信息
	user, err := model.GetUserByOpenid(appid, openid)
	if err == nil && user.Openid == openid && !user.ID.IsZero() {
		//登录成功
		RetSuc(ctx, 0, "success", user, 1)
		return
	}

	//新用户则存储用户信息
	user, err = model.AddUser(appid, openid)
	if err == nil && !user.ID.IsZero() {
		//注册成功
		RetSuc(ctx, 0, "success", user, 1)
		return
	}

	RetErr(ctx, 400, "login error")
}

func (u UserController) GetRankUser(ctx *gin.Context) {
	//获取appid,code
	appid := ctx.Param("appid")
	code := ctx.Param("code")
	if code == "" || appid == "" || config.GetWechatInfo(appid) == nil {
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

func (u UserController) UpdateUser(ctx *gin.Context) {
	//带默认值
	uid := ctx.PostForm("uid")
	appid := ctx.PostForm("appid")
	if uid == "" || appid == "" || config.GetWechatInfo(appid) == nil {
		RetErr(ctx, 400, "uid appi error")
		return
	}

	name := ctx.DefaultPostForm("name", "")
	url := ctx.DefaultPostForm("url", "")
	province := ctx.DefaultPostForm("province", "")
	scoreStr := ctx.DefaultPostForm("score", "-1")
	score, _ := strconv.Atoi(scoreStr)
	ret := model.UpdateUser(appid, uid, name, url, province, score)
	if ret {
		RetSuc(ctx, 0, "success", "update user", 1)
	} else {
		RetErr(ctx, 400, "not update")
	}
}
