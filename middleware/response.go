package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/captainlee1024/fast-gin/settings"

	"github.com/gin-gonic/gin"
)

//context需要设置的东西
const (
	ContextResponse = "response"
)

// ResponseCode 响应状态码
type ResponseCode int

// 状态码 1000 以下为通用码，1000 以上为用户自定义码
const (
	SuccessCode ResponseCode = iota
	UndefErrorCode
	ValidErrorCode
	InternalErrorCode

	InvalidRequestErrorCode ResponseCode = 401
	CustomizeCode           ResponseCode = 1000

	GROUPALL_SAVE_FLOWERROR ResponseCode = 2001
)

// Response 响应结构体
type Response struct {
	ErrorCode ResponseCode `json:"errno"`
	ErrorMsg  string       `json:"errmsg"`
	Data      interface{}  `json:"data"`
	TraceID   interface{}  `json:"trace_id"`
	Stack     interface{}  `json:"stack"`
}

// ResponseSuccess 成功时返回
func ResponseSuccess(c *gin.Context, data interface{}) {
	traceContext := GetGinTraceContext(c)
	traceID := ""
	if traceContext != nil {
		traceID = traceContext.TraceID
	}

	resp := &Response{
		ErrorCode: SuccessCode,
		ErrorMsg:  "",
		Data:      data,
		TraceID:   traceID,
	}

	c.JSON(http.StatusOK, resp)
	response, _ := json.Marshal(resp)
	c.Set(ContextResponse, string(response))

}

// ResponseError 错误时返回
func ResponseError(c *gin.Context, code ResponseCode, err error) {
	traceContext := GetGinTraceContext(c)
	traceID := ""
	if traceContext != nil {
		traceID = traceContext.TraceID
	}

	stack := ""
	if c.Query("is_debug") == "1" || settings.ConfBase.Mode == "dev" {
		stack = strings.Replace(fmt.Sprintf("%+v", err), err.Error()+"\n", "", -1)
	}

	resp := &Response{
		ErrorCode: code,
		ErrorMsg:  err.Error(),
		Data:      "",
		TraceID:   traceID,
		Stack:     stack,
	}

	c.JSON(200, resp)
	response, _ := json.Marshal(resp)
	c.Set(ContextResponse, string(response))
	c.AbortWithError(200, err)
}
