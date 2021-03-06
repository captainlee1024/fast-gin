// IOC 控制反转，是一种思想 ->它的典型实现方法就是依赖注入
/*
A 是父组件 (业务层/文中小明), B 是子组件(底层/文中手机), 那么 A 依赖 B. B 控制着 A.

A 如果不采用构造函数、属性/接口或者工厂模式等方式, 而是用到 B 时再实例 B, 那么一旦想替换 B 就会在 A 很多地方都要改动.

这个问题很好解决, 就是 A 再封装一层(也就是依赖注入, 即 A 依赖的 B 注入到 A, 那么控制就翻转了), 再封装一层的具体实现也很多样:
• 属性, 把 B 变为 A 的一个属性, Python 中就是这么提倡的;
• 接口, A 不直接调用 B, 而是调用接口, B 去实现接口, 这样 B 即便由下文中的 iphone6 变成 iphonex 也没关系;
• 构造函数和工厂模式, 我没想出合适的例子, 但是大概应该也是类似的意思.
*/
// 推荐使用 google 的库 github.com/google/wire
// 使用依赖注入有如下好处
// 1. 方便单元测试 （例如 上面我们进行测试的时候，可以 mock 一个需要的参数传如就可以进行测试了，方便，好控制）
// 2. 一次初始化多次复用（例如 redis、gRPC 连接，以参数的形式传进来，这样一个连接可以传给多个人进行复用）

package controller

import (
	"errors"
	v1 "github.com/captainlee1024/fast-gin/api/fastgin/v1"
	"github.com/captainlee1024/fast-gin/internal/fastgin/data"
	"github.com/captainlee1024/fast-gin/internal/fastgin/do"
	"github.com/captainlee1024/fast-gin/internal/fastgin/dto"
	"github.com/captainlee1024/fast-gin/internal/fastgin/middleware"
	"github.com/captainlee1024/fast-gin/internal/fastgin/service"
	"github.com/gin-gonic/gin"
)

// FastGinRegister 注册路由处理器
// 这种方式不方便给同一个 controller 文件里的各个服务进行分组，只能分到同一组
func FastGinRegister(router *gin.RouterGroup) {
	repo := data.NewFastGinRepo()
	fus := service.NewFastGinUsercase(repo)
	fastGin := NewFastGinController(fus)

	router.GET("/", fastGin.IndexHandler)
	router.GET("error", fastGin.ErrorHandler)
	router.GET("/ping", fastGin.PingHandler)
	router.GET("/listpage", fastGin.GetFastGinListHandler)
	router.POST("/add", fastGin.AddFastGinHandler)
	router.POST("/get", fastGin.GetFastGinHandler)
	router.POST("/edit", fastGin.EditFastGinHandler)
	router.POST("/remove", fastGin.RemoveFastGinHandler)

}

// FastGinController 是 FastGin API 的实现类。
// must be embedded to have forward compatible implementations.
type FastGinController struct {
	v1.UnimplementedFastGinServer
	fus *service.FastGinUsecase
}

// NewFastGinController .
// NewFastGinController
// @Summary Index 接口
// @Description 这是 FastGin 的欢迎页面接口
// @Tags FastGinDemo
// @Accept application/json
// @Produce application/json
// @Security
// @Success 200
// @Router / [get]
func NewFastGinController(fus *service.FastGinUsecase) *FastGinController {
	return &FastGinController{fus: fus}
}

//func NewFastGinController2() *FastGinController {
//	return NewFastGinController(service.NewFastGinUsercase(data.NewFastGinRepo()))
//}

// IndexHandler 处理器
func (fastGin *FastGinController) IndexHandler(c *gin.Context) {
	bizData, err := fastGin.fus.IndexBiz(c)
	if err != nil {
		middleware.ResponseError(c, 2000, err)
	}
	middleware.ResponseSuccess(c, bizData)
}

// PingHandler 处理器
func (fastGin *FastGinController) PingHandler(c *gin.Context) {
	bizData, err := fastGin.fus.PingPong(c)
	if err != nil {
		middleware.ResponseError(c, 2000, err)
	}
	middleware.ResponseSuccess(c, bizData)
}

// ErrorHandler 处理器
func (fastGin *FastGinController) ErrorHandler(c *gin.Context) {
	bizData, err := fastGin.fus.ErrorBiz(c)
	if err != nil {
		middleware.ResponseError(c, 2000, err)
	}
	middleware.ResponseSuccess(c, bizData)
}

// AddFastGinHandler 创建一个 FastGin
func (fastGin *FastGinController) AddFastGinHandler(c *gin.Context) {
	fgDto := new(dto.FastGinDto)
	if err := fgDto.BindingValidParams(c); err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}

	fgDo := &do.FastGinDo{
		DemoName: fgDto.DemoName,
		Info:     fgDto.Info,
	}
	err := fastGin.fus.AddFastGin(fgDo, c)
	if err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}
	middleware.ResponseSuccess(c, "Add FastGin success!")
}

// GetFastGinHandler .
func (fastGin *FastGinController) GetFastGinHandler(c *gin.Context) {
	fgDto := new(dto.FastGinDto)
	if err := fgDto.BindingValidParams(c); err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}

	fgDo := &do.FastGinDo{
		FastGinID: fgDto.FastGinID,
		DemoName:  fgDto.DemoName,
		Info:      fgDto.Info,
	}

	fgDo, err := fastGin.fus.GetFastGin(fgDo, c)
	if err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}
	responseDto := &dto.FastGinDto{
		DemoName: fgDo.DemoName,
		Info:     fgDo.Info,
	}

	middleware.ResponseSuccess(c, responseDto)
}

// GetFastGinListHandler .
func (fastGin *FastGinController) GetFastGinListHandler(c *gin.Context) {

	//// 分页处理
	//if !fgListDo.PageIsOk() {
	//	page, _ := public.GetPageInfo(c)
	//	fgListDo.Page = page
	//}
	//if !fgListDo.PageSizeIsOk() {
	//	_, size := public.GetPageInfo(c)
	//	fgListDo.PageSize = size
	//}

	fgListDto := new(dto.FastGinListPageDto)
	if err := fgListDto.BindingValidParems(c); err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}

	fgListDo := &do.FastGinListPageDo{
		Page:     fgListDto.Page,
		PageSize: fgListDto.PageSize,
	}
	fgDtos, err := fastGin.fus.ListFastGin(fgListDo, c)
	if err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}
	//fgDto := new(fgd)

	middleware.ResponseSuccess(c, fgDtos)
}

// EditFastGinHandler 修改 FastGin 信息
func (fastGin *FastGinController) EditFastGinHandler(c *gin.Context) {
	fgDto := new(dto.FastGinDto)
	if err := fgDto.BindingValidParams(c); err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}

	fgDo := &do.FastGinDo{
		FastGinID: fgDto.FastGinID,
		DemoName:  fgDto.DemoName,
		Info:      fgDto.Info,
	}

	err := fastGin.fus.EditFastGin(fgDo, c)
	if err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}

	middleware.ResponseSuccess(c, "success! 信息修改成功！")
}

// RemoveFastGinHandler 删除 FastGin
func (fastGin *FastGinController) RemoveFastGinHandler(c *gin.Context) {
	fgDto := new(dto.FastGinDto)
	if err := fgDto.BindingValidParams(c); err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}

	fgDo := &do.FastGinDo{
		FastGinID: fgDto.FastGinID,
	}

	err := fastGin.fus.RemoveFastGin(fgDo, c)
	if err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}
	middleware.ResponseSuccess(c, "success! 删除成功！")
}

// GetJWTTokenHandler .
func (fastGin *FastGinController) GetJWTTokenHandler(c *gin.Context) {
	// 1. 参数校验
	fgDto := new(dto.FastGinDto)
	if err := fgDto.BindingValidParams(c); err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}

	// 对象转换
	// dto -> do
	fgDo := &do.FastGinDo{
		FastGinID: fgDto.FastGinID,
	}

	fgDo, err := fastGin.fus.GetJWTToken(fgDo, c)
	if err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}

	// 对象转换
	// dto -> do
	responseDto := &dto.FastGinDto{
		FastGinID: fgDo.FastGinID,
		DemoName:  fgDo.DemoName,
	}
	middleware.ResponseSuccess(c, responseDto)
}

// AuthFastGinJWTHandler 测试 JWT 认证中间件的服务
func (fastGin *FastGinController) AuthFastGinJWTHandler(c *gin.Context) {
	// 1. 参数校验
	fgDto := new(dto.FastGinDto)
	if err := fgDto.BindingValidParams(c); err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}

	// 2. 对象转换
	// dto -> do
	fgDo := &do.FastGinDo{
		FastGinID: fgDto.FastGinID,
		DemoName:  fgDto.DemoName,
	}

	// 3. 逻辑处理
	err := fastGin.fus.AuthJWT(fgDo, c)
	if err != nil {
		middleware.ResponseError(c, 2001, errors.New("认证失败"))
		return
	}
	middleware.ResponseSuccess(c, "你已经通过了JWT的认证哦！")
}

// AuthFastGinSessionHandler 测试 Session 认证中间件的服务
func (fastGin *FastGinController) AuthFastGinSessionHandler(c *gin.Context) {
	middleware.ResponseSuccess(c, "你已经通过了 Session 认证哦！")
}
