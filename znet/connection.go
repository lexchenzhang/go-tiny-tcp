package znet

import (
	"errors"
	"fmt"
	"io"
	"net"
)

type IConnection interface {
	Start()
	Stop()
	GetTCPConn() *net.TCPConn
	GetConnID() uint32
	RemoteAddr() net.Addr
	SendMsg(msgID uint32, data []byte) error
}

type HandleFunc func(*net.TCPConn, []byte, int) error

type Connection struct {
	conn       *net.TCPConn
	connID     uint32
	isClose    bool
	exitChan   chan bool
	msgHandler IMsgHandler
}

func NewConnection(conn *net.TCPConn, connID uint32, msgHandler IMsgHandler) *Connection {
	return &Connection{
		conn:       conn,
		connID:     connID,
		msgHandler: msgHandler,
		isClose:    false,
		exitChan:   make(chan bool, 1),
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

func (c *Connection) SendMsg(msgID uint32, data []byte) error {
	if c.isClose {
		return errors.New("connection closed")
	}
	dp := NewDataPack()

	binaryMsg, err := dp.Pack(NewMessage(msgID, data))
	if err != nil {
		fmt.Println("pack error, msg id = ", msgID)
		return errors.New("pack msg error")
	}

	if _, err := c.conn.Write(binaryMsg); err != nil {
		fmt.Println("send msg error, msg id = ", msgID)
		return errors.New("send msg error")
	}
	return nil
}

func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is running...")
	defer c.Stop()
	for {
		dp := NewDataPack()
		buf := make([]byte, dp.GetHeadLen())
		_, err := io.ReadFull(c.conn, buf)
		if err != nil {
			fmt.Println("read from conn failed, err:", err)
			break
		}
		msgHeader, err := dp.Unpack(buf)
		if err != nil {
			fmt.Println("unpack data failed, err:", err)
			break
		}
		if msgHeader.GetDataLen() > 0 {
			// read content of pack data
			msg := msgHeader.(*Message)
			msg.Data = make([]byte, msg.GetDataLen())
			_, err := io.ReadFull(c.conn, msg.Data)
			if err != nil {
				fmt.Println("server unpack data err", err)
				return
			}

			fmt.Println("->Recv MsgID[", msg.GetID(), "] Len[", msg.GetDataLen(), "] Data[", string(msg.GetData()), "]")

			req := NewRequest(c, msg)

			go c.msgHandler.DoMsgHandler(req)
		}
	}
}

func (c *Connection) StartWriter() {}
