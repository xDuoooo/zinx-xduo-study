package main

import (
	"fmt"
	"net"
	"time"
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
		//2. 连接调用Write写数据
		_, err := conn.Write([]byte("Hello Zinx V0.1"))

		buf := make([]byte, 512)

		cnt, err := conn.Read(buf)
		if err != nil {
			fmt.Println("read buf error")
			return
		}

		fmt.Printf(" server call back: %s, cnt = %d\n", buf[:cnt], cnt)

		time.Sleep(1 * time.Second)
	}

}
