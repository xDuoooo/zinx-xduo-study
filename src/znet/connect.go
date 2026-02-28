package znet

import (
	"net"
	"zinx-xduo-study/src/ziface"
)

// Connection 当前连接模块
type Connection struct {
	//当前连接的Socket TCP套接字
	Conn *net.TCPConn

	//连接的ID
	ConnID uint32

	//当前连接状态
	isClosed bool

	//当前连接所绑定的处理业务方法API
	handlerAPI ziface.HandleFunc

	//告知当前连接已经退出的/停止 channel
	EXitChan chan bool
}

// NewConnection 初始化连接模块的方法
func NewConnection(conn *net.TCPConn, connID uint32, callBackAPI ziface.HandleFunc) *Connection {
	c := &Connection{
		Conn:       conn,
		ConnID:     connID,
		handlerAPI: callBackAPI,
		isClosed:   false,
		EXitChan:   make(chan bool, 1),
	}
	return c
}
