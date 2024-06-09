package controller

import (
	"strconv"
	"web/model"

	"github.com/gin-gonic/gin"
)

type UserController struct {
}

func (this UserController) Login(ctx *gin.Context) {
	openid := ctx.Param("code")
	if openid == "" {
		return
	}

	user, err := model.GetUserByOpenid(openid)
	if err == nil && user.Openid == openid {
		//登录成功
		RetSuc(ctx, 0, "success", user, 1)
		return
	}

	user, err = model.AddUser(openid)
	if err == nil && user.ID.String() != "" {
		//注册成功
		RetSuc(ctx, 0, "success", user, 1)
		return
	}

	RetErr(ctx, 400, "login error")
}

func (this UserController) GetRankUser(ctx *gin.Context) {
	results := model.GetRankUser()
	if len(results) != 0 {
		RetErr(ctx, 400, "no rank")
	} else {
		RetSuc(ctx, 0, "success", results, 1)
	}
}

func (this UserController) UpdateUser(ctx *gin.Context) {
	//带默认值
	uid := ctx.PostForm("uid")
	name := ctx.DefaultPostForm("name", "")
	url := ctx.DefaultPostForm("url", "")
	province := ctx.DefaultPostForm("province", "")
	scoreStr := ctx.DefaultPostForm("score", "0")
	score, _ := strconv.Atoi(scoreStr)
	ret := model.UpdateUser(uid, name, url, province, score)
	if ret {
		RetSuc(ctx, 0, "success", "update user", 1)
	} else {
		RetErr(ctx, 400, "not update")
	}
}
