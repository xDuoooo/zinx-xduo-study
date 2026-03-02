package main

import (
	"fmt"
	"io"
	"net"
	"time"
	"zinx-xduo-study/src/znet"
)

/*
模拟客户端
*/
func main() {
	fmt.Println("client start...")
	time.Sleep(1 * time.Second)

	//1. 直接连接远程服务器，得到一个connect连接

	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("client start err,exit!")
		return
	}

	for {
		//封包消息

		dp := znet.NewDataPack()
		pack, err := dp.Pack(znet.NewMessage(1, []byte("ZinxV0.6 client Test Message")))
		if err != nil {
			fmt.Println("Pack error:", err)
			break
		}
		_, err = conn.Write(pack)
		if err != nil {
			fmt.Println("read buf error")
			break
		}
		//读取服务器回复的message
		head := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(conn, head); err != nil {
			fmt.Println("read head error: ", err)
			break
		}
		messageHead, err := dp.UnPack(head)
		if err != nil {
			fmt.Println("client unpack message error: ", err)
			break
		}
		if messageHead.GetMsgLen() > 0 {
			msg := messageHead.(*znet.Message)
			msg.Data = make([]byte, messageHead.GetMsgLen())
			_, err := io.ReadFull(conn, msg.Data)
			if err != nil {
				fmt.Println("client read message error: ", err)
				break
			}
			fmt.Println(msg.Id, msg.DataLen, string(msg.Data))

		}

		time.Sleep(1 * time.Second)
	}

}
