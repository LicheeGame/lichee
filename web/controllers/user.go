package controllers

import "github.com/gin-gonic/gin"

type UserController struct {
}

func (this UserController) GetUserInfo(ctx *gin.Context) {
	RetSuc(ctx, 0, "success", "user info", 1)
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
