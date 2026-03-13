package znet

import (
	"errors"
	"fmt"
	"io"
	"net"
	"strings"
	"sync"
	"zinx-xduo-study/src/utils"
	"zinx-xduo-study/src/ziface"
)

// Connection 当前连接模块
type Connection struct {
	//当前Conn隶属于哪个Server
	TcpServer ziface.IServer

	//当前连接的Socket TCP套接字
	Conn *net.TCPConn

	//连接的ID
	ConnID uint32

	//当前连接状态
	isClosed bool
	// 用互斥锁保护 isClosed 读写，防止并发竞态
	closedLock sync.Mutex

	//告知当前连接已经退出的/停止 channel
	EXitChan chan bool

	// 改为有缓冲 channel（缓冲大小16），防止 SendMessage 在 Writer 未消费时永久阻塞
	msgChan chan []byte

	//该连接处理的方法Router ， key：msgID value：router
	MsgHandler ziface.IMsgHandler

	//连接属性集合
	property map[string]interface{}
	//保护连接属性集合的锁
	propertyLock sync.RWMutex
}

func (c *Connection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	delete(c.property, key)
}

func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()
	if value, ok := c.property[key]; ok {
		return value, nil
	} else {
		return nil, errors.New("no property found")
	}
}

func (c *Connection) SetProperty(key string, value interface{}) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()
	c.property[key] = value
}

// NewConnection 初始化连接模块的方法
func NewConnection(tcpServer ziface.IServer, conn *net.TCPConn, connID uint32, handler ziface.IMsgHandler) *Connection {
	c := &Connection{
		TcpServer:  tcpServer,
		Conn:       conn,
		ConnID:     connID,
		isClosed:   false,
		EXitChan:   make(chan bool, 1),
		msgChan:    make(chan []byte, 16), // 有缓冲，避免 SendMessage 阻塞
		MsgHandler: handler,
		property:   make(map[string]interface{}),
	}
	c.TcpServer.GetConnManager().Add(c)
	return c
}

// StartReader 连接的读数据业务方法
// 原方法名 StratReader 为拼写错误，已更正为 StartReader
func (c *Connection) StartReader() {
	fmt.Println("[Reader] Goroutine is running...")
	defer fmt.Println("connID = ", c.ConnID, "[Reader is exit],remove addr is ", c.RemoteAddr().String())
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
			// 区分正常断开和真正的错误
			// io.EOF：客户端正常关闭连接（发送了 FIN）
			// connection reset by peer：客户端强制关闭（如 Ctrl+C）
			// 这两种都属于正常断开，只打印 info，不视为错误
			if err == io.EOF || strings.Contains(err.Error(), "connection reset by peer") {
				fmt.Println("[Info] ConnID = ", c.ConnID, "client disconnected, addr = ", c.RemoteAddr().String())
			} else {
				fmt.Println("[Error] read msg head error: ", err)
			}
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
				if err == io.EOF || strings.Contains(err.Error(), "connection reset by peer") {
					fmt.Println("[Info] ConnID = ", c.ConnID, "client disconnected while reading data")
				} else {
					fmt.Println("[Error] read msg data error: ", err)
				}
				break
			}
			pack.SetData(data)
		}

		//得到当前conn数据的Request请求数据
		req := Request{
			conn: c,
			msg:  pack,
		}
		//go func(request ziface.IRequest) {
		//从路由中，找到注册绑定的Conn对应的router调用
		//c.MsgHandler.DoMsgHandler(request)
		//将消息发送给对应的channel
		//}(&req)
		if utils.GlobalObject.WorkerPoolSize > 0 {
			//已经开启了工作池机制，将消息 发送给Worker工作池处理即可
			c.MsgHandler.SendMsgToTaskQueue(&req)
		} else {
			go c.MsgHandler.DoMsgHandler(&req)
		}

	}
}

// StartWriter 写消息Goroutine,专门发送给客户端消息的模块
func (c *Connection) StartWriter() {
	fmt.Println("[Writer Goroutine is running]")
	defer fmt.Println(c.RemoteAddr().String(), "[connection Writer exit!]")

	for {
		select {
		case data := <-c.msgChan:
			//有数据要写给客户端
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("Send data error, ", err)
			}
		case <-c.EXitChan:
			fmt.Println("Writer Stop()... ConnID = ", c.ConnID)
			//代表Reader已经退出，该协程也要退出
			return
		}
	}
}

// SendMessage 提供一个SendMsg方法 将我们要发送给客户端的数据，先进行封包，再发送
func (c *Connection) SendMessage(msgId uint32, data []byte) error {
	// 用锁安全读取 isClosed
	c.closedLock.Lock()
	isClosed := c.isClosed
	c.closedLock.Unlock()
	if isClosed {
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
	// 使用 select + default 防止 msgChan 满时永久阻塞
	select {
	case c.msgChan <- bytes:
	default:
		return errors.New("send msg channel is full, drop message")
	}
	return nil
}

// Start 启动连接 让当前的连接准备开始工作
func (c *Connection) Start() {
	fmt.Println("Conn Start()... ConnID = ", c.ConnID)
	c.TcpServer.CallBeforeConnCreateFunc(c)
	//启动从当前连接的读数据业务
	go c.StartReader() // 同步更新调用名
	go c.StartWriter()
}

// Stop 停止连接 结束当前连接的工作
func (c *Connection) Stop() {
	fmt.Println("Conn Stop()... ConnID = ", c.ConnID)

	// 用互斥锁保护 isClosed，防止多个 goroutine 并发调用 Stop() 产生竞态
	c.closedLock.Lock()
	if c.isClosed {
		c.closedLock.Unlock()
		return
	}
	c.isClosed = true
	c.closedLock.Unlock()

	//调用开发者注册的 销毁连接之前 需要执行的业务Hook函数
	c.TcpServer.CallAfterConnDeployFunc(c)
	//关闭socket连接
	c.Conn.Close()
	//将当前连接从ConnManager中删除
	c.TcpServer.GetConnManager().Remove(c)
	//告知Writer关闭（先 Remove 再通知，避免 Writer 还未退出时 channel 被关闭）
	c.EXitChan <- true
	//回收资源
	close(c.EXitChan)
	close(c.msgChan)
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
