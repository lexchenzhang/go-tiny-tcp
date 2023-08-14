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
}

type Server struct {
	Name       string
	Port       int
	IPVersion  string
	IP         string
	MsgHandler IMsgHandler
}

func (s *Server) start() {
	go func() {
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
			fmt.Println("accept tcp conn: ", conn.RemoteAddr())

			dealConn := NewConnection(conn, cid, s.MsgHandler)
			go dealConn.Start()
			cid++
		}
	}()
}

func (s *Server) stop() {}

func (s *Server) StartAndServe() {
	s.start()
	s.printInfo()
	select {}
}

func (s *Server) AddRouter(msgID uint32, router IRouter) {
	s.MsgHandler.AddRouter(msgID, router)
}

func NewServer() IServer {
	return &Server{
		Name:       utils.GlobalObject.Name,
		Port:       utils.GlobalObject.TcpProt,
		IPVersion:  "tcp4",
		IP:         utils.GlobalObject.Host,
		MsgHandler: NewMsgHandler(),
	}
}

func (s *Server) printInfo() {
	log.Println("Server ", s.Name, " is running on ", s.IP, ":", s.Port)
}
