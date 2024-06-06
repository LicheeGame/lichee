package controller

import (
	"fmt"
	"strconv"
	"web/logger"
	"web/model"

	"github.com/gin-gonic/gin"
)

type UserController struct {
}

func (this UserController) GetUserInfo(ctx *gin.Context) {
	logger.Logger.Info("GetUserInfo")

	idStr := ctx.Param("openid")
	openid, _ := strconv.Atoi(idStr)

	//查询数据库
	user, _ := model.GetUserTest(openid)

	RetSuc(ctx, 0, "success", fmt.Sprintf("user info:%d %d", user.Openid, user.Score), 1)
}

func (this UserController) GetList(ctx *gin.Context) {
	RetErr(ctx, 400, "user list")
}

func (this UserController) AddUser(ctx *gin.Context) {
	RetSuc(ctx, 0, "success", "add user", 1)
}

func (this UserController) DeleteUser(ctx *gin.Context) {
	RetSuc(ctx, 0, "success", "delete user", 1)
}
