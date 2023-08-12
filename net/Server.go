package net

type IServer interface {
	Start()
	Stop()
	Serve()
}

type Server struct {
	Port int
	Addr string
}

func (s *Server) Start() {}

func (s *Server) Stop() {}

func (s *Server) Serve() {}

func NewServer(name string) IServer {
	return &Server{
		Port: 8888,
		Addr: "127.0.0.1",
	}
}
