package znet

type IRequest interface {
	GetConn() IConnection
	GetData() []byte
	GetMsgID() uint32
}

type Request struct {
	conn IConnection
	msg  IMessage
}

func NewRequest(conn IConnection, msg IMessage) *Request {
	return &Request{
		conn: conn,
		msg:  msg,
	}
}

func (r *Request) GetConn() IConnection {
	return r.conn
}

func (r *Request) GetData() []byte {
	return r.msg.GetData()
}

func (r *Request) GetMsgID() uint32 {
	return r.msg.GetID()
}
