// Package main provides ...
package main

import (
	"errors"
	"time"

	"github.com/captainlee1024/fast-gin/log"

	"github.com/captainlee1024/fast-gin/settings"
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

// @title Fast-Gin（这里写标题）
// @version 1.0
// @description Go Web 通用脚手架
// @termsOfService http://swagger.io/terms/

// @contact.name CaptainLee1024（这里换成你的信息）
// @contact.url http://blog.leecoding.club
// @contact.email 644052732@qq.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1
func main() {
	// 1. 加载配置
	// 2. 初始化日志
	// 3. 初始化 MySQL 连接
	// 4. 初始化 Redis 连接
	if err := settings.Init("./conf/dev/"); err != nil {
		// log.Fatal(err)
		panic(err)
	}
	// if err := settings.InitModule("./conf/dev/", []string{"base", "redis", "mysql"}); err != nil {
	// 	log.Fatal(err)
	// }

	// 5. 注册路由
	// 6. 启动服务（开启平滑下线）
	//time.Sleep(time.Second * 10)

	// test debug
	log.Log.Debug("/debug", settings.NewTrace(), log.DLTagUndefind,
		map[string]interface{}{
			"message":  "debug 测试替换日志默认Caller",
			"error":    errors.New("text string"),
			"balabala": "xxxx"})

	// todo sth
	log.Log.Info("/test", settings.NewTrace(), log.DLTagUndefind,
		map[string]interface{}{
			"message":  "todo sth",
			"error":    errors.New("text string"),
			"balabala": "xxxx"})

	time.Sleep(time.Second * 10)

}
