// Package do DO（Domain Object）：领域对象，应该使用贫血模型
package do

// 注意：模型一定要用贫血模型
// 失血模型导致 model 变成了废物，作用非常小
// 充血模型太太太过臃肿
// 贫血模型，使得 model 实现一些验证逻辑，这样可以大大减少逻辑层代码的复杂度，而且不会变得混乱
// 例如一个用户 model ，它的是否被封号，是否注销，是否存在这些逻辑，不应该平铺在你的逻辑层，可以
// 让你的 model 去实现一些验证这些逻辑的方法，然后在逻辑层调用就行了

// FastGinDo 业务实体
type FastGinDo struct {
	FastGinID int64
	DemoName  string
	Info      string
}

// IDIsOk 用于辅助逻辑处理层逻辑校验
func (fastGin *FastGinDo) IDIsOk() bool {
	return fastGin.FastGinID != 0
}

//InfoIsOk 用于辅助逻辑处理层逻辑校验
//这里的校验逻辑是 Info 字段不能为空
func (fastgin *FastGinDo) InfoIsOk() bool {
	return fastgin.Info != ""
}

//DemoNameIsOk 用户名是否存在
func (fastgin *FastGinDo) DemoNameIsOk() bool {
	return fastgin.DemoName != ""
}

// FastGinListPageDo
type FastGinListPageDo struct {
	PageSize  int
	Page      int
	FastGinID int64
	DemoName  string
}

// PageSizeIsOk 是否设置了每页数量
// 如果设置了返回 true，否则返回 false
func (f *FastGinListPageDo) PageSizeIsEmpty() bool {
	return f.PageSize == 0
}

// PageIsOk 是否设置了从第几页开始查询
// 如果设置了返回 true，否则返回 false
func (f *FastGinListPageDo) PageIsEmpty() bool {
	return f.Page == 0
}
