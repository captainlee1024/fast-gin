// Package dto DTO（Data Transfer Object）：数据传输对象
package dto

import (
	"github.com/captainlee1024/fast-gin/internal/pkg/public"
	"github.com/gin-gonic/gin"
)

// FastGinDto 传输实体
type FastGinDto struct {
	FastGinID int64  `json:"fast_gin_id" form:"fast_gin_id" comment:"ID" validate:""`
	DemoName  string `json:"demo_name" form:"demo_name" comment:"名称" validate:"required" example:"testing"`
	Info      string `json:"info" form:"info" comment:"详细信息" validate:"" example:"fast-gin 是一个通用企业脚手架"`
}

// BindingValidParams 验证逻辑
func (params *FastGinDto) BindingValidParams(c *gin.Context) error {
	return public.DefaultGetValidParams(c, params)
}

// FastGinListPageDto 分页获取所有 FastGin 简要信息（只展示 ID 和名称）
type FastGinListPageDto struct {
	PageSize  int    `json:"page_size" form:"page_size" comment:"每页记录数" validate:"" example:"10"`
	Page      int    `json:"page" form:"page" comment:"页数" validate:"" example:"1"`
	FastGinID int64  `json:"fast_gin_id" form:"fast_gin_id" comment:"ID" validate:""`
	DemoName  string `json:"demo_name" form:"demo_name" comment:"名称" validate:"" example:"testing"`
}

// FastGinListPageDto 验证逻辑
func (params *FastGinListPageDto) BindingValidParems(c *gin.Context) error {
	return public.DefaultGetValidParams(c, params)
}
