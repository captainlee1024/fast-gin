package router

import (
	"errors"

	"github.com/captainlee1024/fast-gin/middleware"
	"github.com/captainlee1024/fast-gin/settings"
	"github.com/gin-gonic/gin"
)

// SetUp 初始化路由
func SetUp() *gin.Engine {
	// 当系统设置为 relase 的时候，为发布模式，其他配置都将设置成 debug 模式
	if settings.ConfBase.Mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(middleware.RequestLog(),
		middleware.GinRecovery(true),
		middleware.TranslationMiddleware())

	r.GET("/ping", func(c *gin.Context) {
		middleware.ResponseSuccess(c, "pong")
	})

	r.GET("/ping/error", func(c *gin.Context) {
		middleware.ResponseError(c, 2000, errors.New("test responserror"))
	})

	return r
}
