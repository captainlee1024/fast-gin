// Package v1 包含了服务的接口以及各个服务接口需要用到的 DTO
// DTO（Data Transfer Object）：数据传输对象
// 由于不是分布式项目...
package v1

import "github.com/gin-gonic/gin"

// FastGinServer is the server API for FastGin Service
// All implementations must emed UnimplementedShopServer
// for forward compatibility
type FastGinServer interface {
	// Index is the fastginindex service
	Index(c *gin.Context)
	// Error is the fastginerror service
	Error(c *gin.Context)
	// Ping is the fastginping service
	Ping(c *gin.Context)
}
