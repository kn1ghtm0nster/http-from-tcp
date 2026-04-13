package server

import (
	"fmt"
	"net"
	"strconv"
	"sync/atomic"

	"github.com/kn1ghtm0nster/http-from-tcp/internal/response"
)

const (
	ServerStartedState serverState = iota
	ServerStoppedState
)

type serverState int

type Server struct {
	Listener net.Listener
	state serverState
	serverClosed atomic.Bool
}

func Serve(port int) (*Server, error) {
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		return nil, err
	}

	server := &Server{
		state: 			ServerStartedState,
		Listener: 		listener,
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

	err := response.WriteStatusLine(conn, response.StatusOK)
	if err != nil {
		fmt.Printf("ERROR WRITING STATUS LINE: %v", err)
		return
	}

	headers := response.GetDefaultHeaders(0)
	err = response.WriteHeaders(conn, headers)
	if err != nil {
		fmt.Printf("ERROR WRITING HEADERS: %v", err)
		return
	}

	_, err = conn.Write([]byte("\r\n"))
	if err != nil {
		fmt.Printf("ERROR WRITING BODY: %v", err)
		return
	}
}