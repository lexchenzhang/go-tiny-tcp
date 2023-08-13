package znet

type IRouter interface {
	PreHandleFunc(req IRequest)
	HandleFunc(req IRequest)
	PostHandleFunc(req IRequest)
}

type BaseRouter struct{}

func (br *BaseRouter) PreHandleFunc(req IRequest)  {}
func (br *BaseRouter) HandleFunc(req IRequest)     {}
func (br *BaseRouter) PostHandleFunc(req IRequest) {}
