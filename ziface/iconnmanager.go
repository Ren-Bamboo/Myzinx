package ziface

type IConnManager interface {
	//添加连接
	Add(IConnection)
	//删除连接
	Remove(IConnection)
	// 通过连接ID，下线对应的连接
	Offline(connID uint32) bool
	//根据ID查询连接
	Get(uint32) (IConnection, error)
	//获取连接个数
	Count() uint32
	//清理所有连接
	Clear()
	// 获得与Connetion通信的通道
	GetTCC() chan map[bool]IConnection
}
