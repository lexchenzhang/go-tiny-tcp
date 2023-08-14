package znet

import (
	"fmt"
	"strconv"

	"github.com/lexchenzhang/go-tiny-tcp/utils"
)

type IMsgHandler interface {
	DoMsgHandler(req IRequest)
	AddRouter(msgId uint32, router IRouter)
	StartOneWorker(workerID int, taskQueue chan IRequest)
	StartWorkerPool()
	SendRequestToTaskQueue(req IRequest)
}

type MsgHandler struct {
	routers map[uint32]IRouter
	// MQ
	taskQueue []chan IRequest
	// worker number
	workerPoolSize uint32
}

func NewMsgHandler() *MsgHandler {
	return &MsgHandler{
		routers:        make(map[uint32]IRouter),
		workerPoolSize: utils.GlobalObject.WorkerPoolSize,
		taskQueue:      make([]chan IRequest, utils.GlobalObject.WorkerPoolSize),
	}
}

func (mh *MsgHandler) AddRouter(msgId uint32, router IRouter) {
	if _, ok := mh.routers[msgId]; ok {
		panic("repeat add router msgID[" + strconv.Itoa(int(msgId)) + "]")
	}
	mh.routers[msgId] = router
}

func (mh *MsgHandler) DoMsgHandler(req IRequest) {
	if router, ok := mh.routers[req.GetMsgID()]; ok {
		router.PreHandleFunc(req)
		router.HandleFunc(req)
		router.PostHandleFunc(req)
	}
}

func (mh *MsgHandler) StartOneWorker(workerID int, taskQueue chan IRequest) {
	fmt.Println("worker start, id=", workerID)
	for {
		select {
		case req := <-taskQueue:
			mh.DoMsgHandler(req)
		}
	}
}

func (mh *MsgHandler) StartWorkerPool() {
	for i := 0; i < int(mh.workerPoolSize); i++ {
		mh.taskQueue[i] = make(chan IRequest, utils.GlobalObject.MaxWorkerTaskLen)
		go mh.StartOneWorker(i, mh.taskQueue[i])
	}
}

func (mh *MsgHandler) SendRequestToTaskQueue(req IRequest) {
	workerID := req.GetConn().GetConnID() % mh.workerPoolSize
	mh.taskQueue[workerID] <- req
}
