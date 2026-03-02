package znet

import (
	"errors"
	"fmt"
	"io"
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

	//告知当前连接已经退出的/停止 channel
	EXitChan chan bool

	//该连接处理的方法Router
	Router ziface.IRouter
}

// NewConnection 初始化连接模块的方法
func NewConnection(conn *net.TCPConn, connID uint32, router ziface.IRouter) *Connection {
	c := &Connection{
		Conn:     conn,
		ConnID:   connID,
		isClosed: false,
		EXitChan: make(chan bool, 1),
		Router:   router,
	}
	return c
}

// StratReader 连接的读数据业务方法
func (c *Connection) StratReader() {
	fmt.Println("Reader Goroutine is running...")
	defer fmt.Println("connID = ", c.ConnID, "Reader is exit,remove addr is ", c.RemoteAddr().String())
	defer c.Stop()
	for {
		//读取客户端的数据到buf中，最大512字节
		//buf := make([]byte, )
		//_, err := c.Conn.Read(buf)
		//if err != nil {
		//	fmt.Println("receive buf err", err)
		//	continue
		//}
		////得到当前连接数据的Request数据
		//req := Request{
		//	conn: c,
		//	data: buf,
		//}

		//拆包解包对象
		dp := NewDataPack()
		//读取客户端Msg Head 二进制流 8个字节
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			fmt.Println("read msg head error: ", err)
			break
		}
		//拆包， 得到msgID 和msgDataLen 放在msg消息中
		pack, err := dp.UnPack(headData)
		if err != nil {
			fmt.Println("unpack error: ", err)
			break
		}
		//根据dataLen 读取Data,放在msg.data中
		var data []byte
		if pack.GetMsgLen() > 0 {
			data = make([]byte, pack.GetMsgLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				fmt.Println("read msg data error")
				break
			}
			pack.SetData(data)
		}

		//得到当前conn数据的Request请求数据
		req := Request{
			conn: c,
			msg:  pack,
		}
		go func(request ziface.IRequest) {
			//从路由中，找到注册绑定的Conn对应的router调用
			c.Router.PreHandle(request)
			c.Router.Handle(request)
			c.Router.PostHandle(request)
		}(&req)

	}
}

// SendMessage 提供一个SendMsg方法 将我们要发送给客户端的数据，先进行封包，再发送
func (c *Connection) SendMessage(msgId uint32, data []byte) error {
	if c.isClosed {
		return errors.New("connection closed when send msg")
	}
	//将data 进行封包
	pack := NewDataPack()
	msg := NewMessage(msgId, data)
	//已经序列化好的二进制,即需要发送的数据
	bytes, err := pack.Pack(msg)
	if err != nil {
		fmt.Println("pack error msg id = ", msgId)
		return errors.New("pack error msg")
	}
	//将数据发送给客户端
	if _, err := c.Conn.Write(bytes); err != nil {
		fmt.Println("write msg id :", msg, "error : ", err)
		return errors.New("send to conn error")
	}
	return nil
}

// Start 启动连接 让当前的连接准备开始工作
func (c *Connection) Start() {
	fmt.Println("Conn Start()... ConnID = ", c.ConnID)
	//启动从当前连接的读数据业务
	go c.StratReader()
	//TODO 启动从当前连接的写数据的业务
}

// Stop 停止连接 结束当前连接的工作
func (c *Connection) Stop() {
	fmt.Println("Conn Stop()... ConnID = ", c.ConnID)

	//如果当前连接已经关闭
	if c.isClosed == true {
		return
	}
	c.isClosed = true
	//关闭socket连接
	c.Conn.Close()
	//回收资源
	close(c.EXitChan)
}

// GetTCPConnection 获取当前连接绑定的socket connection
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

// GetConnID 获取当前连接模块的连接ID
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

// RemoteAddr 获取远程客户端的 TCP状态 IP port
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}
