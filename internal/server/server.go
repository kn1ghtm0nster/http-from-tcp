package server

import (
	"fmt"
	"net"
	"strconv"
	"sync/atomic"

	"github.com/kn1ghtm0nster/http-from-tcp/internal/request"
	"github.com/kn1ghtm0nster/http-from-tcp/internal/response"
)

const (
	ServerStartedState serverState = iota
	ServerStoppedState
)

type serverState int

type Server struct {
	Listener net.Listener
	handler Handler
	state serverState
	serverClosed atomic.Bool
}

type Handler func(res *response.Writer, req *request.Request)

func Serve(port int, handler Handler) (*Server, error) {
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		return nil, err
	}

	server := &Server{
		state: 			ServerStartedState,
		Listener: 		listener,
		handler:		handler,
	}

	go server.listen()

	return server, nil
}

func (s *Server) Close() error {
	s.serverClosed.Store(true)
	err := s.Listener.Close()
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) listen() {
	for {
		conn, err := s.Listener.Accept()
		if err != nil {
			if s.serverClosed.Load() {
				return
			}
			fmt.Printf("CONNECTION ERROR: %v", err)
			continue
		}
		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()

	writer := response.NewWriter(conn)
	req, err := request.RequestFromReader(conn)
	if err != nil {
		writer.WriteStatusLine(response.BadRequest)
		writer.WriteHeaders(response.GetDefaultHeaders(len("Bad Request")))
		writer.WriteBody([]byte("Bad Request"))
		return
	}

	s.handler(writer, req)
}