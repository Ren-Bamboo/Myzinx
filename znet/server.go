package znet

import (
	"fmt"
	"net"

	"github.com/Ren-Bamboo/Myzinx/utils"
	"github.com/Ren-Bamboo/Myzinx/ziface"
)

type Server struct {
	//	服务器属性
	//	服务器名称
	Name string
	//	服务器绑定IP版本
	IPVersion string
	//	服务器监听的IP
	IP string
	//	服务器监听的端口
	Port int
	//	管理消息的MsgHandler
	MsgHandler ziface.IMsgHandler
	//	管理连接的ConnManager
	ConnManager ziface.IConnManager

	// 对外开放的钩子函数：在ConnStart后
	HookConnStart func(connection ziface.IConnection)
	// 对外开放的钩子函数：在ConnStop前
	HookConnStop func(connection ziface.IConnection)
}

// 创建Server，返回IServer接口指针
func NewServer(name string) ziface.IServer {
	return &Server{
		Name:          utils.GlobalObject.Name,
		IP:            utils.GlobalObject.Host,
		Port:          utils.GlobalObject.Port,
		MsgHandler:    NewMsgHandler(),
		ConnManager:   NewConnManager(),
		HookConnStart: nil,
		HookConnStop:  nil,
	}
}

// Start 实现接口
// 启动服务器
func (s *Server) Start() {
	fmt.Println("Server Start")
	fmt.Println("当前服务端配置")
	utils.GlobalObject.ShowConfig()

	// 开启工作池模式
	s.MsgHandler.StartWorkPool()

	//	绑定、监听
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.IP, s.Port))
	if err != nil {
		fmt.Println("ERROR in net.Listen ", err)
		return
	}
	//	循环接收客户端连接
	connid := uint32(0)
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("ERROR in listener.Accept()")
			continue // 退出这个连接
		}

		// 连接管理
		if s.ConnManager.Count() >= utils.GlobalObject.MaxConn {
			//TODO 返回一个最大连接数
			// 超过最大连接数
			fmt.Println("超过最大连接数量")
			conn.Close()
			continue
		}

		//	处理当前连接的业务类
		delConn := NewConnection(s.HookConnStart, s.HookConnStop, s.ConnManager.GetTCC(), conn, connid, s.MsgHandler)

		connid += 1
		// 启动连接处理
		go delConn.Start()
	}
}

// Stop 停止服务器
func (s *Server) Stop() {
	fmt.Println("[Server Stop]")

	// 连接清理
	s.ConnManager.Clear()
}

// Server 运行服务器
func (s *Server) Server() {
	// 启动Server
	go s.Start()
	//	其他的业务操作
	// TODO
	//阻塞等待
	select {}
}

func (s *Server) AddRouter(msgID uint32, router ziface.IRouter) {
	s.MsgHandler.AddHandler(msgID, router)
}

func (s *Server) GetConnManager() ziface.IConnManager {
	return s.ConnManager
}

func (s *Server) SetHookStart(hook func(connection ziface.IConnection)) {
	fmt.Println("设置HookConnStart")
	s.HookConnStart = hook
}
func (s *Server) SetHookStop(hook func(connection ziface.IConnection)) {
	fmt.Println("设置HookConnStop")
	s.HookConnStop = hook
}

func (s *Server) CallHookStart(connection ziface.IConnection) {
	if s.HookConnStart != nil {
		s.HookConnStart(connection)
	}
}
func (s *Server) CallHookStop(connection ziface.IConnection) {
	if s.HookConnStop != nil {
		s.HookConnStop(connection)
	}
}
