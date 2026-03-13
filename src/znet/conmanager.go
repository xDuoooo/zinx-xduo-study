package znet

import (
	"errors"
	"fmt"
	"sync"
	"zinx-xduo-study/src/ziface"
)

type ConnManager struct {
	connections map[uint32]ziface.IConnection //已经创建的Connection集合
	connLock    sync.RWMutex                  //保护连接集合的读写锁
}

func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]ziface.IConnection),
	}
}

// Add 添加连接
func (connMgr *ConnManager) Add(connection ziface.IConnection) {
	//保护共享资源map，加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	//讲conn假如到ConnManager中
	connMgr.connections[connection.GetConnID()] = connection
	// 已持有写锁，直接用 len() 避免调用 Len() 再次加读锁导致死锁
	fmt.Println("connection add to ConnManager successfully: connection num = ", len(connMgr.connections))
}

// Remove 删除连接
func (connMgr *ConnManager) Remove(connection ziface.IConnection) {
	//保护共享资源map，加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()
	delete(connMgr.connections, connection.GetConnID())
}

// Get 根据ConnID获取连接
func (connMgr *ConnManager) Get(connID uint32) (ziface.IConnection, error) {
	//保护共享资源map，加读锁
	connMgr.connLock.RLock()
	defer connMgr.connLock.RUnlock()

	if conn, ok := connMgr.connections[connID]; ok {
		return conn, nil
	} else {
		return nil, errors.New("connection not FOUND")
	}
}

// Len 得到当前连接总数
func (connMgr *ConnManager) Len() int {
	// 加读锁，防止并发读取 map 数据竞争
	connMgr.connLock.RLock()
	defer connMgr.connLock.RUnlock()
	return len(connMgr.connections)
}

// ClearConn 清除并终止所有连接
func (connMgr *ConnManager) ClearConn() {
	// 先持锁拷贝所有连接并清空 map，然后解锁再调用 Stop()
	// 原因：Stop() 内部会调用 Remove()，Remove() 也要获取写锁，持锁期间调用会死锁
	connMgr.connLock.Lock()
	conns := make([]ziface.IConnection, 0, len(connMgr.connections))
	for _, conn := range connMgr.connections {
		conns = append(conns, conn)
	}
	count := len(connMgr.connections)                         // 在清空前记录数量，否则清空后打印永远为 0
	connMgr.connections = make(map[uint32]ziface.IConnection) // 清空 map
	connMgr.connLock.Unlock()                                 // 先解锁，再 Stop

	for _, conn := range conns {
		conn.Stop() // Remove() 在此被调用时锁已释放，不会死锁
	}
	fmt.Println("Clear All connections success! cleared conn num = ", count)
}
