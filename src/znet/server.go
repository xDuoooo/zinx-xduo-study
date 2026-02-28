package znet

import "zinx-xduo-study/src/ziface"

// Server IServer 的接口实现，定义一个Server的服务器模块
type Server struct {
	//服务器名称
	Name string
	//服务器绑定的IP版本
	IPVersion string
	//服务器监听的IP
	IP string
	//服务器监听的端口
	Port int
}

func (s *Server) Start() {

}
func (s *Server) Stop() {

}
func (s *Server) Server() {

}

/*
初始化Server模块的方法
*/
func NewServer(name string) ziface.IServer {
	s := &Server{
		Name:      name,
		IPVersion: "icp4",
		IP:        "0.0.0.0",
		Port:      8999,
	}
	return s
}
