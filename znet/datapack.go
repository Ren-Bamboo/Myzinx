package znet

import (
	"bytes"
	"encoding/binary"
	"errors"

	"github.com/Ren-Bamboo/Myzinx/utils"
	"github.com/Ren-Bamboo/Myzinx/ziface"
)

type DataPack struct {
}

// 创建DataPack
func NewDataPack() *DataPack {
	return &DataPack{}
}

// 获取包头长度
func (dp *DataPack) GetHeadLen() uint32 {
	// Message Len: 4Byte + Message ID: 4Byte
	return 8
}

// 封包：将Message按照TLV格式封装，转为[]byte传输
func (dp *DataPack) Pack(msg ziface.IMessage) ([]byte, error) {
	dataBuffer := bytes.NewBuffer([]byte{})
	//		将Len、ID、Data按格式写入Buffer中
	if err := binary.Write(dataBuffer, binary.LittleEndian, msg.GetLen()); err != nil {
		return nil, errors.New("error in binary.Write 1")
	}
	if err := binary.Write(dataBuffer, binary.LittleEndian, msg.GetID()); err != nil {
		return nil, errors.New("error in binary.Write 2")
	}
	if err := binary.Write(dataBuffer, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, errors.New("error in binary.Write 3")
	}
	return dataBuffer.Bytes(), nil
}

// 解包：这里只解出Message的长度和ID
func (dp *DataPack) Unpack(binaryData []byte) (ziface.IMessage, error) {
	dataBuffer := bytes.NewReader(binaryData)
	msg := &Message{}

	if err := binary.Read(dataBuffer, binary.LittleEndian, &msg.Len); err != nil {
		return nil, err
	}
	if err := binary.Read(dataBuffer, binary.LittleEndian, &msg.ID); err != nil {
		return nil, err
	}
	// 判断长度是否超过预定
	if utils.GlobalObject.MaxPackageSize > 0 && msg.Len > utils.GlobalObject.MaxPackageSize {
		return nil, errors.New("error in func (dp *DataPack) Unpack(binaryData []byte) 超过最大预定包长度")
	}
	return msg, nil
}

// 解包2：解出数据，并拼接
func (dp *DataPack) UnpackData(binaryData []byte, msg ziface.IMessage) (ziface.IMessage, error) {
	dataBuffer := make([]byte, msg.GetLen())
	if err := binary.Read(bytes.NewReader(binaryData), binary.LittleEndian, dataBuffer); err != nil {
		return nil, errors.New("error in binary.Read(bytes.NewReader(binaryData), binary.LittleEndian, dataBuffer)")
	}
	msg.SetData(dataBuffer)
	return msg, nil
}
