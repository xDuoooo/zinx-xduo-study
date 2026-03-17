package main

import (
	"fmt"
	"zinx-xduo-study/src/mmo_game_zinx/api"
	"zinx-xduo-study/src/mmo_game_zinx/core"
	"zinx-xduo-study/src/ziface"
	"zinx-xduo-study/src/znet"
)

// 服务器启动入口
func OnConnectionAdd(conn ziface.IConnection) {
	player := core.NewPlayer(conn)
	//给客户端发送MsgId：1的消息
	player.SyncPid()
	//给客户端发送MsgId：200消息
	player.BroadCastStartPosition()
	//将当前新上线的玩家添加到World中
	core.WorldMgrObj.AddPlayer(player)
	//将该连接绑定一个PID 以便于以后的接口使用
	conn.SetProperty("pid", player.Pid)
	//同步周边玩家，告知他们当前玩家已经上线，广播当前玩家的位置信息
	player.SyncSurrounding()
	fmt.Printf("player id:%d is online\n", player.Pid)
}
func OnConnectionLost(conn ziface.IConnection) {
	// Get the "pID" property of the current connection
	// 获取当前连接的PID属性
	pID, _ := conn.GetProperty("pID")
	var playerID int32
	if pID != nil {
		playerID = pID.(int32)
	}

	// Get the corresponding player object based on the player ID
	// 根据pID获取对应的玩家对象
	player := core.WorldMgrObj.GetPlayerByPID(playerID)

	// Trigger the player's disconnection business logic
	// 触发玩家下线业务
	if player != nil {
		player.LostConnection()
	}

	fmt.Println("====> Player ", playerID, " left =====")

}
func main() {
	s := znet.NewServer("MMO Game Zinx")

	//连接创建和销毁Hook函数
	s.SetBeforeConnCreateFunc(OnConnectionAdd)
	s.SetAfterConnDeployFunc(OnConnectionLost)
	//注册一些路由业务
	s.AddRouter(2, &api.WorldChatAPi{})
	s.AddRouter(3, &api.MoveApi{})
	//启动服务
	s.Server()
}
