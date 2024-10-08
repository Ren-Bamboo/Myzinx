package main

import (
	"fmt"
	"io"
	"net"
	"time"

	"github.com/Ren-Bamboo/Myzinx/znet"
)

func recv(conn net.Conn) {
	dp := znet.DataPack{}
	for {
		// 接收消息头
		bufferHead := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(conn, bufferHead); err != nil {
			fmt.Println("error in io.ReadFull(conn, bufferHead)", err)
			break
		}
		msg, err := dp.Unpack(bufferHead)
		if err != nil {
			fmt.Println("error in dp.Unpack(bufferHead)", err)
			break
		}
		// 读取消息体
		if msg.GetLen() > 0 {
			bufferData := make([]byte, msg.GetLen())
			if _, err := io.ReadFull(conn, bufferData); err != nil {
				fmt.Println("error in ReadFull(conn, bufferData)", err)
				break
			}
			if _, err := dp.UnpackData(bufferData, msg); err != nil {
				fmt.Println("error in dp.UnpackData(bufferData, msg", err)
				break
			}
		}
		fmt.Printf("Receive MSG\nID: %d, Len: %d, Data: %s\n", msg.GetID(), msg.GetLen(), msg.GetData())
	}
}
func send(conn net.Conn) {
	dp := znet.DataPack{}
	for i := 1; i <= 2; i++ {
		// 封装消息
		buffer := []byte("hello word ")
		msg := znet.NewMessage(uint32(i), buffer)
		// 发送消息
		if msgByte, err := dp.Pack(msg); err != nil {
			fmt.Println("error in dp.Pack(msg)", err)
			break
		} else {
			if _, err := conn.Write(msgByte); err != nil {
				fmt.Println("error in .Write(msgByte)", err)
				break
			}
		}
		// 执行时间阻塞
		time.Sleep(2 * time.Second)
	}
}

// 模拟客户端
func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8889")
	if err != nil {
		fmt.Println("Error in net.Dial", err)
		return
	}
	// 开启接收业务
	go recv(conn)

	// 处理业务
	go send(conn)

	// 阻塞
	select {}
	defer conn.Close()
}
