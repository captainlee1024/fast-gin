package middleware

import (
	"time"

	mylog "github.com/captainlee1024/fast-gin/log"
	"github.com/gin-gonic/gin"
)

// GinLogger 接收gin框架默认的日志
func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		c.Next()

		cost := time.Since(start)
		// zap.L().Info(path,
		// 	//logger
		// 	zap.Int("status", c.Writer.Status()),
		// 	zap.String("method", c.Request.Method),
		// 	zap.String("path", path),
		// 	zap.String("query", query),
		// 	zap.String("ip", c.ClientIP()),
		// 	zap.String("user-agent", c.Request.UserAgent()),
		// 	zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
		// 	zap.Duration("cost", cost),
		// )
		mylog.Log.Info(path, mylog.NewTrace(), "dltag string", map[string]interface{}{
			"status":     c.Writer.Status(),
			"method":     c.Request.Method,
			"path":       path,
			"query":      query,
			"ip":         c.ClientIP(),
			"user-agent": c.Request.UserAgent(),
			"errors":     c.Errors.ByType(gin.ErrorTypePrivate).String(),
			"cost":       cost,
		})
	}
}
