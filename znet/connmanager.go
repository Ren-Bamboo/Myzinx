package znet

import (
	"errors"
	"fmt"
	"sync"

	"github.com/Ren-Bamboo/Myzinx/ziface"
)

type ConnManager struct {
	// connection 连接集
	connections map[uint32]ziface.IConnection

	// 与connection通信的通道true表示加入、false表示退出
	toConnChan chan map[bool]ziface.IConnection

	// 读写锁
	rwLock sync.RWMutex
}

// 监控Connection发送的信号，进行加入或移除“连接集”
func (cm *ConnManager) ListenSignal() {
	// 循环听取信号
	for signal := range cm.toConnChan {
		if v, ok := signal[true]; ok {
			// 加入CM
			fmt.Printf("Connection-%d，加入ConnManager\n", v.GetConnID())
			cm.Add(v)
		} else if v, ok = signal[false]; ok {
			// 退出CM
			fmt.Printf("Connection-%d，退出ConnManager\n", v.GetConnID())
			cm.Remove(v)
		} else {
			fmt.Printf("error in ListenSignal() ")
			panic("错误发生")
		}
		fmt.Println("当前维持连接的对象：", cm.connections)
	}
}

func NewConnManager() ziface.IConnManager {
	cm := &ConnManager{
		connections: make(map[uint32]ziface.IConnection),
		toConnChan:  make(chan map[bool]ziface.IConnection), // 不能要有缓存的，万一缓冲中有移除连接信号（Connection调用了stop），而又请求了查询，就会出问题。
	}
	// 进行监控
	go cm.ListenSignal()

	return cm
}

func (cm *ConnManager) GetTCC() chan map[bool]ziface.IConnection {
	return cm.toConnChan
}

// 添加连接
func (cm *ConnManager) Add(connection ziface.IConnection) {
	// 使用写锁
	cm.rwLock.Lock()
	defer cm.rwLock.Unlock()

	if _, ok := cm.connections[connection.GetConnID()]; !ok {
		cm.connections[connection.GetConnID()] = connection
	} else {
		fmt.Println("error in Add(connection ziface.IConnection)连接已经存在")
	}
}

// 移除连接，主要用于客户端主动断开后，删除记录连接
func (cm *ConnManager) Remove(connection ziface.IConnection) {
	// 使用写锁
	cm.rwLock.Lock()
	defer cm.rwLock.Unlock()

	if _, ok := cm.connections[connection.GetConnID()]; ok {
		delete(cm.connections, connection.GetConnID())
	} else {
		fmt.Println("error in Remove(connection ziface.IConnection) 连接不存在")
	}
}

// 根据ID下线Connection
func (cm *ConnManager) Offline(connID uint32) bool {
	// 使用写锁
	cm.rwLock.Lock()
	defer cm.rwLock.Unlock()
	if conn, ok := cm.connections[connID]; ok {
		// 断开连接，Stop会自动发送信号，清除连接记录
		conn.Stop()
		return true
	}
	fmt.Println("error in Offline(connection ziface.IConnection) 连接不存在")
	return false
}

// 根据ID查询连接
func (cm *ConnManager) Get(connID uint32) (ziface.IConnection, error) {
	// 使用读锁
	cm.rwLock.RLock()
	defer cm.rwLock.RUnlock()

	if conn, ok := cm.connections[connID]; ok {
		return conn, nil
	} else {
		fmt.Println("error in Remove(connection ziface.IConnection) 连接不存在")
		return nil, errors.New("没有对应Connection")
	}
}

// 获取连接个数
func (cm *ConnManager) Count() uint32 {
	// 使用读锁
	cm.rwLock.RLock()
	defer cm.rwLock.RUnlock()

	return uint32(len(cm.connections))
}

// 清理所有连接，服务端主动断开
func (cm *ConnManager) Clear() {
	// 使用写锁
	cm.rwLock.Lock()
	defer cm.rwLock.Unlock()

	//  断开记录的Connection
	for connID, conn := range cm.connections {
		conn.Stop()
		delete(cm.connections, connID)
	}
	// 清除自身申请的资源
	close(cm.toConnChan)
}
