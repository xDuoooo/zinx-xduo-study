package core

import (
	"fmt"
	"sync"
)

// 一个AOI地图中的格子
type Grid struct {
	GID       int //格子ID
	MinX      int
	MaxX      int
	MinY      int
	MaxY      int
	playerIDs map[int]bool //格子内玩家或者物体的集合
	pIDLock   sync.RWMutex //保护当前集合的锁
}

// NewGrid 初始化方法
func NewGrid(gID, minX, maxX, minY, maxY int) *Grid {
	return &Grid{
		GID:       gID,
		MinY:      minY,
		MaxY:      maxY,
		MaxX:      maxX,
		MinX:      minX,
		playerIDs: make(map[int]bool),
	}
}

// Add 增加玩家
func (g *Grid) Add(playerID int) {
	g.pIDLock.Lock()
	defer g.pIDLock.Unlock()
	g.playerIDs[playerID] = true
}

// Remove 删除玩家
func (g *Grid) Remove(playerID int) {
	g.pIDLock.Lock()
	defer g.pIDLock.Unlock()
	delete(g.playerIDs, playerID)
}
func (g *Grid) GetPlayerIDs() (playerIDs []int) {
	g.pIDLock.RLock()
	defer g.pIDLock.RUnlock()

	for k, _ := range g.playerIDs {
		playerIDs = append(playerIDs, k)
	}
	return playerIDs
}

func (g *Grid) String() string {
	return fmt.Sprintf("Grid id: %d, minX:%d, maxX:%d, minY:%d, maxY:%d, playerIDs:%v",
		g.GID, g.MinX, g.MaxX, g.MinY, g.MaxY, g.playerIDs)
}
