package ziface

type IDataPack interface {
	// 获取包头长度
	GetHeadLen() uint32
	// 封包
	Pack(msg IMessage) ([]byte, error)
	// 解包
	Unpack([]byte) (IMessage, error)
	// 解包：数据部分
	UnpackData(binaryData []byte, msg IMessage) (IMessage, error)
}
