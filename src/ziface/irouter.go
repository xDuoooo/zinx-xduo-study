package ziface

// IRouter 路由抽象接口
// 路由里的数据都是IRequest
type IRouter interface {
	// PreHande 在处理业务之前的hook
	PreHande(request IRequest)
	// Handle 在处理业务的hook
	Handle(request IRequest)
	// PostHandle 在处理业务之后的hook
	PostHandle(request IRequest)
}
