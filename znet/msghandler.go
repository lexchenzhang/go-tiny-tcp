package znet

import "strconv"

type IMsgHandler interface {
	DoMsgHandler(req IRequest)
	AddRouter(msgId uint32, router IRouter)
}

type MsgHandler struct {
	routers map[uint32]IRouter
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

func NewMsgHandler() *MsgHandler {
	return &MsgHandler{
		routers: make(map[uint32]IRouter),
	}
}
