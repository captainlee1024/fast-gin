package middleware

import (
	"fmt"

	"github.com/captainlee1024/fast-gin/internal/fastgin/settings"
	"github.com/gin-gonic/gin"
)

// IPAuthMiddleware IP 白名单中间件
func IPAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		isMatched := false
		// fmt.Printf("\n\n\n%v\n\n\n", settings.GetStringSliceConf("base.http.allow_ip"))
		for _, host := range settings.GetStringSliceConf("base.http.allow_ip") {
			if c.ClientIP() == host {
				isMatched = true
			}
		}

		if !isMatched {
			ResponseError(c, CodeInternalError, fmt.Errorf("%v not in iplist", c.ClientIP()))
			c.Abort()
			return
		}
		c.Next()
	}
}
