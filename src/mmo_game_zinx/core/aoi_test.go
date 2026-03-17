package core

import (
	"fmt"
	"testing"
)

func TestNewAOIManager(t *testing.T) {
	//初始化AOIManager
	aoiMgr := NewAOIManager(0, 250, 5, 0, 250, 5)

	//打印AOIManager
	fmt.Println(aoiMgr)
}
func TestAOIManager_GetSurroundGridByGid(t *testing.T) {
	//初始化AOIManager
	aoiMgr := NewAOIManager(0, 250, 5, 0, 250, 5)

	for gid, _ := range aoiMgr.Grids {
		Grids := aoiMgr.GetSurroundGridByGid(gid)
		gids := make([]int, 0, len(Grids))
		for _, value := range Grids {
			gids = append(gids, value.GID)
		}
		fmt.Printf("gid: %d 周围的格子ID: %v\n", gid, gids)
	}
}
