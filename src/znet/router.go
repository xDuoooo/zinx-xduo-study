package znet

import "zinx-xduo-study/src/ziface"

// BaseRouter 实现router时候，先嵌入这个BaseRouter基类，然后根据需要对这个基类的方法进行重写
type BaseRouter struct{}

//这里之所以BaseRouter的方法都为空
//是因为有的Router不希望有PreHandler、PostHander这两个业务
//所以Router全部继承BaseRouter的好处就是，不需要都实现，可以选择性实现

// PreHande 在处理业务之前的hook
func (br *BaseRouter) PreHande(request ziface.IRequest) {

}

// Handle 在处理业务的hook
func (br *BaseRouter) Handle(request ziface.IRequest) {

}

// PostHandle 在处理业务之后的hook
func (br *BaseRouter) PostHandle(request ziface.IRequest) {

}
