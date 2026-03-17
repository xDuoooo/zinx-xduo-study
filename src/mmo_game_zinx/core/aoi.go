package core

import "fmt"

// AOIManager AOI区域管理模块
type AOIManager struct {
	//区域左边界坐标
	MinX int
	//区域右边界坐标
	MaxX int
	//X方向格子数量
	CntsX int
	//区域上边界坐标
	MinY int
	//区域下边界坐标
	MaxY int
	//Y方向格子数量
	CntsY int
	//当前区域有哪些格子
	Grids map[int]*Grid
}

// NewAOIManager 初始化一个AOI区域管理模块
func NewAOIManager(minX, maxX, cntsX, minY, maxY, cntsY int) *AOIManager {
	aoi := &AOIManager{
		MinX:  minX,
		MaxX:  maxX,
		CntsX: cntsX,
		MinY:  minY,
		MaxY:  maxY,
		CntsY: cntsY,
		Grids: make(map[int]*Grid),
	}
	//给AOI初始化区域的格子中所有的格子进行编号 和 初始化工作
	for y := 0; y < cntsY; y++ {
		for x := 0; x < cntsX; x++ {
			gid := y*cntsX + x
			//初始化gid格子
			aoi.Grids[gid] = NewGrid(gid, aoi.MinX+x*aoi.gridWidth(),
				aoi.MinX+(x+1)*aoi.gridWidth(),
				aoi.MinY+y*aoi.gridHeight(),
				aoi.MinY+(y+1)*aoi.gridHeight())
		}
	}

	return aoi
}

// 得到每个格子在X轴方向的宽度
func (m *AOIManager) gridWidth() int {
	return (m.MaxX - m.MinX) / m.CntsX
}

// 得到每个格子在Y轴方向的宽度
func (m *AOIManager) gridHeight() int {
	return (m.MaxY - m.MinY) / m.CntsY
}

// 打印格子信息
func (m *AOIManager) String() string {
	s := fmt.Sprintf("AOIManager: \n MinX: %d, MaxX: %d, cntsX: %d,minY:%d,maxY:%d cntsY:%d \n Grids in AOIManager:\n", m.MinX, m.MaxX, m.CntsX, m.MinY, m.MaxY, m.CntsY)
	for _, grid := range m.Grids {
		s += fmt.Sprintln(grid)
	}
	return s
}

// 根据GID得到九宫格周边的GID
func (m *AOIManager) GetSurroundGridByGid(gID int) (Grids []*Grid) {
	//判断是否在格子中
	if _, ok := m.Grids[gID]; !ok {
		return
	}
	Grids = append(Grids, m.Grids[gID])
	//GID的X坐标
	idX := gID % m.CntsX
	//看看左边是否有格子
	if idX-1 >= 0 {
		//尝试把放在返回之中
		Grids = append(Grids, m.Grids[gID-1])
	}
	//看看右边是否有格子
	if idX+1 <= m.CntsX-1 {
		Grids = append(Grids, m.Grids[gID+1])
	}
	gridX := make([]int, 0, len(Grids))
	for _, value := range Grids {
		gridX = append(gridX, value.GID)
	}
	for _, value := range gridX {
		idY := value / m.CntsX
		//看看上边是否有格子
		if idY-1 >= 0 {
			//尝试把放在返回值之中
			Grids = append(Grids, m.Grids[value-m.CntsX])
		}
		//看看右边是否有格子
		if idY+1 <= m.CntsY-1 {
			//尝试把放在返回值之中
			Grids = append(Grids, m.Grids[value+m.CntsX])
		}
	}
	return Grids
}

//通过x，y坐标得到对应的格子

func (m *AOIManager) GetGidByPos(x, y float32) int {
	idx := (int(x) - m.MinX) / m.gridWidth()
	idy := (int(y) - m.MinY) / m.gridHeight()
	return idy*m.CntsX + idx
}

// 通过x,y坐标得到周边九宫格全部player的ID
func (m *AOIManager) GetSurroundPlayerIDsByPos(x, y float32) (playerIDs []int) {
	//得到当前玩家格子ID
	gid := m.GetGidByPos(x, y)
	//通过GID得到周边九宫格信息
	Grids := m.GetSurroundGridByGid(gid)
	//将九宫格的信息的ID全部放在result中
	for _, v := range Grids {
		playerIDs = append(playerIDs, v.GetPlayerIDs()...)
	}
	return playerIDs
}

// AddPidToGrid 添加一个PlayerID到一个格子中
func (m *AOIManager) AddPidToGrid(pid int, gid int) {
	m.Grids[gid].Add(pid)
}

// RemovePidFromGrid 移除一个格子的playerID
func (m *AOIManager) RemovePidFromGrid(pid int, gid int) {
	m.Grids[gid].Remove(pid)
}

// 通过GID获取全部的PlayerID
func (m *AOIManager) GetPidsToGid(gid int) []int {
	return m.Grids[gid].GetPlayerIDs()
}

// 通过坐标把一个player添加到一个格子中
func (m *AOIManager) AddPidToGridByPos(pID int, x, y float32) {
	gID := m.GetGidByPos(x, y)
	m.Grids[gID].Add(pID)
}

// 通过坐标把一个player添加到一个格子中
func (m *AOIManager) RemovePidFromGridByPos(pID int, x, y float32) {
	gID := m.GetGidByPos(x, y)
	m.Grids[gID].Remove(pID)
}

//定义一些AOI的边界值

const (
	AOI_MIN_X  int = 85
	AOI_MAX_X  int = 410
	AOI_CNTS_X int = 10
	AOI_MIN_Y  int = 75
	AOI_MAX_Y  int = 400
	AOI_CNTS_Y int = 20
)
