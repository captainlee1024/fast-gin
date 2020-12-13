package router

import (
	"github.com/captainlee1024/fast-gin/internal/fastgin/controller"
	"github.com/captainlee1024/fast-gin/pkg/jwt"

	"github.com/captainlee1024/fast-gin/internal/fastgin/middleware"
	"github.com/captainlee1024/fast-gin/internal/fastgin/settings"
	"github.com/gin-gonic/gin"
)

// SetUp 初始化路由
func SetUp() *gin.Engine {
	// 当系统设置为 relase 的时候，为发布模式，其他配置都将设置成 debug 模式
	if settings.ConfBase.Mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	v1 := r.Group("/fastgin")
	v1.Use(
		middleware.RequestLog(),
		middleware.GinRecovery(true),
		middleware.TranslationMiddleware(),
	)
	{
		controller.FastGinRegister(v1)
	}

	// 非登录接口
	apiNormalGroup := r.Group("/api")
	apiNormalGroup.Use(
		middleware.RequestLog(),
		middleware.GinRecovery(true),
		middleware.TranslationMiddleware(),
		middleware.IPAuthMiddleware(),
	)
	{
		apiNormalGroup.GET("/jwt", func(c *gin.Context) {
			// 注意：保存到 Redis 的时候，Key 是转换成 string 类型之后的 ID 这里是 11
			jwt, _ := jwt.GenToken(int64(11), "Jojo")
			middleware.ResponseSuccess(c, jwt)
			return
		})
	}

	// 登录接口
	apiAuthGroup := r.Group("/api")
	apiAuthGroup.Use(
		middleware.RequestLog(),
		middleware.GinRecovery(true),
		middleware.TranslationMiddleware(),
		middleware.JWTAuthMiddleware(),
	)
	{
		apiAuthGroup.GET("/jwtauth", func(c *gin.Context) {
			middleware.ResponseSuccess(c, "JWTAuth Success!")
			return
		})

	}

	return r
}
