package znet

import (
	"fmt"
	"strconv"
	"zinx-xduo-study/src/ziface"
)

// MsgHandler 消息处理模块的实现
type MsgHandler struct {
	//存放每一个MsgID 所对应的处理方法
	Apis map[uint32]ziface.IRouter
}

// DoMsgHandler 调度 执行对应的Router消息处理方法
func (mh *MsgHandler) DoMsgHandler(request ziface.IRequest) {
	router, ok := mh.Apis[request.GetMsgID()]
	if !ok {
		fmt.Println("api msgID = ", request.GetMsgID(), "IS NOT FOUND Need Register")
	}
	router.PreHandle(request)
	router.Handle(request)
	router.PostHandle(request)
}

// AddRouter 为消息增加具体的处理机制
func (mh *MsgHandler) AddRouter(msgID uint32, router ziface.IRouter) {
	//当前ID已经被注册，就没有必要添加了
	if _, ok := mh.Apis[msgID]; ok {
		//id已经注册了
		panic("repeat api , msgID = " + strconv.Itoa(int(msgID)))
	}
	mh.Apis[msgID] = router
	fmt.Println("Add api MsgID = ", msgID, " success!!!")
}

func NewMsgHandler() *MsgHandler {
	return &MsgHandler{
		Apis: make(map[uint32]ziface.IRouter),
	}
}
