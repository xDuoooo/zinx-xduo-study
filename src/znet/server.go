package znet

import (
	"fmt"
	"net"
	"zinx-xduo-study/src/ziface"
)

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
	go func() {
		fmt.Printf("[Start] Server Listener at IP :%s, Port %d, is starting", s.IP, s.Port)
		//1. 获取一个TCP的Addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr error: ", err)
		}
		//2. 监听服务器的地址
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listen", s.IPVersion, "err", err)
		}
		fmt.Println("start Zinx server success", s.Name, "success, Listening...")
		//3. 阻塞等待客户端进行连接，处理客户端连接业务()
		for {
			//如果有客户端连接进来，阻塞会返回
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err", err)
				continue
			}
			//已经与客户端建立连接，做一些业务,做一个最基本的最大512字节长度的回显业务
			go func() {
				for {
					buf := make([]byte, 512)
					cnt, err := conn.Read(buf)
					if err != nil {
						fmt.Println("receive buf err", err)
						continue
					}
					//接收到服务器写的消息
					fmt.Printf("receive client buf %s,cnt = %d\n", buf[:cnt], cnt)
					//回显功能
					if _, err := conn.Write(buf[:cnt]); err != nil {
						fmt.Println("write back buf err", err)
						continue
					}
				}
			}()
		}
	}()
}
func (s *Server) Stop() {
	//TODO 将一些服务器的资源、状态或者一些已经开辟的连接信息 进行停止或者回收
}
func (s *Server) Server() {
	//启动Server的服务功能
	s.Start()

	//TODO 做一些启动服务器之后的额外业务

	//阻塞状态
	select {}
}

/*
初始化Server模块的方法
*/
func NewServer(name string) ziface.IServer {
	s := &Server{
		Name:      name,
		IPVersion: "tcp4",
		IP:        "0.0.0.0",
		Port:      8999,
	}
	return s
}
