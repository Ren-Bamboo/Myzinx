package znet

import (
	"errors"
	"fmt"
	"io"
	"net"
	"sync"

	"github.com/Ren-Bamboo/Myzinx/utils"
	"github.com/Ren-Bamboo/Myzinx/ziface"
)

type Connection struct {
	// 套接字Conn
	Conn net.Conn
	// 连接ID
	ConnID uint32
	// 连接状态
	isClosed bool
	//	处理消息的MsgHandler
	MsgHandler ziface.IMsgHandler

	// hook函数
	HookConnStart func(connection ziface.IConnection) // 建立连接后、执行处理业务前
	HookConnStop  func(connection ziface.IConnection) // 执行业务后、断开连接前

	// 属性字段
	Property map[string]interface{}
	// 属性操作锁
	propertyLock sync.RWMutex

	// 与ConnManager通信的通道
	ToCMChan chan map[bool]ziface.IConnection
	//	退出信号的channel
	ExitChan chan bool
	// Reader向Writer发送数据的channel（既是信号、又是数据）
	RWChan chan []byte
}

// 创建Connection的方法
func NewConnection(hookstart, hookstop func(connection ziface.IConnection), tocmchan chan map[bool]ziface.IConnection, conn net.Conn, connId uint32, msghandler ziface.IMsgHandler) *Connection {
	connection := &Connection{
		HookConnStart: hookstart,
		HookConnStop:  hookstop,
		Conn:          conn,
		ConnID:        connId,
		MsgHandler:    msghandler,
		isClosed:      false,
		Property:      make(map[string]interface{}),
		ToCMChan:      tocmchan,
		ExitChan:      make(chan bool, 1),
		RWChan:        make(chan []byte),
	}

	// 通知ConnManager，加入管理
	connection.ToCMChan <- map[bool]ziface.IConnection{true: connection}
	return connection
}

// Connection 的读取业务
func (c *Connection) StartReader() {
	fmt.Println("[StartReader() is running......]")
	defer fmt.Println("[StartReader() is exited]")
	defer c.Stop()

	for {
		// TLV格式读取消息
		//创建拆封包对象
		dp := NewDataPack()
		// 读取消息头
		bufferHead := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetConnection(), bufferHead); err != nil {
			fmt.Println("error in io.ReadFull(c.GetConnection()", err)
			break
		}
		msg, err := dp.Unpack(bufferHead)
		if err != nil {
			fmt.Println("error in dp.Unpack(bufferHead)", err)
			break
		}
		// 读取消息数据
		if msg.GetLen() > 0 {
			dataBuffer := make([]byte, msg.GetLen())
			if _, err := io.ReadFull(c.GetConnection(), dataBuffer); err != nil {
				fmt.Println("error in io.ReadFull(c.GetConnection(), dataBuffer)", err)
				break
			}
			// 数据拼接
			if _, err = dp.UnpackData(dataBuffer, msg); err != nil {
				fmt.Println("error in dp.UnpackData(dataBuffer, msg)", err)
				break
			}
		}
		// 调用Router处理业务
		// 封装Request
		request := NewRequest(c, msg)
		// request处理模式判断
		if utils.GlobalObject.WorkPoolSize > 0 {
			// 工作池模式
			c.MsgHandler.SendToWorkPool(request)
		} else {
			// 非工作池模式
			//Router处理Request业务
			go c.MsgHandler.DoHandler(request)
		}
	}
}

// 写业务
func (c *Connection) StartWriter() {
	fmt.Println("[StartWriter() is running......]")
	defer fmt.Println("[Exit StartWriter()]")
	for {
		select {
		case data := <-c.RWChan:
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("error in c.Conn.Write(data); err != nil", err)
				return
			}
		case <-c.ExitChan:
			return
		}
	}
}

// 实现IConnection接口-Start()
func (c *Connection) Start() {
	fmt.Println("[Conn start]——>ConnID:", c.ConnID)

	// 开启读业务
	go c.StartReader()

	// 开启写业务
	go c.StartWriter()

	// 调用HookStart
	c.CallHook(c.HookConnStart)
}

// 实现IConnection接口-Start()
func (c *Connection) Stop() {
	fmt.Println("[Conn stop]——>ConnID:", c.ConnID)
	defer fmt.Printf("[Connection is closed]——>closed ID: %d closed Addr: %s\n", c.ConnID, c.Conn.RemoteAddr().String())
	if c.isClosed {
		return
	}
	c.isClosed = true

	// 调用HookStop
	c.CallHook(c.HookConnStop)

	// 关闭连接
	c.Conn.Close()

	// 通知ConnManager，移除连接
	c.ToCMChan <- map[bool]ziface.IConnection{false: c}

	// 作为通知消息
	c.ExitChan <- true

	// 回收资源
	close(c.ExitChan)
	close(c.RWChan)
}

// 实现IConnection接口-GetConnection()
func (c *Connection) GetConnection() net.Conn {
	return c.Conn
}

// 实现IConnection接口-GetConnID()
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

// 实现IConnection接口-GetRoteAddr()
func (c *Connection) GetRoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

// 实现IConnection接口-Send()
func (c *Connection) Send(msgID uint32, data []byte) error {
	if c.isClosed {
		return errors.New("连接已经断开")
	}
	dp := DataPack{}
	// 创建消息
	msg := NewMessage(msgID, data)
	// pack 消息
	byteMsg, err := dp.Pack(msg)
	if err != nil {
		fmt.Println("error in dp.UnpackData(dataBuffer, msg)", err)
		return errors.New("dp.Pack(msg)错误")
	}
	// 向通道写
	c.RWChan <- byteMsg
	// 没有错误
	return nil
}

// 设置属性
func (c *Connection) SetProperty(key string, val interface{}) {
	// 添加写锁
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	c.Property[key] = val
}

// 删除属性
func (c *Connection) RemoveProperty(key string) {
	// 添加写锁
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	delete(c.Property, key)
}

// 查询属性
func (c *Connection) GetProperty(key string) (interface{}, error) {
	//	添加读锁
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()

	if v, ok := c.Property[key]; ok {
		return v, nil
	}
	return nil, errors.New("error in GetProperty(key string) ")
}

// 调用hook函数
func (c *Connection) CallHook(call func(connection ziface.IConnection)) {
	if call != nil {
		call(c)
	}
}
