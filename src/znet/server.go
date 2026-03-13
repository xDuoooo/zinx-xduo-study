package znet

import (
	"fmt"
	"net"
	"zinx-xduo-study/src/utils"
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
	//当前Server绑定的Router
	MsgHandler ziface.IMsgHandler
	//当前server的连接管理器
	connManager ziface.IConnManager
	//连接创建前调用的函数
	BeforeConnCreateFunc func(connection ziface.IConnection)
	//连接销毁所调用的函数
	AfterConnDeployFunc func(connection ziface.IConnection)
}

func (s *Server) AddRouter(msgID uint32, router ziface.IRouter) {
	s.MsgHandler.AddRouter(msgID, router)
	fmt.Println("Add Router Success!")
}

func (s *Server) Start() {

	fmt.Printf("[START] Server name: %s,listenner at IP: %s, Port %d is starting\n", s.Name, s.IP, s.Port)
	fmt.Printf("[Zinx] Version: %s, MaxConn: %d, MaxPacketSize: %d\n",
		utils.GlobalObject.Version,
		utils.GlobalObject.MaxConn,
		utils.GlobalObject.MaxPackageSize)
	go func() {
		//0. 开启消息队列及其Worker工作池
		s.MsgHandler.StartWorkerPool()

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
		var cid uint32
		cid = 0
		//3. 阻塞等待客户端进行连接，处理客户端连接业务()
		for {
			//如果有客户端连接进来，阻塞会返回
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err", err)
				continue
			}
			//判断当前已有连接是否超过最大连接个数,如果超过最大连接，那么则关闭此连接
			if s.connManager.Len() >= utils.GlobalObject.MaxConn {
				fmt.Println("too many server connections !!! MaxConn = ", utils.GlobalObject.MaxConn)
				conn.Close()
				//TODO 给用户发一个失败的消息
				continue
			}

			//将处理新连接的业务方法 和 conn 进行绑定 得到我们的连接模块
			dealConnection := NewConnection(s, conn, cid, s.MsgHandler)
			cid++
			//启动当前的连接业务处理
			go dealConnection.Start()
		}
	}()
}
func (s *Server) Stop() {
	//将一些服务器的资源、状态或者一些已经开辟的连接信息 进行停止或者回收
	fmt.Println("[STOP] Zinx Server name: ", s.Name)
	s.connManager.ClearConn()
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
	// 先 Reload 从配置文件读取配置（如果文件不存在则使用默认值）
	// 必须在初始化其他组件之前调用，这样 WorkerPoolSize 等配置才会生效
	// 注意: 如果 conf/zinx.json 不存在, Reload 内部的 panic 会中止程序，建议写入配置文件后使用
	utils.GlobalObject.Reload()
	s := &Server{
		// 使用传入的 name 参数，而非始终读取全局配置
		Name:        name,
		IPVersion:   "tcp4",
		IP:          utils.GlobalObject.Host,
		Port:        utils.GlobalObject.TCPPort,
		MsgHandler:  NewMsgHandler(),
		connManager: NewConnManager(),
	}
	return s
}
func (s *Server) GetConnManager() ziface.IConnManager {
	return s.connManager
}

func (s *Server) SetBeforeConnCreateFunc(hookFunc func(connection ziface.IConnection)) {
	s.BeforeConnCreateFunc = hookFunc
}
func (s *Server) SetAfterConnDeployFunc(hookFunc func(connection ziface.IConnection)) {
	s.AfterConnDeployFunc = hookFunc
}

func (s *Server) CallBeforeConnCreateFunc(connection ziface.IConnection) {
	if s.BeforeConnCreateFunc != nil {
		fmt.Println("---> Call BeforeConnCreateFunc ...")
		s.BeforeConnCreateFunc(connection)
	}
}

func (s *Server) CallAfterConnDeployFunc(connection ziface.IConnection) {
	if s.AfterConnDeployFunc != nil {
		fmt.Println("---> Call AfterConnDeployFunc ...")
		s.AfterConnDeployFunc(connection)
	}
}
