// 依赖倒置：业务逻辑层依赖于数据持久化层，但是不应该依赖一个实现，应该依赖于它的抽象

package service

import (
	"errors"
	"fmt"
	"github.com/captainlee1024/fast-gin/internal/fastgin/do"
	"github.com/captainlee1024/fast-gin/internal/pkg/public"
	"github.com/captainlee1024/fast-gin/pkg/jwt"
	"github.com/gin-gonic/gin"
)

// FastGinDo 业务实体 放在do 层了
//type FastGinDo struct {
//	DemoName string
//	Info     string
//}

// FastGinDoRepo 存储接口
type FastGinDoRepo interface {
	SaveFastGin(*do.FastGinDo, *gin.Context) error
	GetFastGinByID(int64, *gin.Context) (*do.FastGinDo, error)
	GetFastGin(*do.FastGinDo, *gin.Context) (*do.FastGinDo, error)
	GetFastGinList(int, int, *gin.Context) ([]*do.FastGinDo, error)
	UpdateFastGin(*do.FastGinDo, *gin.Context) error
	DeleteFastGin(*do.FastGinDo, *gin.Context) error
	SetAToken(string, string, *gin.Context) error
}

// FastGinUsecase is .
type FastGinUsecase struct {
	repository FastGinDoRepo
}

// NewFastGinUsercase 创建一个 FastGinUsercase
func NewFastGinUsercase(repo FastGinDoRepo) *FastGinUsecase {
	return &FastGinUsecase{repository: repo}
}

// IndexBiz fastgin 欢迎页面
func (uc *FastGinUsecase) IndexBiz(c *gin.Context) (data string, err error) {
	return "Welcome to FastGin!", nil
}

// PingPong 处理 ping 业务
func (uc *FastGinUsecase) PingPong(c *gin.Context) (data string, err error) {
	// logic business
	return "Pong", nil
}

// ErrorBiz 用于测试 ResponseError
func (uc *FastGinUsecase) ErrorBiz(c *gin.Context) (data string, err error) {
	// logic business
	return "", errors.New("test Middleware ResponseError...")
}

// Create 添加一个 FastGin
func (uc *FastGinUsecase) AddFastGin(fdo *do.FastGinDo, c *gin.Context) (err error) {
	// logic business
	if !fdo.DemoNameIsOk() && !fdo.InfoIsOk() {
		return errors.New("AddFastGin fdo DemoName/Info 不能为空！")
	}
	return uc.repository.SaveFastGin(fdo, c)
}

// Get 获取 FastGin
func (uc *FastGinUsecase) GetFastGin(fdo *do.FastGinDo, c *gin.Context) (fastGin *do.FastGinDo, err error) {
	fastGin = &do.FastGinDo{}
	if !fdo.InfoIsOk() && !fdo.DemoNameIsOk() {
		return nil, errors.New("GetFastGin fdo DemoName/Info 至少有一个不为空")
	}
	fastGin, err = uc.repository.GetFastGin(fdo, c)
	return
}

// ListFastGin 获取所有 FastGin
func (uc *FastGinUsecase) ListFastGin(fgListDo *do.FastGinListPageDo, c *gin.Context) (listFastGin []*do.FastGinDo, err error) {
	// 分页处理
	if fgListDo.PageIsEmpty() {
		page, _ := public.GetPageInfo(c)
		fgListDo.Page = page
	}
	if fgListDo.PageSizeIsEmpty() {
		_, size := public.GetPageInfo(c)
		fgListDo.PageSize = size
	}

	fastGins, err := uc.repository.GetFastGinList(fgListDo.Page, fgListDo.PageSize, c)
	if err != nil {
		return nil, err
	}

	// 初始化返回值定义的变量，那里只是声明，并没有申请内存
	listFastGin = make([]*do.FastGinDo, 0, len(fastGins))
	//for _, fastGin := range fastGins {
	//	listFastGin = append(listFastGin, fastGin)
	//}
	listFastGin = fastGins
	return listFastGin, err
}

// EditFastGin 修改信息
func (uc *FastGinUsecase) EditFastGin(fdo *do.FastGinDo, c *gin.Context) (err error) {
	return uc.repository.UpdateFastGin(fdo, c)

}

// RemvoeFastGin 删除一个 FastGin
func (uc *FastGinUsecase) RemoveFastGin(fdo *do.FastGinDo, c *gin.Context) (err error) {
	return uc.repository.DeleteFastGin(fdo, c)
}

func (uc *FastGinUsecase) GetJWTToken(fdo *do.FastGinDo, c *gin.Context) (fastGin *do.FastGinDo, err error) {
	if !fdo.IDIsOk() {
		return nil, errors.New("必须填写正确的　ID 才能获取通过认证的 JWTToken！")
	}
	fastGin = new(do.FastGinDo)
	fastGin, err = uc.repository.GetFastGinByID(fdo.FastGinID, c)
	if err != nil {
		return nil, err
	}

	token, err := jwt.GenToken(fastGin.FastGinID, fastGin.DemoName)
	if err != nil {
		return nil, err
	}
	err = uc.repository.SetAToken(fmt.Sprint(fastGin.FastGinID), token, c)
	if err != nil {
		return nil, err
	}
	return
}

// AuthJWT 认证逻辑
func (uc *FastGinUsecase) AuthJWT(fdo *do.FastGinDo, c *gin.Context) error {
	return nil
}
