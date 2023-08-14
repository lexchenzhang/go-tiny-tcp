package znet

import (
	"fmt"
	"log"
	"net"

	"github.com/lexchenzhang/go-tiny-tcp/utils"
)

type IServer interface {
	start()
	stop()
	StartAndServe()
	AddRouter(msgID uint32, router IRouter)
	GetConnMgr() IConnManager
	SetOnConnStart(f func(conn IConnection))
	SetOnConnStop(f func(conn IConnection))
	CallOnConnStart(conn IConnection)
	CallOnConnStop(conn IConnection)
}

type Server struct {
	Name        string
	Port        int
	IPVersion   string
	IP          string
	MsgHandler  IMsgHandler
	ConnMgr     IConnManager
	OnConnStart func(conn IConnection)
	OnConnStop  func(conn IConnection)
}

func (s *Server) start() {
	go func() {
		s.MsgHandler.StartWorkerPool()
		//1 gain TCP addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr err: ", err)
			return
		}

		//2 listen
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			panic(err)
		}

		var cid uint32
		cid = 0

		//3 handle conn
		for {
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("accept tcp err: ", err)
				continue
			}

			//TODO::inform client can't accept more connections
			if s.ConnMgr.Len() >= utils.GlobalObject.MaxConn {
				conn.Close()
				fmt.Println("too many connections, close conn: ", conn.RemoteAddr())
				continue
			} else {
				fmt.Println("max conn allowed is ", utils.GlobalObject.MaxConn)
			}

			fmt.Println("accept tcp conn: ", conn.RemoteAddr())

			dealConn := NewConnection(s, conn, cid, s.MsgHandler)
			go dealConn.Start()
			cid++
		}
	}()
}

func (s *Server) stop() {
	fmt.Println("[Server Stoped]")
	s.ConnMgr.ClearConn()
}

func (s *Server) StartAndServe() {
	s.start()
	s.printInfo()
	select {}
}

func (s *Server) AddRouter(msgID uint32, router IRouter) {
	s.MsgHandler.AddRouter(msgID, router)
}

func (s *Server) GetConnMgr() IConnManager {
	return s.ConnMgr
}

func NewServer() IServer {
	return &Server{
		Name:       utils.GlobalObject.Name,
		Port:       utils.GlobalObject.TcpProt,
		IPVersion:  "tcp4",
		IP:         utils.GlobalObject.Host,
		MsgHandler: NewMsgHandler(),
		ConnMgr:    NewConnManager(),
	}
}

func (s *Server) printInfo() {
	log.Println("Server ", s.Name, " is running on ", s.IP, ":", s.Port)
}

func (s *Server) SetOnConnStart(f func(IConnection)) {
	s.OnConnStart = f
}

func (s *Server) SetOnConnStop(f func(IConnection)) {
	s.OnConnStop = f
}

func (s *Server) CallOnConnStart(conn IConnection) {
	if s.OnConnStart != nil {
		s.OnConnStart(conn)
	}
}

func (s *Server) CallOnConnStop(conn IConnection) {
	if s.OnConnStop != nil {
		s.OnConnStop(conn)
	}
}
