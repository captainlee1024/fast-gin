// Package v1 包含了服务的接口以及各个服务接口需要用到的 DTO
// DTO（Data Transfer Object）：数据传输对象
// 这里把 DTO 导读放到一个包内 /dto
package v1

import (
	"errors"
	"github.com/captainlee1024/fast-gin/internal/fastgin/middleware"
	"github.com/gin-gonic/gin"
)

// FastGinServer is the server API for FastGin Service
// All implementations must emed UnimplementedFastGinServer
// for forward compatibility
type FastGinServer interface {
	// IndexHandler is the fastginindex service
	IndexHandler(c *gin.Context)
	// ErrorHandler is the fastginerror service
	ErrorHandler(c *gin.Context)
	// PingHandler is the fastginping service
	PingHandler(c *gin.Context)
	// GetFastGinListHandler is the listpage service
	GetFastGinListHandler(c *gin.Context)
	// AddFastGinHandler is the add service
	AddFastGinHandler(c *gin.Context)
	// GetFastGinHandler is the get a FastGin service
	GetFastGinHandler(c *gin.Context)
	// EditFastGinHandler is the .
	EditFastGinHandler(c *gin.Context)
	// RemoveFastGinHandler is the .
	RemoveFastGinHandler(c *gin.Context)

	// 认证相关
	// GetJWTTolenHandler 获取 JWT
	GetJWTTokenHandler(c *gin.Context)
	// AuthFastGinJWTHandler 通过 JWT 认证的服务
	AuthFastGinJWTHandler(c *gin.Context)
	// AuthFastGinSessionHandler 通过 Session 认证的服务
	AuthFastGinSessionHandler(c *gin.Context)
}

// UnimplementedFastGinServer can be embedded to have forward compatible implementations.
type UnimplementedFastGinServer struct {
}

func (*UnimplementedFastGinServer) IndexHandler(c *gin.Context) {
	middleware.ResponseError(c, 2001, errors.New("method IndexHandler not implemented"))
}

func (*UnimplementedFastGinServer) ErrorHandler(c *gin.Context) {
	middleware.ResponseError(c, 2001, errors.New("method ErrorHandler not implemented"))
}

func (*UnimplementedFastGinServer) PingHandler(c *gin.Context) {
	middleware.ResponseError(c, 2001, errors.New("method PingHandler not implemented"))
}

func (*UnimplementedFastGinServer) GetFastGinListHandler(c *gin.Context) {
	middleware.ResponseError(c, 2001, errors.New("method GetFastGinListHandler not implemented"))
}

func (*UnimplementedFastGinServer) AddFastGinHandler(c *gin.Context) {
	middleware.ResponseError(c, 2001, errors.New("method AddFastGinHandler not implemented"))
}

func (*UnimplementedFastGinServer) GetFastGinHandler(c *gin.Context) {
	middleware.ResponseError(c, 2001, errors.New("method GetFastGinHandler not implemented"))
}

func (*UnimplementedFastGinServer) EditFastGinHandler(c *gin.Context) {
	middleware.ResponseError(c, 2001, errors.New("method EditFastGinHandler not implemented"))
}

func (*UnimplementedFastGinServer) RemoveFastGinHandler(c *gin.Context) {
	middleware.ResponseError(c, 2001, errors.New("method RemoveFastGinHandler not implemented"))
}

func (*UnimplementedFastGinServer) GetJWTTokenHandler(c *gin.Context) {
	middleware.ResponseError(c, 2001, errors.New("method GetJWTTokenHandler not implemented"))
}

func (*UnimplementedFastGinServer) AuthFastGinJWTHandler(c *gin.Context) {
	middleware.ResponseError(c, 2001, errors.New("method AuthFastGinJWTHandler not implemented"))
}

func (*UnimplementedFastGinServer) AuthFastGinSessionHandler(c *gin.Context) {
	middleware.ResponseError(c, 2001, errors.New("method AuthFastGinSessionHandler not implemented"))
}
