package znet

import (
	"fmt"
	"strconv"
	"zinx-xduo-study/src/utils"
	"zinx-xduo-study/src/ziface"
)

// MsgHandler 消息处理模块的实现
type MsgHandler struct {
	//存放每一个MsgID 所对应的处理方法
	Apis map[uint32]ziface.IRouter

	//当前业务工作Worker的数量
	WorkerPoolSize uint32
	//负责Worker取任务的消息队列
	TaskQueue []chan ziface.IRequest
}

// DoMsgHandler 调度 执行对应的Router消息处理方法
func (mh *MsgHandler) DoMsgHandler(request ziface.IRequest) {
	router, ok := mh.Apis[request.GetMsgID()]
	if !ok {
		fmt.Println("api msgID = ", request.GetMsgID(), "IS NOT FOUND Need Register")
		return // 路由不存在时直接返回，防止 nil pointer panic
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
		Apis:           make(map[uint32]ziface.IRouter),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize,
		TaskQueue:      make([]chan ziface.IRequest, utils.GlobalObject.WorkerPoolSize),
	}
}

// StartWorkerPool 启动一个Worker工作池(开启工作池动作只能发生一次，一个zinx框架只能有一个worker工作池)
func (mh *MsgHandler) StartWorkerPool() {
	//根据WorkerPoolSize 分别开启Worker，每个Worker用一个Go来承载
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		//一个worker被启动
		//1. 当前worker对应的channel消息队列开辟空间 第 0 个worker 对应第0个chan
		mh.TaskQueue[i] = make(chan ziface.IRequest, utils.GlobalObject.MaxWorkerTaskLen)
		//2. 启动当前的Worker，阻塞等待进行消息消费
		go mh.StartOneWorker(i, mh.TaskQueue[i])
	}
}

// StartOneWorker 启动一个Worker工作流程
func (mh *MsgHandler) StartOneWorker(workerID int, taskQueue chan ziface.IRequest) {
	fmt.Println("Worker ID = ", workerID, "is started ... ")
	//不断阻塞等待对应消息队列的消息
	for {
		select {
		case request := <-taskQueue:
			mh.DoMsgHandler(request)
		}
	}
}

func (mh *MsgHandler) SendMsgToTaskQueue(request ziface.IRequest) {
	//1. 将消息平均分配给不同的worker
	workerID := request.GetConnection().GetConnID() % mh.WorkerPoolSize
	fmt.Println("Add ConnID = ", request.GetConnection().GetConnID(),
		"request MsgID = ", request.GetMsgID(),
		"to WorkerID = ", workerID)
	//2. 将消息发送给对应的worker的TaskQueue即可
	mh.TaskQueue[workerID] <- request
}
