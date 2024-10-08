package ziface

type IRouter interface {
	//	处理Request（Connection）业务前的方法
	PreHandle(IRequest)
	//	处理Request（Connection）业务的方法
	Handle(IRequest)
	//	处理Request（Connection）业务后的方法
	PastHandle(IRequest)
}
