package controller

import (
	"errors"

	"github.com/captainlee1024/fast-gin/internal/fastgin/middleware"
	"github.com/gin-gonic/gin"
)

// FastGinController (微服务中是 FastGinService)
type FastGinController struct {
}

// FastGinRegister 注册路由处理器
func FastGinRegister(router *gin.RouterGroup) {
	fastGin := FastGinController{}
	router.GET("/", fastGin.Index)
	router.GET("error", fastGin.Error)
}

// Index index 处理器
func (fastGin *FastGinController) Index(c *gin.Context) {
	middleware.ResponseSuccess(c, "index Success test")
	return
}

// Error error 处理器
func (fastGin *FastGinController) Error(c *gin.Context) {
	middleware.ResponseError(c, 2000, errors.New("test middleware ResponseError"))
}

// NewFastGinController （微服务中是 NewFastGinService）
func NewFastGinController() *FastGinController {
	return &FastGinController{}
}

// FastGin fast-gin 测试路由
func FastGin(c *gin.Context) {

}
