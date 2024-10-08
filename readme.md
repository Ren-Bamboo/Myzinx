# 轻量级服务器框架Myzinx

学习来源于：https://www.yuque.com/aceld/npyr8s/hdktpg

区别不同：

Myzinx1.0与Zinx-v0.8前相似，不同主要在于：
	1）ConnManager与Connection之间的控制流程
	2）Connection调用Hook函数

* 针对1）：Zinx-0.8为Connection结构传入了Server，笔者参考其他意见认为此操作可能增加了耦合度。因此，Myzinx采用了channel控制Connection与ConnManager之间的交互方式。
* 针对2）：为了避免在Connection中使用Server，将Hook函数作为Connection的属性。



## 基于Myzinx框架开发的使用示例

注意：首先需要明确地对消息ID进行定义，然后再根据消息ID进行业务处理开发，**完整的示例将在项目下的“TestRun”文件夹中。**



对于客户端（主要通过框架处理包的封装和请求）：

~~~go
// 数据处理流：data---(DataPack)--->Message--->packBytes
// data 原始数据
buffer := []byte("hello word ") 
// message 封装的消息
msg := znet.NewMessage(uint32(i), buffer) 
// packBytes 打包后的Bytes
msgByte, err := dp.Pack(msg) 
~~~



对于服务端（主要编写Router和Hook函数）：

~~~go
// 数据处理流：Message---(封装)--->Request---(MsgHandler)--->Router

// 创建Router
type RecallRouter struct {
	znet.BaseRouter
}
func (pr *RecallRouter) Handle(request ziface.IRequest) {
	fmt.Println("来自消息：", request.GetData())

	buffer := []byte(fmt.Sprintf("Handle Msg Type: %d", request.GetMsgID()))
	buffer = append(buffer, request.GetData()...)
	if err := request.GetConn().Send(202, buffer); err != nil {
		fmt.Println("error in request.GetConn().Send(1,buffer)", err)
		return
	}
}
// 创建HookStart
func hookStart(conn ziface.IConnection) {
	fmt.Println("Hook Start")
	conn.Send(10, []byte("hookStart********"))

	//	设置Connection属性
	conn.SetProperty("name", "bamboo")
	conn.SetProperty("github", "bamboo@github.com")
	conn.SetProperty("phone", "1804446++++")
}
// ——————————————————————————————————————————————-
// 创建服务器
server := znet.NewServer("v1.1") 

// 将消息ID绑定Router
server.AddRouter(2, &RecallRouter{})

// 设置hook函数
server.SetHookStart(hookStart)

// 启动服务
server.Server()
~~~



配置文件：

/utils/globalobj.go：全局配置文件，包括默认的监听ip和端口等

conf/zinx.json：用户配置文件，开发者可配置相关参数覆盖原全局配置