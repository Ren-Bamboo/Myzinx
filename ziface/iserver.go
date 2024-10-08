package ziface

// IServer 定义接口类，抽象层
type IServer interface {
	// Start 启动服务器
	Start()
	//	停止服务器
	Stop()
	//	运行服务器
	Server()
	//	添加消息ID对应的Router
	AddRouter(msgID uint32, router IRouter)
	// 返回连接管理器
	GetConnManager() IConnManager
	// 设置HookStart
	SetHookStart(hook func(connection IConnection))
	// 设置HookStop
	SetHookStop(hook func(connection IConnection))
	//调用HookStart
	CallHookStart(connection IConnection)
	//调用HookStop
	CallHookStop(connection IConnection)
}
