package znet

import (
	"fmt"
	"io"
	"net"
	"testing"
)

// datapack 拆包封包的单元测试
func TestDataPack(t *testing.T) {
	/*
		模拟的服务器
	*/
	//1 创建socketTCP
	listener, err := net.Listen("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("server listener err: ", err)
		return
	}

	//2 从客户端读取数据，拆包处理
	//创建一个go去承载 负责从客户端处理业务的模块
	go func() {
		for {
			accept, err := listener.Accept()
			if err != nil {
				fmt.Println("server accept error", err)
			}
			go func(conn net.Conn) {
				//处理客户端的请求
				//拆包的过程
				//第一次 读 读head
				dp := NewDataPack()
				for {
					headData := make([]byte, dp.GetHeadLen())
					_, err := io.ReadFull(conn, headData)
					if err != nil {
						fmt.Println("read head error")
						return
					}
					//第二次读 根据head中的datalen 再读取data内容
					msgHead, err := dp.UnPack(headData)
					if err != nil {
						fmt.Println("server unpack error")
						return
					}
					if msgHead.GetMsgLen() > 0 {
						//msg 有数据的 需要进行第二次读取
						msg := msgHead.(*Message)
						msg.SetData(make([]byte, msg.GetMsgLen()))
						_, err := io.ReadFull(conn, msg.Data)
						if err != nil {
							fmt.Println("server unpack data err: ", err)
							return
						}
						//完整的一个的消息已经读取完毕
						fmt.Println("-----> Recv msgID: ", msg.Id, "dataLen: ", msg.DataLen, "msgData: ", string(msg.Data))
					}

				}
			}(accept)
		}
	}()
	/*
		模拟的客户端
	*/
	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("client dial err: ", err)
	}
	//创建一个封包对象 dp
	for {
		dp := NewDataPack()

		//模拟粘包过程，封装两个msg一同发送
		msg1 := &Message{
			Id:      1,
			DataLen: 4,
			Data:    []byte{'z', 'i', 'n', 'x'},
		}
		sendData1, err := dp.Pack(msg1)
		//封装第一个包

		//封装第二个包
		msg2 := &Message{
			Id:      1,
			DataLen: 7,
			Data:    []byte{'n', 'i', 'h', 'a', 'o', '!', '!'},
		}
		sendData2, err := dp.Pack(msg2)
		//将两个包粘在一起
		sendData1 = append(sendData1, sendData2...)
		_, err = conn.Write(sendData1)
		if err != nil {
			return
		}
	}

}

// TestStickyPacket 演示不使用封包协议时的TCP粘包现象
// 客户端发送5条独立消息，服务端不做任何拆包处理，直接读取原始字节
// 观察服务端收到的数据边界是否和发送时一致
func TestStickyPacket(t *testing.T) {
	done := make(chan struct{})

	// ========== 模拟服务端 ==========
	listener, err := net.Listen("tcp", "127.0.0.1:7778")
	if err != nil {
		fmt.Println("server listen err:", err)
		return
	}

	go func() {
		conn, err := listener.Accept()
		if err != nil {
			return
		}
		defer conn.Close()

		fmt.Println("\n========== 服务端开始读取（无任何拆包处理）==========")
		readCount := 0
		buf := make([]byte, 256) // 固定缓冲区，完全不管消息边界
		for {
			n, err := conn.Read(buf)
			if err != nil {
				break
			}
			readCount++
			// 直接打印读到的原始内容，可以看到消息边界是混乱的
			fmt.Printf("第 %d 次 Read：读到 %d 字节，内容: %q\n", readCount, n, string(buf[:n]))
		}
		fmt.Printf("========== 共调用了 %d 次 Read（客户端发送了5条消息）==========\n", readCount)
		close(done)
	}()

	// ========== 模拟客户端 ==========
	conn, err := net.Dial("tcp", "127.0.0.1:7778")
	if err != nil {
		fmt.Println("client dial err:", err)
		return
	}

	// 定义5条独立消息（客户端认为每条是单独的逻辑消息）
	messages := []string{
		"Hello",
		"World",
		"TCP is stream",
		"no boundary",
		"sticky packet!",
	}

	fmt.Println("\n========== 客户端发送5条独立消息 ==========")
	for i, msg := range messages {
		fmt.Printf("发送第 %d 条消息: %q\n", i+1, msg)
		conn.Write([]byte(msg))
	}

	conn.Close() // 关闭连接，触发服务端 Read 返回 EOF
	<-done       // 等待服务端读完所有数据

	fmt.Println("\n结论：客户端发了5次，但服务端 Read 的次数不一定是5次")
	fmt.Println("这就是粘包——TCP把多条消息合并成了一个流，消息边界丢失了")
	fmt.Println("解决：使用 DataPack 在数据前加 DataLen+MsgId 头，接收方才能准确切割消息")
}
