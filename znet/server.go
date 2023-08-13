package znet

import (
	"fmt"
	"net"
)

type IServer interface {
	start()
	stop()
	StartAndServe()
	AddRouter(router IRouter)
}

type Server struct {
	Port      int
	IPVersion string
	IP        string
	Router    IRouter
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

			dealConn := NewConnection(conn, cid, s.Router)
			go dealConn.Start()
			cid++
		}
	}()
}

func (s *Server) stop() {}

func (s *Server) StartAndServe() {
	s.start()
	select {}
}

func (s *Server) AddRouter(router IRouter) {
	s.Router = router
}

func NewServer(name string) IServer {
	return &Server{
		Port:      8888,
		IPVersion: "tcp4",
		IP:        "127.0.0.1",
		Router:    nil,
	}
}
