package server

import (
	"fmt"
	"log"
	"net"
	"sync/atomic"

	"github.com/IronWill79/http-server/internal/response"
)

type Server struct {
	closed   atomic.Bool
	listener net.Listener
}

func Serve(port int) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}
	s := Server{
		listener: listener,
	}
	go func() {
		s.listen()
	}()
	return &s, nil
}

func (s *Server) Close() error {
	s.closed.Store(true)
	return s.listener.Close()
}

func (s *Server) listen() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go func(c net.Conn) {
			if s.closed.Load() {
				return
			}
			s.handle(c)
		}(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	response.WriteStatusLine(conn, 200)
	h := response.GetDefaultHeaders(0)
	response.WriteHeaders(conn, h)
	conn.Close()
}
