package ziface

import "net"

type IConnection interface {
	// 启动连接
	Start()
	// 停止连接
	Stop()

	// 设置属性
	SetProperty(string, interface{})
	// 删除属性
	RemoveProperty(string)
	// 查询属性
	GetProperty(string) (interface{}, error)

	// 获取当前连接的Conn对象
	GetConnection() net.Conn
	// 获取连接ID
	GetConnID() uint32
	// 得到客户端连接的地址和端口
	GetRoteAddr() net.Addr

	// 发送消息的方法
	Send(msgID uint32, data []byte) error
}
