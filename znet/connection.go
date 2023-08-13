package znet

import (
	"fmt"
	"net"
)

type IConnection interface {
	Start()
	Stop()
	GetTCPConn() *net.TCPConn
	GetConnID() uint32
	RemoteAddr() net.Addr
	Send(data []byte) error
}

type HandleFunc func(*net.TCPConn, []byte, int) error

type Connection struct {
	conn     *net.TCPConn
	connID   uint32
	isClose  bool
	exitChan chan bool
	router   IRouter
}

func NewConnection(conn *net.TCPConn, connID uint32, router IRouter) *Connection {
	return &Connection{
		conn:     conn,
		connID:   connID,
		router:   router,
		isClose:  false,
		exitChan: make(chan bool, 1),
	}
}

func (c *Connection) Start() {
	if c.isClose {
		return
	}
	go c.StartReader()
	go c.StartWriter()
}
func (c *Connection) Stop() {
	if c.isClose {
		return
	}
	c.isClose = true
	c.conn.Close()
	close(c.exitChan)
}
func (c *Connection) GetTCPConn() *net.TCPConn {
	return c.conn
}
func (c *Connection) GetConnID() uint32 {
	return c.connID
}
func (c *Connection) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}
func (c *Connection) Send(data []byte) error {
	return nil
}
func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is running...")
	defer fmt.Println("connID = ", c.connID, " Reader is exit, remote addr is ", c.RemoteAddr().String())
	defer c.Stop()
	for {
		buf := make([]byte, 1024)
		_, err := c.conn.Read(buf)
		if err != nil {
			fmt.Println("connID = ", c.connID, " Reader is exit, remote addr is ", c.RemoteAddr().String(), " err is ", err)
			continue
		}

		req := NewRequest(c, buf)

		go func(request IRequest) {
			c.router.PreHandleFunc(request)
			c.router.HandleFunc(request)
			c.router.PostHandleFunc(request)
		}(req)
	}
}
func (c *Connection) StartWriter() {}
