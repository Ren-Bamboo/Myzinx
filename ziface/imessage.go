package ziface

type IMessage interface {
	// 设置消息ID
	SetID(uint32)
	// 设置消息长度
	SetLen(uint32)
	// 设置消息数据
	SetData([]byte)

	// 获取消息ID
	GetID() uint32
	// 获取消息长度
	GetLen() uint32
	// 获取消息数据
	GetData() []byte
}
