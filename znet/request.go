package znet

import "github.com/Ren-Bamboo/Myzinx/ziface"

// 将连接和数据绑定在一起
type Request struct {
	Conn ziface.IConnection
	msg  ziface.IMessage
}

func NewRequest(conn ziface.IConnection, msg ziface.IMessage) ziface.IRequest {
	return &Request{
		Conn: conn,
		msg:  msg,
	}
}

// 获取当前Connection
func (r *Request) GetConn() ziface.IConnection {
	return r.Conn
}

// 获取消息数据
func (r *Request) GetData() []byte {
	return r.msg.GetData()
}

// 获取消息ID
func (r *Request) GetMsgID() uint32 {
	return r.msg.GetID()
}
