package middleware

import (
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"

	mylog "github.com/captainlee1024/fast-gin/internal/fastgin/log"
	"github.com/gin-gonic/gin"
)

// GinRecovery recover掉项目可能出现的panic
func GinRecovery(stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					// zap.L().Error(c.Request.URL.Path,
					// 	zap.Any("error", err),
					// 	zap.String("request", string(httpRequest)),
					// )
					// If the connection is dead, we can't write a status to it.
					c.Error(err.(error)) // nolint: errcheck
					mylog.Log.Error(c.Request.URL.Path, mylog.NewTrace(), "`dltag string`", map[string]interface{}{
						"error":   err,
						"request": string(httpRequest),
					})
					c.Abort()
					return
				}

				if stack {
					// zap.L().Error("[Recovery from panic]",
					// 	zap.Any("error", err),
					// 	zap.String("request", string(httpRequest)),
					// 	zap.String("stack", string(debug.Stack())),
					// )
					mylog.Log.Error("[Recovery from panic]", mylog.NewTrace(), "dltag string", map[string]interface{}{
						"error":   err,
						"request": string(httpRequest),
						"stack":   string(debug.Stack()),
					})
				} else {
					// zap.L().Error("[Recovery from panic]",
					// 	zap.Any("error", err),
					// 	zap.String("request", string(httpRequest)),
					// )
					mylog.Log.Error("[Recovery from panic]", mylog.NewTrace(), "dltag string", map[string]interface{}{
						"error":   err,
						"request": string(httpRequest),
					})
				}
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
