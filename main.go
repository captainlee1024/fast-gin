package main

import (
	"context"
	"fmt"
	"github.com/captainlee1024/fast-gin/pkg/snowflake"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/captainlee1024/fast-gin/internal/fastgin/data/mysql"
	"github.com/captainlee1024/fast-gin/internal/fastgin/data/redis"
	mylog "github.com/captainlee1024/fast-gin/internal/fastgin/log"
	"github.com/captainlee1024/fast-gin/internal/fastgin/router"
	"github.com/captainlee1024/fast-gin/internal/fastgin/settings"
)

/* swagger main 函数注释格式（写项目相关介绍信息）
// @title 这里写标题
// @version 1.0
// @description 这里写描述信息
// @termsOfService http://swagger.io/terms/

// @contact.name 这里写联系人信息
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host 这里写接口服务的host
// @BasePath 这里写base path
*/

// @title Fast-Gin
// @version 1.0
// @description Fast-Gin 是一个能够帮助你快速进行开发的 Web 通用脚手架
// @termsOfService http://swagger.io/terms/

// @contact.name CaptainLee1024（这里换成你的信息）
// @contact.url http://blog.leecoding.club
// @contact.email 644052732@qq.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host 127.0.0.1:8080
// @BasePath /fastgin/v1
func main() {
	// 1. 加载配置
	// 2. 初始化日志
	if err := settings.Init("./configs/dev/"); err != nil {
		// log.Fatal(err)
		panic(err)
	}

	trace := mylog.NewTrace()

	// 3. 初始化 MySQL 连接
	if err := mysql.InitDBPool(); err != nil {
		mylog.Log.Error("mysql", trace, mylog.DLTagUndefind, map[string]interface{}{
			"error": err,
		})
	}
	// 释放 mysql 资源，并且刷新缓冲里的日志信息
	defer func() {
		log.Println("------------------------------------------------------------------------")
		log.Printf("[INFO] %s\n", " start destroy resources.")
		mysql.Close()
		mylog.Log.L.Sync()
		log.Printf("[INFO] %s\n", " success destroy resources.")
	}()

	// 4. 初始化 Redis 连接
	defaultConn, err := redis.ConnFactory("default")
	if err != nil {
		mylog.Log.Error("redis", trace, mylog.DLTagUndefind, map[string]interface{}{
			"error": err,
		})
	}
	defer defaultConn.Close()

	// 初始化雪花算法
	if err := snowflake.Init(settings.ConfBase.StartTime, settings.ConfBase.MachineID); err != nil {
		mylog.Log.Error("initSnowflake", trace, mylog.DLTagUndefind, map[string]interface{}{
			"error": err,
		})
		return
	}

	// 5. 注册路由
	r := router.SetUp()

	// 6. 启动服务（开启平滑下线）
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%v", settings.ConfBase.Port),
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			mylog.Log.Fatal("listen", mylog.NewTrace(), mylog.DLTagUndefind, map[string]interface{}{
				"err": err,
			})
		}
	}()

	// 等待中断信号来优雅关闭服务器，为关闭服务器操作设置一个5秒的延时
	quit := make(chan os.Signal, 1)
	// kill 默认会发送 syscall.SIGTERM 信号
	// kill -2 发送 syscall.SIGINT 信号，我们常用的Ctrl+C就是触发系统SIGINT信号
	// kill -9 发送 syscall.SIGKILL 信号，但是不能被捕获，所以不需要添加它
	// signal.Notify把收到的 syscall.SIGINT或syscall.SIGTERM 信号转发给quit
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // 此处不会阻塞
	<-quit

	shoutdownTrace := mylog.NewTrace()
	mylog.Log.Info("Shoutdown", shoutdownTrace, mylog.DLTagUndefind, map[string]interface{}{
		"msg": "Shoutdown Server ...",
	})

	// 创建一个 5 秒超时的 context
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		mylog.Log.Fatal("Shoutdown", shoutdownTrace, mylog.DLTagUndefind, map[string]interface{}{
			"error": err,
		})
	}

	mylog.Log.Info("Server exiting", shoutdownTrace, mylog.DLTagUndefind, map[string]interface{}{
		"msg": "Server exiting",
	})

	/*
		// test mylog debug
		mylog.Log.Debug("/debug", mylog.NewTrace(), mylog.DLTagUndefind,
			map[string]interface{}{
				"message":  "debug 测试替换日志默认Caller",
				"error":    errors.New("text string"),
				"balabala": "xxxx"})

		// test mylog info
		mylog.Log.Info("/test", mylog.NewTrace(), mylog.DLTagUndefind,
			map[string]interface{}{
				"message":  "todo sth",
				"error":    errors.New("text string"),
				"balabala": "xxxx"})

		// test mylog error
		mylog.Log.Error("/error", mylog.NewTrace(), mylog.DLTagUndefind,
			map[string]interface{}{
				"message":  "error 级别日志测试",
				"error":    errors.New("text string"),
				"balabala": "xxxx"})

		// time.Sleep(time.Second * 10)
	*/
}
