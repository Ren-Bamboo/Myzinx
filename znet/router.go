package znet

import (
	"github.com/Ren-Bamboo/Myzinx/ziface"
)

// 创建的目的是：有些Router不需要实现IRouter接口中所有方法，因此，创建一个Base，然后继承即可
type BaseRouter struct {
}

// 处理Request（Connection）业务前的方法
func (br *BaseRouter) PreHandle(request ziface.IRequest) {}

// 处理Request（Connection）业务的方法
func (br *BaseRouter) Handle(request ziface.IRequest) {
}

// 处理Request（Connection）业务后的方法
func (br *BaseRouter) PastHandle(request ziface.IRequest) {}
