// 依赖倒置：业务逻辑层依赖于数据持久化层，但是不应该依赖一个实现，应该依赖于它的抽象

package service

import (
	"errors"
	"github.com/captainlee1024/fast-gin/internal/fastgin/do"
)

// FastGinDo 业务实体 放在do 层了
//type FastGinDo struct {
//	DemoName string
//	Info     string
//}

// FastGinDoRepo 存储接口
type FastGinDoRepo interface {
	SaveFastGin(*do.FastGinDo) error
	GetFastGin(*do.FastGinDo) (*do.FastGinDo, error)
	ListFastGin() ([]*do.FastGinDo, error)
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
func (uc *FastGinUsecase) IndexBiz() (data string, err error) {
	return "Welcome to FastGin!", nil
}

// PingPong 处理 ping 业务
func (uc *FastGinUsecase) PingPong() (data string, err error) {
	// logic business
	return "Pong", nil
}

// ErrorBiz 用于测试 ResponseError
func (uc *FastGinUsecase) ErrorBiz() (data string, err error) {
	// logic business
	return "", errors.New("test Middleware ResponseError...")
}

// Create 添加一个 FastGin
func (uc *FastGinUsecase) Create(fdo *do.FastGinDo) (err error) {
	// logic business
	return uc.repository.SaveFastGin(fdo)
}

// Get 获取 FastGin
func (uc *FastGinUsecase) Get(fdo *do.FastGinDo) (fastGin *do.FastGinDo, err error) {
	fastGin = &do.FastGinDo{}
	fastGin, err = uc.repository.GetFastGin(fdo)
	return
}

// ListFastGin 获取所有 FastGin
func (uc *FastGinUsecase) ListFastGin(fdo *do.FastGinDo) (listFastGin []*do.FastGinDo, err error) {

	fastGins, err := uc.repository.ListFastGin()
	if err != nil {
		return nil, err
	}
	// 初始化返回值定义的变量，那里只是声明，并没有申请内存
	listFastGin = make([]*do.FastGinDo, 0, len(fastGins))
	listFastGin = fastGins
	return
}
