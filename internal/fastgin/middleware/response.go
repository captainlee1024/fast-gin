package middleware

import (
	"encoding/json"
	"fmt"
	"github.com/captainlee1024/fast-gin/internal/pkg/public"
	"net/http"
	"strings"

	"github.com/captainlee1024/fast-gin/internal/fastgin/settings"

	"github.com/gin-gonic/gin"
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
	traceContext := public.GetGinTraceContext(c)
	traceID := ""
	if traceContext != nil {
		traceID = traceContext.TraceID
	}

	resp := &Response{
		ErrorCode: CodeSuccess,
		ErrorMsg:  "",
		Data:      data,
		TraceID:   traceID,
	}

	c.JSON(http.StatusOK, resp)
	response, _ := json.Marshal(resp)
	c.Set(public.CtxResponseKey, string(response))

}

// ResponseError 错误时返回
func ResponseError(c *gin.Context, code ResponseCode, err error) {
	traceContext := public.GetGinTraceContext(c)
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
	c.Set(public.CtxResponseKey, string(response))
	c.AbortWithError(200, err)
}

// 定义状态码
// const (
// 	CodeInvalidParam    = iota // 请求参数有误
// 	CodeUserExist              // 用户已存在
// 	CodeUserNotExist           // 用户不存在
// 	CodeInvalidPassword        // 密码错误
// 	CodeServerBusy             // 服务器繁忙，例如数据库连接错误的时候，不需要吧具体的信息返回给前端用户
// CodeInvalidParam:    "请求参数有误",
// CodeUserExist:       "用户已存在",
// CodeUserNotExist:    "用户不存在",
// CodeInvalidPassword: "用户名或密码错误",
// CodeServerBusy:      "服务繁忙",
// )

// ResponseCode 响应状态码
type ResponseCode int64

// 状态码 1000 以下为通用码，1000 以上为用户自定义码
const (
	CodeSuccess    ResponseCode = iota // success
	CodeUndefError                     // 未知的错误
	CodeValidError
	CodeInternalError // 内部错误

	CodeInvalidRequestError ResponseCode = 401
	CodeCustomize           ResponseCode = 1000

	CodeNeedLogin      ResponseCode = 1100 + iota // 用户需要登录
	CodeInvalidToken                              // 无效 token
	CodeLoginElsewhere                            // 用户在别处登录

	GROUPALL_SAVE_FLOWERROR ResponseCode = 2001
)

// 定义状态码对应信息
var codeMsgMap = map[ResponseCode]string{
	CodeSuccess:       "success",
	CodeUndefError:    "未知的错误",
	CodeValidError:    "验证错误",
	CodeInternalError: "内部错误",

	CodeInvalidRequestError: "",
	CodeCustomize:           "",

	CodeNeedLogin:      "请登录",
	CodeInvalidToken:   "无效的token",
	CodeLoginElsewhere: "账号已在其它客户端登录，重新登录",
}

// Msg 获取状态码对应的提示信息
func (c ResponseCode) Msg() string {
	msg, ok := codeMsgMap[c]
	if !ok {
		msg = codeMsgMap[CodeInternalError]
	}
	return msg
}
