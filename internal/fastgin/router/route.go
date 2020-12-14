package router

import (
	"github.com/captainlee1024/fast-gin/internal/fastgin/controller"
	"github.com/captainlee1024/fast-gin/internal/fastgin/data"
	mylog "github.com/captainlee1024/fast-gin/internal/fastgin/log"
	"github.com/captainlee1024/fast-gin/internal/fastgin/service"
	"github.com/gin-gonic/contrib/sessions"

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

	v1 := r.Group("/fastgin/v1")
	v1.Use(
		middleware.RequestLog(),
		middleware.GinRecovery(true),
		middleware.TranslationMiddleware(),
		middleware.IPAuthMiddleware(),
	)
	{
		//controller.FastGinRegister(v1)

		fastGin := controller.NewFastGinController(service.NewFastGinUsercase(data.NewFastGinRepo()))

		v1.GET("/", fastGin.IndexHandler)
		v1.GET("error", fastGin.ErrorHandler)
		v1.GET("/ping", fastGin.PingHandler)
		v1.GET("/listpage", fastGin.GetFastGinListHandler)
		v1.POST("/add", fastGin.AddFastGinHandler)
		v1.POST("/get", fastGin.GetFastGinHandler)
		v1.POST("/edit", fastGin.EditFastGinHandler)
		v1.POST("/remove", fastGin.RemoveFastGinHandler)
		v1.POST("/jwt/get", fastGin.GetJWTTokenHandler)

		// 使用 JWT 方式认证的接口
		authJWT := v1.Group("", middleware.JWTAuthMiddleware())
		{
			authJWT.POST("/jwt/auth", fastGin.AuthFastGinJWTHandler)
		}
		// TODO: Session 登录认证待完善
		// 使用 Session 方式认证的接口
		store, err := sessions.NewRedisStore(10, "tcp", settings.GetStringConf("base.session.redis_server"), settings.GetStringConf("base.session.redis_password"), []byte("secret"))
		if err != nil {
			mylog.Log.Fatal("sessions.NewRedisStore", mylog.NewTrace(), mylog.DLTagUndefind, map[string]interface{}{"error": err})
		}
		authSession := v1.Group("",
			sessions.Sessions("mysession", store),
			middleware.SessionAuthMiddleware())
		{
			authSession.GET("/session/auth", fastGin.AuthFastGinSessionHandler)
		}

	}

	return r
}
