package main

import "zinx-xduo-study/src/znet"

/*
基于Zinx框架来开发的 服务端应用程序
*/
func main() {
	//1. 创建一个server句柄 使用 Zinx的Api
	s := znet.NewServer("[Zinx V0.1]")
	//2. 启动Server
	s.Server()
}
