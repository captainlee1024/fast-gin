package dto

import (
	"github.com/captainlee1024/fast-gin/pkg/public"
	"github.com/gin-gonic/gin"
)

// ParamSigUp 注册请求参数
type ParamSigUp struct {
	Username string `json:"username" binding:"required"`               // 用户名
	Password string `json:"password" binding:"required,checkPassword"` // 用户密码
	//RePassword string `json:"re_password" binding:"required,eqfield=Password"` // 确认密码
	RePassword string `json:"confirm_password" binding:"required,eqfield=Password"` // 确认密码
}

// ParamLogin 登录请求参数
type ParamLogin struct {
	Username string `json:"username" binding:"required"` // 用户名
	Password string `json:"password" binding:"required"` // 用户密码
}

// ListPageInput 分页参数
type ListPageInput struct {
	PageSize int    `form:"page_size" json:"page_size" comment:"每页记录数" validate:"" example:"10"`
	Page     int    `form:"page" json:"page" comment:"页数" validate:"required" example:"1"`
	Name     string `form:"name" json:"name" comment:"姓名" validate:"" example:""`
}

// BindingValidParams 验证逻辑
func (params *ListPageInput) BindingValidParams(c *gin.Context) error {
	return public.DefaultGetValidParams(c, params)
}

// AddUserInput 添加用户的 DTO
type AddUserInput struct {
	Name  string `form:"name" json:"name" comment:"姓名" validate:"required"`
	Sex   int    `form:"sex" json:"sex" comment:"性别" validate:""`
	Age   int    `form:"age" json:"age" comment:"年龄" validate:"required,gt=10"`
	Birth string `form:"birth" json:"birth" comment:"生日" validate:"required"`
	Addr  string `form:"addr" json:"addr" comment:"地址" validate:"required"`
}

// BindingValidParams 添加用户字段验证逻辑
func (params *AddUserInput) BindingValidParams(c *gin.Context) error {
	return public.DefaultGetValidParams(c, params)
}

// EditUserInput 编辑/修改用户的 DTO
type EditUserInput struct {
	Id    int    `form:"id" json:"id" comment:"ID" validate:"required"`
	Name  string `form:"name" json:"name" comment:"姓名" validate:"required"`
	Sex   int    `form:"sex" json:"sex" comment:"性别" validate:""`
	Age   int    `form:"age" json:"age" comment:"年龄" validate:"required,gt=10"`
	Birth string `form:"birth" json:"birth" comment:"生日" validate:"required"`
	Addr  string `form:"addr" json:"addr" comment:"地址" validate:"required"`
}

// BindingValidParams 修改用户信息字段的验证逻辑
func (params *EditUserInput) BindingValidParams(c *gin.Context) error {
	return public.DefaultGetValidParams(c, params)
}

// RemoveUserInput 删除用户的 DTO
type RemoveUserInput struct {
	IDS string `form:"ids" json:"ids" comment:"IDS" validate:"required"`
}

// BindingValidParams 删除用户字段的验证逻辑
func (params *RemoveUserInput) BindingValidParams(c *gin.Context) error {
	return public.DefaultGetValidParams(c, params)
}
