package znet

import (
	"errors"
	"fmt"
	"io"
	"net"

	"github.com/lexchenzhang/go-tiny-tcp/utils"
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
	server     IServer
	conn       *net.TCPConn
	connID     uint32
	isClose    bool
	exitChan   chan bool
	msgHandler IMsgHandler
	msgChan    chan []byte
}

func NewConnection(server IServer, conn *net.TCPConn, connID uint32, msgHandler IMsgHandler) *Connection {
	c := &Connection{
		server:     server,
		conn:       conn,
		connID:     connID,
		msgHandler: msgHandler,
		isClose:    false,
		exitChan:   make(chan bool, 1),
		msgChan:    make(chan []byte),
	}
	c.server.GetConnMgr().Add(c)
	return c
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
	c.exitChan <- true
	c.server.GetConnMgr().Remove(c)
	close(c.exitChan)
	close(c.msgChan)
	fmt.Println("drop connection:", c.connID)
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

	c.msgChan <- binaryMsg

	return nil
}

func (c *Connection) StartReader() {
	fmt.Println("[Reader Goroutine is running]")
	defer fmt.Println("[Reader Goroutine is exited]")
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

			if utils.GlobalObject.WorkerPoolSize > 0 {
				c.msgHandler.SendRequestToTaskQueue(req)
			} else {
				go c.msgHandler.DoMsgHandler(req)
			}

		}
	}
}

func (c *Connection) StartWriter() {
	fmt.Println("[Writer Goroutine is running]")
	defer fmt.Println(c.RemoteAddr().String(), " [Writer Goroutine is exited]")
	for {
		select {
		case msg := <-c.msgChan:
			if _, err := c.conn.Write(msg); err != nil {
				fmt.Println("write to conn failed, err:", err)
				return
			}
		case <-c.exitChan:
			return
		}
	}
}
