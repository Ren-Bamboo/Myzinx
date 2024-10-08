package znet

import (
	"fmt"

	"github.com/Ren-Bamboo/Myzinx/utils"
	"github.com/Ren-Bamboo/Myzinx/ziface"
)

type MsgHandler struct {
	// 映射：消息ID对应的Router
	APIs map[uint32]ziface.IRouter
	// 工作池的work大小
	WorkPoolSize uint32
	// Work领取任务的消息队列
	TaskQueue []chan ziface.IRequest
}

func NewMsgHandler() ziface.IMsgHandler {
	return &MsgHandler{
		APIs:         make(map[uint32]ziface.IRouter),
		WorkPoolSize: utils.GlobalObject.WorkPoolSize,
		TaskQueue:    make([]chan ziface.IRequest, utils.GlobalObject.WorkPoolSize),
	}
}

func (mh *MsgHandler) StartWorkPool() {
	//	创建池资源，只执行一次
	for i := uint32(0); i < utils.GlobalObject.WorkPoolSize; i++ {
		mh.TaskQueue[i] = make(chan ziface.IRequest, utils.GlobalObject.MaxTaskSize)
		go mh.StartWorkProcedure(i)
	}
}

func (mh *MsgHandler) StartWorkProcedure(workID uint32) {
	fmt.Println("[StartWorkProcedure()]——>workID:", workID)
	// 轮询队列，处理request
	for request := range mh.TaskQueue[workID] {
		mh.DoHandler(request)
	}
}

func (mh *MsgHandler) SendToWorkPool(request ziface.IRequest) {
	// 均衡算法，计算放入那个work的队列
	toWorkID := request.GetConn().GetConnID() % utils.GlobalObject.WorkPoolSize
	fmt.Printf("ConnID-%d的request放入Work-%d的队列中", request.GetConn().GetConnID(), toWorkID)
	// 将request放入队列
	mh.TaskQueue[toWorkID] <- request
}

func (mh *MsgHandler) AddHandler(msgID uint32, router ziface.IRouter) {
	// 判断映射是否存在
	if _, ok := mh.APIs[msgID]; ok {
		fmt.Println("已经存在msgID,", msgID)
		panic("error in mh.APIs[msgID]; ok ")
	}
	// 存放映射
	mh.APIs[msgID] = router
}

func (mh *MsgHandler) DoHandler(request ziface.IRequest) {
	// 判断映射是否存在
	router, ok := mh.APIs[request.GetMsgID()]
	if !ok {
		fmt.Println("不存在msgID,", request.GetMsgID())
		panic("error in mh.APIs[request.GetMsgID()]; !ok ")
	}

	//	处理Msg
	router.PreHandle(request)
	router.Handle(request)
	router.PastHandle(request)
}
