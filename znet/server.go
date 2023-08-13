package znet

import (
	"errors"
	"fmt"
	"net"
)

type IServer interface {
	start()
	stop()
	StartAndServe()
}

type Server struct {
	Port      int
	IPVersion string
	IP        string
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

			dealConn := NewConnection(conn, cid, handleConn)
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

func NewServer(name string) IServer {
	return &Server{
		Port:      8888,
		IPVersion: "tcp4",
		IP:        "127.0.0.1",
	}
}

func handleConn(conn *net.TCPConn, data []byte, n int) error {
	fmt.Println("handle conn: ", conn.RemoteAddr())
	if _, err := conn.Write(data[:n]); err != nil {
		fmt.Println("write err: ", err)
		return errors.New("write err")
	}
	return nil
}
