package redis

const (
	Prefix                 = "fast-gin:"    // 项目前缀
	KeyFastGinJWTSetPrefix = "fastgin:jwt:" // set 用户登录时的 accessToken
)

// 获取 redis 存储的 key
func GetRedisKey(key string) string {
	return Prefix + key
}
