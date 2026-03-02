package ziface

//消息管理抽象层

type IMsgHandler interface {
	// DoMsgHandler 调度 执行对应的Router消息处理方法
	DoMsgHandler(request IRequest)
	// AddRouter 为消息增加具体的处理机制
	AddRouter(msgID uint32, router IRouter)
}
