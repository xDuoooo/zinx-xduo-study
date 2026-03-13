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
	// SendMessage 发送数据，将数据发送给远程的客户端
	SendMessage(uint32, []byte) error
	// SetProperty 设置连接属性
	SetProperty(key string, value interface{})
	// GetProperty 获取
	GetProperty(key string) (interface{}, error)
	// RemoveProperty 移除
	RemoveProperty(key string)
}

// HandleFunc 定义一个处理连接业务的方法
type HandleFunc func(*net.TCPConn, []byte, int) error
