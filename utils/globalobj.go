package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/Ren-Bamboo/Myzinx/ziface"
)

type GlobalObj struct {
	/*
		Server
	*/
	Server ziface.IServer // 当前Zinx全局的Server对象
	Host   string         // 服务器监听的IP
	Port   int            // 服务器监听的端口
	Name   string         // 服务器名称

	/*
		Myzinx
	*/
	Version        string // 服务器版本
	MaxConn        uint32 // 客户端最大连接数量
	MaxPackageSize uint32 // Myzinx框架中数据包最大长度
	WorkPoolSize   uint32 // 工作池中work数量
	MaxTaskSize    uint32 // 工作池对应的最大任务数量
}

func (g *GlobalObj) LoadUConfig() {
	file, err := os.Open("conf/zinx.json")
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("用户配置文件不存在，使用默认配置")
			return
		}
		panic(err)
	}
	data, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, GlobalObject)
	if err != nil {
		panic(err)
	}
}

func (g *GlobalObj) ShowConfig() {
	fmt.Printf("Version:%s\nHost:%s,Port:%d\nName:%s\nMaxConn:%d\nMaxPackageSize:%d\nWorkPoolSize:%d\nMaxTaskSize:%d\n", g.Version, g.Host, g.Port, g.Name, g.MaxConn, g.MaxPackageSize, g.WorkPoolSize, g.MaxTaskSize)
}

// 全局对外的GlobalObj
var GlobalObject *GlobalObj

func init() {
	//默认配置设置
	GlobalObject = &GlobalObj{
		Host: "0.0.0.0",
		Port: 8889,
		Name: "Myzinx",

		Version:        "latest:v1.1",
		MaxConn:        1000,
		MaxPackageSize: 4096,
		WorkPoolSize:   10,
		MaxTaskSize:    1024,
	}
	//	加载用户配置文件：conf/zinx.json
	GlobalObject.LoadUConfig()
}
