package common

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 是所有接口统一的返回结构：
//
//	{"code":0,"msg":"","data":{...}}
type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// OK 返回成功响应，业务数据放在 data 字段。
func OK(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, Response{
		Code: 0,
		Msg:  "",
		Data: data,
	})
}

// Fail 返回失败响应，携带业务错误码与提示信息，data 为 null。
func Fail(ctx *gin.Context, code int, msg string) {
	ctx.JSON(http.StatusOK, Response{
		Code: code,
		Msg:  msg,
		Data: nil,
	})
}
