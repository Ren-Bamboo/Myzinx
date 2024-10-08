package znet

import "github.com/Ren-Bamboo/Myzinx/ziface"

type Message struct {
	ID   uint32 // 消息ID
	Len  uint32 // 消息长度
	Data []byte // 消息数据
}

func NewMessage(id uint32, data []byte) ziface.IMessage {
	return &Message{
		ID:   id,
		Len:  uint32(len(data)),
		Data: data,
	}
}

// 设置消息ID
func (m *Message) SetID(id uint32) {
	m.ID = id
}

// 设置消息长度
func (m *Message) SetLen(len uint32) {
	m.Len = len
}

// 设置消息数据
func (m *Message) SetData(data []byte) {
	m.Data = data
}

// 获取消息ID
func (m *Message) GetID() uint32 {
	return m.ID
}

// 获取消息长度
func (m *Message) GetLen() uint32 {
	return m.Len
}

// 获取消息数据
func (m *Message) GetData() []byte {
	return m.Data
}
