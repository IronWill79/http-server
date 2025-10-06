package server

import (
	"fmt"
	"log"
	"net"
	"sync/atomic"

	"github.com/IronWill79/http-server/internal/request"
	"github.com/IronWill79/http-server/internal/response"
)

type Server struct {
	closed   atomic.Bool
	handler  Handler
	listener net.Listener
}

func Serve(port int, handler Handler) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}
	s := &Server{
		handler:  handler,
		listener: listener,
	}
	go s.listen()
	return s, nil
}

func (s *Server) Close() error {
	s.closed.Store(true)
	return s.listener.Close()
}

func (s *Server) listen() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if s.closed.Load() {
				return
			}
			log.Printf("Error accepting connection: %v", err)
			continue
		}
		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	req, err := request.RequestFromReader(conn)
	if err != nil {
		hErr := &HandlerError{
			Status:  response.StatusCodeBadRequest,
			Message: err.Error(),
		}
		hErr.Write(conn)
		return
	}
	w := response.NewWriter(conn)
	s.handler(w, req)
}
