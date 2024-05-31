package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type JsonStruct struct {
	Code  int         `json:"code"`
	Msg   interface{} `json:"msg"`
	Data  interface{} `json:"data"`
	Count int         `json:"count"`
}

type JsonErrStruct struct {
	Code int         `json:"code"`
	Msg  interface{} `json:"msg"`
}

func RetSuc(ctx *gin.Context, code int, msg interface{}, data interface{}, count int) {
	ret := &JsonStruct{Code: code, Msg: msg, Data: data, Count: count}
	ctx.JSON(http.StatusOK, ret)
}

func RetErr(ctx *gin.Context, code int, msg interface{}) {
	ret := &JsonErrStruct{Code: code, Msg: msg}
	ctx.JSON(http.StatusOK, ret)
}
