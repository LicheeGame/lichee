package controllers

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type UserController struct {
}

func (this UserController) GetUserInfo(ctx *gin.Context) {
	id := ctx.Param("id")
	name := ctx.Param("name")
	RetSuc(ctx, 0, "success", fmt.Sprintf("user info:%d %s", id, name), 1)
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
