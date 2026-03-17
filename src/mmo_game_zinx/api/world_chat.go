package api

import (
	"fmt"
	"zinx-xduo-study/src/mmo_game_zinx/core"
	"zinx-xduo-study/src/mmo_game_zinx/pb"
	"zinx-xduo-study/src/ziface"
	"zinx-xduo-study/src/znet"

	"google.golang.org/protobuf/proto"
)

//世界聊天 路由业务

type WorldChatAPi struct {
	znet.BaseRouter
}

func (wc *WorldChatAPi) Handle(request ziface.IRequest) {
	pbMsg := &pb.Talk{}
	//解析数据
	err := proto.Unmarshal(request.GetData(), pbMsg)
	if err != nil {
		fmt.Println("Talk Unmarshal error ", err)
		return
	}

	pid, err := request.GetConnection().GetProperty("pid")

	player := core.WorldMgrObj.GetPlayerByPID(pid.(int32))

	player.Talk(pbMsg.Content)
}
