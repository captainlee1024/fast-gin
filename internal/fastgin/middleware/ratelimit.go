package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
)

// RateLimitMiddleware 限流中间件
// fillInterval 填充速率 cap 总容量
func RateLimitMiddleware(fillInterval time.Duration, cap int64) func(c *gin.Context) {
	bucket := ratelimit.NewBucket(fillInterval, cap)
	return func(c *gin.Context) {
		// 如果取不到令牌就返回响应（也可以进行等待，一直等待，或者等待多少时间返回响应）

		//if bucket.Take(1) > 0 { // 返回的是取到令牌所需要等待的时间是多少
		//}

		// 每次请求来了就取一个，如果==0说明没有了，要等，就返回rate limit...，执行后面的函数
		// 如果不大==0就放行
		if bucket.TakeAvailable(1) == 0 { // 返回移除的令牌数或者0
			c.String(http.StatusOK, "rate limit...")
			c.Abort()
			return
		}
		// 取到令牌就放行
		c.Next()
	}
}
