package main

import (
	"fmt"
	"zinx-xduo-study/src/ziface"
	"zinx-xduo-study/src/znet"
)

/*
基于Zinx框架来开发的 服务端应用程序
*/

type PingRouter struct {
	znet.BaseRouter
}

// PreHandle Test PreHandle
func (br *PingRouter) PreHandle(request ziface.IRequest) {
	fmt.Println("Call Router PreHandle...")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("before ping..."))
	if err != nil {
		fmt.Println("call back before ping error")
	}
}

// Handle Test Handle
func (br *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call Router Handle...")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("ping...ping...ping"))
	if err != nil {
		fmt.Println("call back ping error")
	}
}

// PostHandle Test PostHandle
func (br *PingRouter) PostHandle(request ziface.IRequest) {
	fmt.Println("Call Router PostHandle...")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("Post ping..."))
	if err != nil {
		fmt.Println("call back Post ping error")
	}
}
func main() {
	//1. 创建一个server句柄 使用 Zinx的Api
	s := znet.NewServer("[Zinx V0.2]")
	//2. 自定义路由
	router := PingRouter{}
	s.AddRouter(&router)
	//3. 启动Server
	s.Server()
}
