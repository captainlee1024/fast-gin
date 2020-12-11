package middleware

import (
	"errors"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

// SessionAuthMiddleware session 认证中间件
func SessionAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		if name, ok := session.Get("user").(string); !ok || name == "" {
			ResponseError(c, CodeInternalError, errors.New("user not login"))
			c.Abort()
			return
		}
		c.Next()
	}
}
