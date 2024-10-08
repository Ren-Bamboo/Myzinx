package main

import (
	"fmt"

	"github.com/Ren-Bamboo/Myzinx/ziface"
	"github.com/Ren-Bamboo/Myzinx/znet"
)

// 基于1.1开发
type PingRouter struct {
	znet.BaseRouter
}

func (pr *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("来自消息：", request.GetData())

	buffer := []byte(fmt.Sprintf("Handle Msg Type: %d", request.GetMsgID()))
	buffer = append(buffer, request.GetData()...)
	if err := request.GetConn().Send(201, buffer); err != nil {
		fmt.Println("error in request.GetConn().Send(1,buffer)", err)
		return
	}
}

type RecallRouter struct {
	znet.BaseRouter
}

func (pr *RecallRouter) Handle(request ziface.IRequest) {
	fmt.Println("来自消息：", request.GetData())

	buffer := []byte(fmt.Sprintf("Handle Msg Type: %d", request.GetMsgID()))
	buffer = append(buffer, request.GetData()...)
	if err := request.GetConn().Send(202, buffer); err != nil {
		fmt.Println("error in request.GetConn().Send(1,buffer)", err)
		return
	}
}

// 创建HookStart
func hookStart(conn ziface.IConnection) {
	fmt.Println("Hook Start")
	conn.Send(10, []byte("hookStart********"))

	//	设置Connection属性
	conn.SetProperty("name", "bamboo")
	conn.SetProperty("github", "bamboo@github.com")
	conn.SetProperty("phone", "1804446++++")
}

// 创建HookStop
func hookStop(conn ziface.IConnection) {
	fmt.Println("Hook Stop")

	for _, k := range []string{"name", "github", "phone"} {
		if v, err := conn.GetProperty(k); err == nil {
			fmt.Println("name:", v)
		} else {
			fmt.Println("error", err)
		}
	}
}

func main() {
	server := znet.NewServer("v1.1")

	// 自定义Router
	server.AddRouter(1, &PingRouter{})
	server.AddRouter(2, &RecallRouter{})

	// 设置hook函数
	server.SetHookStart(hookStart)
	server.SetHookStop(hookStop)

	server.Server()
}
