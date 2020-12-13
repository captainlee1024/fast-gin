package data

import (
	"github.com/captainlee1024/fast-gin/internal/fastgin/do"
	"github.com/captainlee1024/fast-gin/internal/fastgin/service"
)

// 在编译的时候可以知道这个对象实现了这个 interface{}
var _ service.FastGinDoRepo = (service.FastGinDoRepo)(nil)

// NewFastGinRepo 创建一个 fastGinRepo ，它是 service.FastGinDoRepo 的实现
func NewFastGinRepo() service.FastGinDoRepo {
	return new(fastGinRepo)
}

type fastGinRepo struct{}

// SaveFastGin 保存 FastGinDo 至数据库
func (fg *fastGinRepo) SaveFastGin(fgDo *do.FastGinDo) (err error) {
	// ...
	return
}

func (fg *fastGinRepo) GetFastGin(fgDo *do.FastGinDo) (fastGin *do.FastGinDo, err error) {
	// ...
	fastGin = &do.FastGinDo{}
	return
}

func (fg *fastGinRepo) ListFastGin() (listFastGin []*do.FastGinDo, err error) {
	// ...
	// 首先初始化返回值定义的变量，那里只是声明，并没有申请内存
	listFastGin = make([]*do.FastGinDo, 0, 2)
	return
}