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
//func (br *PingRouter) PreHandle(request ziface.IRequest) {
//	fmt.Println("Call Router PreHandle...")
//	_, err := request.GetConnection().GetTCPConnection().Write([]byte("before ping..."))
//	if err != nil {
//		fmt.Println("call back before ping error")
//	}
//}

// Handle Test Handle
func (br *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call PingRouter Handle...")
	//先读取客户端的数据，再回写ping..ping..ping
	fmt.Println("receive from client: msgID = ", request.GetMsgID())
	fmt.Println("receive from client: message = ", string(request.GetData()))

	err := request.GetConnection().SendMessage(101, []byte("PingRouter..PingRouter..PingRouter"))
	if err != nil {
		fmt.Println(err)
	}
}

// PostHandle Test PostHandle
//
//	func (br *PingRouter) PostHandle(request ziface.IRequest) {
//		fmt.Println("Call Router PostHandle...")
//		_, err := request.GetConnection().GetTCPConnection().Write([]byte("Post ping..."))
//		if err != nil {
//			fmt.Println("call back Post ping error")
//		}
//	}

// DoConnBegin 创建连接之后执行钩子函数
func DoConnBegin(conn ziface.IConnection) {
	fmt.Println("---> DoConnectionBegin is Called ...")
	err := conn.SendMessage(202, []byte("DoConnection BEGIN "))
	if err != nil {
		fmt.Println(err)
	}
	//给当前的连接设置一些属性
	fmt.Println("Set conn Name,Home ...")
	conn.SetProperty("Name", "xduo-study")
	conn.SetProperty("github", "https://github.com/xDuoooo")
}

func DoAfterConnDeploy(conn ziface.IConnection) {
	fmt.Println("--> DoAfterConnDeploy is Called ...")
	fmt.Println("conn ID = ", conn.GetConnID(), "is deploy")

	property, err := conn.GetProperty("Name")
	if err == nil {
		fmt.Println("Name = ", property)
	}
}

func main() {
	//1. 创建一个server句柄 使用 Zinx的Api
	s := znet.NewServer("[Zinx V0.8]")
	//2. 自定义路由
	pingRouter := &PingRouter{}
	s.AddRouter(0, pingRouter)
	helloRouter := &HelloZinxRouter{}
	s.AddRouter(1, helloRouter)
	//3. 注册连接Hook钩子函数
	s.SetBeforeConnCreateFunc(DoConnBegin)
	s.SetAfterConnDeployFunc(DoAfterConnDeploy)

	//4. 启动Server

	s.Server()
}

type HelloZinxRouter struct {
	znet.BaseRouter
}

// PreHandle Test PreHandle
//func (br *PingRouter) PreHandle(request ziface.IRequest) {
//	fmt.Println("Call Router PreHandle...")
//	_, err := request.GetConnection().GetTCPConnection().Write([]byte("before ping..."))
//	if err != nil {
//		fmt.Println("call back before ping error")
//	}
//}

// Handle Test Handle
func (br *HelloZinxRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call HelloZinxRouter Handle...")
	//先读取客户端的数据，再回写ping..ping..ping
	fmt.Println("receive from client: msgID = ", request.GetMsgID())
	fmt.Println("receive from client: message = ", string(request.GetData()))

	err := request.GetConnection().SendMessage(201, []byte("helloRouter..helloRouter..helloRouter"))
	if err != nil {
		fmt.Println(err)
	}
}
