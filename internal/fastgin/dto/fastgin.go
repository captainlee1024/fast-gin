// Package dto DTO（Data Transfer Object）：数据传输对象
package dto

import (
	"github.com/captainlee1024/fast-gin/internal/pkg/public"
	"github.com/gin-gonic/gin"
)

// FastGinDto 传输实体
type FastGinDto struct {
	DemoName string `json:"demo_name" comment:"名称" validate:"require" example:"testing"`
	Info     string `json:"info" comment:"详细信息" validate:"" example:"fast-gin 是一个通用企业脚手架"`
}

// BindingValidParams 验证逻辑
func (params *FastGinDto) BindingValidParams(c *gin.Context) error {
	return public.DefaultGetValidParams(c, params)
}
