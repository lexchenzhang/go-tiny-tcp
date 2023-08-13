package znet

type IRequest interface {
	GetConn() IConnection
	GetData() []byte
}

type Request struct {
	conn IConnection
	data []byte
}

func NewRequest(conn IConnection, data []byte) *Request {
	return &Request{
		conn: conn,
		data: data,
	}
}

func (r *Request) GetConn() IConnection {
	return r.conn
}

func (r *Request) GetData() []byte {
	return r.data
}
