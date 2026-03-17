package core

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
	"zinx-xduo-study/src/mmo_game_zinx/pb"
	"zinx-xduo-study/src/ziface"

	"google.golang.org/protobuf/proto"
)

// Player 玩家对象
type Player struct {
	Pid  int32
	Conn ziface.IConnection //当前玩家的连接，用于和客户端交互
	X    float32            //平面X
	Y    float32            //高度
	Z    float32            //平面y坐标 (注意不是Y)
	V    float32            //旋转的角度(0-360)
}

var PidGen int32 = 1
var IdLock sync.Mutex

// NewPlayer 创建一个玩家的方法
func NewPlayer(conn ziface.IConnection) *Player {
	IdLock.Lock()
	id := PidGen
	PidGen++
	IdLock.Unlock()

	p := &Player{
		Pid:  id,
		Conn: conn,
		X:    float32(160 + rand.Intn(10)), //随机在160坐标点 基于X轴若干偏移
		Y:    0,
		Z:    float32(140 + rand.Intn(20)),
		V:    0,
	}

	return p

}

// SendMsg 提供一个发送给客户端消息的方法，主要是将pb的protobuf数据序列化之后，再调用zinx的SendMsg方法
func (p *Player) SendMsg(msgId uint32, data proto.Message) {
	//将data数据 Message结构体序列化，转化为二进制
	msg, err := proto.Marshal(data)
	if err != nil {
		fmt.Println("marshal msg err:", err)
		return
	}
	if p.Conn == nil {
		fmt.Println("connection in player is nil")
		return
	}
	err = p.Conn.SendMessage(msgId, msg)
	if err != nil {
		fmt.Println("Player SendMsg error!")
		return
	}
	return
}

// SyncPid 同步playerId给客户端
func (p *Player) SyncPid() {
	data := &pb.SyncPid{
		Pid: p.Pid,
	}
	p.SendMsg(1, data)
}

// 同步player位置给客户端
func (p *Player) BroadCastStartPosition() {
	data := &pb.BroadCast{
		Pid: p.Pid,
		Tp:  2,
		Data: &pb.BroadCast_P{
			P: &pb.Position{
				X: p.X,
				Y: p.Y,
				Z: p.Z,
				V: p.V,
			},
		},
	}
	p.SendMsg(200, data)
}

// Talk
func (p *Player) Talk(content string) {
	protoMsg := &pb.BroadCast{
		Pid:  p.Pid,
		Tp:   1,
		Data: &pb.BroadCast_Content{Content: content},
	}
	players := WorldMgrObj.GetAllPlayers()

	for _, player := range players {
		player.SendMsg(200, protoMsg)
	}

}

func (p *Player) SyncSurrounding() {
	//获取当前玩家周围的玩家有哪些
	pids := WorldMgrObj.AOIManager.GetSurroundPlayerIDsByPos(p.X, p.Z)
	players := make([]*Player, 0, len(pids))
	for _, pid := range pids {
		players = append(players, WorldMgrObj.GetPlayerByPID(int32(pid)))
	}
	//告知其他玩家自己位置信息
	protoMsg := &pb.BroadCast{
		Pid: p.Pid,
		Tp:  2,
		Data: &pb.BroadCast_P{
			P: &pb.Position{
				X: p.X,
				Y: p.Y,
				Z: p.Z,
				V: p.V,
			},
		},
	}
	for _, player := range players {
		player.SendMsg(200, protoMsg)
	}
	playersProtoMsg := make([]*pb.Player, 0, len(players))
	for _, player := range players {
		p := &pb.Player{
			Pid: player.Pid,
			P: &pb.Position{
				X: player.X,
				Y: player.Y,
				Z: player.Z,
				V: player.V,
			},
		}
		playersProtoMsg = append(playersProtoMsg, p)
	}

	// 4.2 Encapsulate SyncPlayers protobuf data
	// 封装SyncPlayer protobuf数据
	SyncPlayersMsg := &pb.SyncPlayers{
		Ps: playersProtoMsg[:],
	}

	// 4.3 Send all player data to the current player to display surrounding players
	// 给当前玩家发送需要显示周围的全部玩家数据
	p.SendMsg(202, SyncPlayersMsg)
}

// UpdatePos Broadcast player position update
// (广播玩家位置移动)
func (p *Player) UpdatePos(x float32, y float32, z float32, v float32) {
	// 触发消失视野和添加视野业务
	// 计算旧格子gID
	oldGID := WorldMgrObj.AOIManager.GetGidByPos(p.X, p.Z)
	// Calculate the new grid gID
	// 计算新格子gID
	newGID := WorldMgrObj.AOIManager.GetGidByPos(x, z)

	// Update the player's position information
	// 更新玩家的位置信息
	p.X = x
	p.Y = y
	p.Z = z
	p.V = v

	if oldGID != newGID {
		// 触发gird切换
		// 把pID从就的aoi格子中删除
		WorldMgrObj.AOIManager.RemovePidFromGrid(int(p.Pid), oldGID)

		// 把pID添加到新的aoi格子中去
		WorldMgrObj.AOIManager.AddPidToGrid(int(p.Pid), newGID)

		_ = p.OnExchangeAoiGrID(oldGID, newGID)
	}

	// Assemble protobuf data, send position to surrounding players
	// 组装protobuf协议，发送位置给周围玩家
	msg := &pb.BroadCast{
		Pid: p.Pid,
		Tp:  4, //Tp:4  Coordinates information after movement(移动之后的坐标信息)
		Data: &pb.BroadCast_P{
			P: &pb.Position{
				X: p.X,
				Y: p.Y,
				Z: p.Z,
				V: p.V,
			},
		},
	}

	// Get all players around the current player
	// (获取当前玩家周边全部玩家)
	players := p.GetSurroundingPlayers()

	// Send MsgID:200 message to each player's client, updating position after movement
	// (向周边的每个玩家发送MsgID:200消息，移动位置更新消息)
	for _, player := range players {
		player.SendMsg(200, msg)
	}
}

func (p *Player) OnExchangeAoiGrID(oldGID, newGID int) error {
	// (获取旧的九宫格成员)
	oldGrIDs := WorldMgrObj.AOIManager.GetSurroundGridByGid(oldGID)

	// Create a hash table for the old nine-grid members to quickly search
	// 为旧的九宫格成员建立哈希表,用来快速查找
	oldGrIDsMap := make(map[int]bool, len(oldGrIDs))
	for _, grID := range oldGrIDs {
		oldGrIDsMap[grID.GID] = true
	}

	// Get members in the new nine-grid from the new grid
	// 获取新的九宫格成员
	newGrIDs := WorldMgrObj.AOIManager.GetSurroundGridByGid(oldGID)

	// Create a hash table for the new nine-grid members to quickly search
	// 为新的九宫格成员建立哈希表,用来快速查找
	newGrIDsMap := make(map[int]bool, len(newGrIDs))
	for _, grID := range newGrIDs {
		newGrIDsMap[grID.GID] = true
	}

	//------ > (处理视野消失) <-------
	offlineMsg := &pb.SyncPid{
		Pid: p.Pid,
	}

	// Find the grid IDs that appear in the old nine-grid but not in the new nine-grid
	// (找到在旧的九宫格中出现,但是在新的九宫格中没有出现的格子)
	leavingGrIDs := make([]*Grid, 0)
	for _, grID := range oldGrIDs {
		if _, ok := newGrIDsMap[grID.GID]; !ok {
			leavingGrIDs = append(leavingGrIDs, grID)
		}
	}

	// Get all players in the disappearing grids
	// (获取需要消失的格子中的全部玩家)
	for _, grID := range leavingGrIDs {
		players := WorldMgrObj.GetPlayersByGID(grID.GID)
		for _, player := range players {

			// Make oneself disappear in the views of other players
			// 让自己在其他玩家的客户端中消失
			player.SendMsg(201, offlineMsg)

			// Make other players' information disappear in one's own client
			// 将其他玩家信息 在自己的客户端中消失
			anotherOfflineMsg := &pb.SyncPid{
				Pid: player.Pid,
			}
			p.SendMsg(201, anotherOfflineMsg)
			time.Sleep(200 * time.Millisecond)
		}
	}

	// ------ > Handle visibility appearance(处理视野出现) <-------

	// Find the grid IDs that appear in the new nine-grid but not in the old nine-grid
	// 找到在新的九宫格内出现,但是没有在就的九宫格内出现的格子
	enteringGrIDs := make([]*Grid, 0)
	for _, grID := range newGrIDs {
		if _, ok := oldGrIDsMap[grID.GID]; !ok {
			enteringGrIDs = append(enteringGrIDs, grID)
		}
	}

	onlineMsg := &pb.BroadCast{
		Pid: p.Pid,
		Tp:  2,
		Data: &pb.BroadCast_P{
			P: &pb.Position{
				X: p.X,
				Y: p.Y,
				Z: p.Z,
				V: p.V,
			},
		},
	}

	// Get all players in the appearing grids
	// 获取需要显示格子的全部玩家
	for _, grID := range enteringGrIDs {
		players := WorldMgrObj.GetPlayersByGID(grID.GID)

		for _, player := range players {
			// Make oneself appear in the views of other players
			// 让自己出现在其他人视野中
			player.SendMsg(200, onlineMsg)

			// Make other players appear in one's own client
			// 让其他人出现在自己的视野中
			anotherOnlineMsg := &pb.BroadCast{
				Pid: player.Pid,
				Tp:  2,
				Data: &pb.BroadCast_P{
					P: &pb.Position{
						X: player.X,
						Y: player.Y,
						Z: player.Z,
						V: player.V,
					},
				},
			}

			time.Sleep(200 * time.Millisecond)
			p.SendMsg(200, anotherOnlineMsg)
		}
	}

	return nil
}
func (p *Player) LostConnection() {
	// 1 Get players in the surrounding AOI nine-grid
	// 获取周围AOI九宫格内的玩家
	players := p.GetSurroundingPlayers()

	// 2 Assemble MsgID:201 message
	// 封装MsgID:201消息
	msg := &pb.SyncPid{
		Pid: p.Pid,
	}

	// 3 Send messages to surrounding players
	// 向周围玩家发送消息
	for _, player := range players {
		player.SendMsg(201, msg)
	}

	// 4 Remove the current player from AOI in the world manager
	// 世界管理器将当前玩家从AOI中摘除
	WorldMgrObj.AOIManager.RemovePidFromGridByPos(int(p.Pid), p.X, p.Z)
}
func (p *Player) GetSurroundingPlayers() []*Player {
	// Get all pIDs in the current AOI area
	// 得到当前AOI区域的所有pID
	pIDs := WorldMgrObj.AOIManager.GetSurroundPlayerIDsByPos(p.X, p.Z)

	// Put all players corresponding to pIDs into the Player slice
	// 将所有pID对应的Player放到Player切片中
	players := make([]*Player, 0, len(pIDs))
	for _, pID := range pIDs {
		players = append(players, WorldMgrObj.GetPlayerByPID(int32(pID)))
	}

	return players
}
