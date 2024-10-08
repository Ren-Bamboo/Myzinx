package ziface

type IRequest interface {
	//	获取当前Connection
	GetConn() IConnection
	//	获取当前数据
	GetData() []byte
	// 获取消息ID
	GetMsgID() uint32
}
