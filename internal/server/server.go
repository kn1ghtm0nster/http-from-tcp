package server

import (
	"bytes"
	"fmt"
	"io"
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

type HandlerError struct {
	StatusCode 	response.StatusCode
	Message 	string

}

type Handler func(w io.Writer, req *request.Request) *HandlerError 

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

	req, err := request.RequestFromReader(conn)
	if err != nil {
		fmt.Printf("ERROR PARSING REQUEST: %v", err)
		return
	}

	bodyBuf := &bytes.Buffer{}
	handlerErr := s.handler(bodyBuf, req)
	if handlerErr != nil {
		err := handlerErr.Write(conn)
		if err != nil {
			fmt.Printf("ERROR WRITING HANDLER ERROR: %v", err)
		}
		return
	}

	err = response.WriteStatusLine(conn, response.StatusOK)
	if err != nil {
		fmt.Printf("ERROR WRITING STATUS LINE: %v", err)
		return
	}

	headers := response.GetDefaultHeaders(bodyBuf.Len())
	err = response.WriteHeaders(conn, headers)
	if err != nil {
		fmt.Printf("ERROR WRITING HEADERS: %v", err)
		return
	}

	_, err = conn.Write(bodyBuf.Bytes())
	if err != nil {
		fmt.Printf("ERROR WRITING BODY: %v", err)
		return
	}
}

func (he *HandlerError) Write(w io.Writer) error {
	err := response.WriteStatusLine(w, he.StatusCode)
	if err != nil {
		return err
	}

	err = response.WriteHeaders(w, response.GetDefaultHeaders(len(he.Message)))
	if err != nil {
		return err
	}

	_, err = w.Write([]byte(he.Message))
	return err
}