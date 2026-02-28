package ziface

import "net"

//定义连接模块的抽象层

type IConnection interface {
	// Start 启动连接 让当前的连接准备开始工作
	Start()
	// Stop 停止连接 结束当前连接的工作
	Stop()
	// GetTCPConnection 获取当前连接绑定的socket connection
	GetTCPConnection() *net.TCPConn
	// GetConnID 获取当前连接模块的连接ID
	GetConnID() uint32
	// RemoteAddr 获取远程客户端的 TCP状态 IP port
	RemoteAddr() net.Addr
	// Send 发送数据，将数据发送给远程的客户端
	Send(data []byte) error
}

// HandleFunc 定义一个处理连接业务的方法
type HandleFunc func(*net.TCPConn, []byte, int) error
