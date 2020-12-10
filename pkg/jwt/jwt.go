package jwt

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// TokenExpireDuration Token过期时间
const TokenExpireDuration = time.Hour * 2

// 定义加密的盐
var mySecret = []byte("我是盐")

// MyClaims 自定义声明结构体并内嵌 jwt.StandarCMyClaims
// jwt 包自带的 jwt.StandardClaims 只包含了官方字段
// 我们这里需要额外记录一个 username 字段，所以要自定义结构体
// 如果想要保存更多的信息，都可以添加到这个结构体中
type MyClaims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	jwt.StandardClaims
}

// GenToken 生成一个 Token
func GenToken(userID int64, username string) (string, error) {
	// 创建一个我们声明的数据
	c := MyClaims{
		UserID:   userID,
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(TokenExpireDuration).Unix(), // 过期时间
			// 从配置文件设置过期时间
			//ExpiresAt: time.Now().Add(time.Duration(viper.GetInt("auth.jwt_expire")) * time.Hour).Unix()
			Issuer: "bluebell", // 签发人
		},
	}
	// 使用指定的签名方法创建对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)

	// 使用指定的 secret 签名并获得完整的编码后的字符串 token
	return token.SignedString(mySecret)
}

// ParseToken 解析 Token
func ParseToken(tokenString string) (*MyClaims, error) {
	// 解析 token
	var mc = new(MyClaims) // 注意，返回值这里声明一个变量是不会帮你申请内存的，所以这里要先new()一下
	token, err := jwt.ParseWithClaims(tokenString, mc, func(token *jwt.Token) (i interface{}, err error) {
		return mySecret, nil
	})
	if err != nil {
		return nil, err
	}
	if token.Valid { // 校验 token
		return mc, nil
	}
	return nil, errors.New("invalid token")
}

// 生成 Access Token 和 Refresh Token
func _GenToken(userID int64, username string) (aToken, rToken string, err error) {
	// 创建一个我们自己的声明
	c := MyClaims{
		userID,
		username,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(TokenExpireDuration).Unix(),
		},
	}

	// 加密并获得完整的编码后的字符串 token
	aToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, c).SignedString(mySecret)

	// refresh token 不需要任何自定义数据
	rToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Second * 30).Unix(), // 过期时间
		Issuer:    "bluebell",                              // 签发人
	}).SignedString(mySecret)
	return
}

// RefreshToken 刷新 Access Token 的接口
// return (newAToken, newRToken string, err error) 如果我们邀请能够续token，
// 那就每次请求都返回一个新的Refreshtoken，这样就会和循环自动续上refreshtoken的有效期
// 如果我们希望refresh到期强制重新登录，不希望能够续token的话，可以只返回新的Accress Token，refreshtoken返回旧的
func RefreshToken(aToken, rToken string) (newAToken, newRToken string, err error) {
	// refresh token 无效，直接返回
	if _, err = jwt.Parse(rToken, func(token *jwt.Token) (i interface{}, err error) {
		return mySecret, nil
	}); err != nil {
		return
	}

	// 从旧 Access Token 中解析出 claims 数据
	var claims MyClaims
	_, err = jwt.ParseWithClaims(aToken, &claims, func(token *jwt.Token) (i interface{}, err error) {
		return mySecret, nil
	})
	v, _ := err.(*jwt.ValidationError)

	// 当 Access Token 是过期错误，并且 Refresh Token 没有过期时就创建一个新的 Access Token
	if v.Errors == jwt.ValidationErrorExpired {
		return _GenToken(claims.UserID, claims.Username)
	}
	return
}
