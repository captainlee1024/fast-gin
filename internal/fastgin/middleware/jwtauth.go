package middleware

import (
	"errors"
	"fmt"
	"github.com/captainlee1024/fast-gin/internal/fastgin/data"
	"github.com/captainlee1024/fast-gin/internal/pkg/public"
	"strings"

	"github.com/captainlee1024/fast-gin/pkg/jwt"
	"github.com/gin-gonic/gin"
)

// JWTAuthMiddleware JWT 认证中间件
// 1. 查看请求头里是否带有 Token
// 2. 查看 Token 格式是否正确
// 3. 解析 Token
// 4. 判断与 Redis 中的 Token 是否相等（这样可以保证只能有一个客户端登录，在其他地方登录 Redis Token 刷新，之前的登录就作废了）
// 5. 都满足，进行下一步
func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		// // redisToken = string(redisToken.(string))
		// fmt.Println(redisToken)
		// 客户端携带 Token 有三种方式 1. 放在请求头中 2. 放在请求体中 3. 放在 URI
		// 这里假设 Token 放在 Header 的 Autherization 中，并使用 Bearer 开头
		// Authorization: Bearer xxx.xxx.xxx 或者是前端团队自己定义的其他格式，例如：X-TOKEN xx.xx.xx
		// 具体根据自己的团队和业务逻辑决定
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			ResponseError(c, CodeNeedLogin, errors.New(CodeNeedLogin.Msg()))
			c.Abort()
			return
		}
		// 按空格分割，获取 JWT
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			ResponseError(c, CodeInvalidToken, errors.New(CodeInvalidToken.Msg()))
			c.Abort()
			return
		}
		// parts[1] 是获取到的 tokenString, 我们使用之前定义好的解析 JWT 的函数来解析它
		mc, err := jwt.ParseToken(parts[1])
		if err != nil {
			ResponseError(c, CodeInvalidToken, errors.New(CodeInvalidToken.Msg()))
			c.Abort()
			return
		}

		// 解析有效之后在判断是否与 Redis 中存的 token 相等
		//trace := GetGinTraceContext(c)
		//redisToken, err := red.String(redis.ConfDo(trace, "default", "GET", fmt.Sprint(mc.UserID)))
		redisToken, err := data.GetAToken(fmt.Sprint(mc.UserID), c)

		// redisToken = string(redisToken.(string))
		if err != nil {
			ResponseError(c, CodeNeedLogin, errors.New(CodeNeedLogin.Msg()))
			c.Abort()
			return
		}

		if parts[1] != redisToken {
			ResponseError(c, CodeLoginElsewhere, errors.New(CodeLoginElsewhere.Msg()))
			c.Abort()
			return
		}
		//fmt.Println(authHeader, parts)
		// 将当前请求的 userID 信息保存到请求的上下文 c 中
		c.Set(public.CtxUserIDKey, mc.UserID)
		c.Next() // 后续的处理函数可以通过 c.Get(controller.CtxUserIDKey) 来获取当前请求用户的 userID
	}
}
