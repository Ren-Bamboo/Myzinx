package ziface

type IMsgHandler interface {
	//添加消息处理Handler的映射
	AddHandler(msgID uint32, router IRouter)
	//处理消息的Handler
	DoHandler(request IRequest)
	//开启消息池资源分配，执行一次
	StartWorkPool()
	//将request发送至工作池中
	SendToWorkPool(request IRequest)
}
