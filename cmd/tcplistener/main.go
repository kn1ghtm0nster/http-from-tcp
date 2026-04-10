package main

import (
	"fmt"
	"log"
	"net"

	"github.com/kn1ghtm0nster/http-from-tcp/internal/request"
)


func main () {
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatalf("ERROR CREATING LISTENER: %v", err)
	}

	defer listener.Close()

	for {
		connection, err := listener.Accept()
		if err != nil {
			log.Fatalf("ERROR ACCEPTING CONNECTION: %v", err)
		}
		
		fmt.Println("NEW CONNECTION ACCEPTED!")
		req, err := request.RequestFromReader(connection)
		if err != nil {
			log.Printf("ERROR READING REQUEST: %v", err)
			connection.Close()
			continue
		}
		fmt.Println("Request line:")
		fmt.Printf("- Method: %s\n- Target: %s\n- Version: %s\n", req.RequestLine.Method, req.RequestLine.RequestTarget, req.RequestLine.HttpVersion)
		fmt.Println("CONNECTION CLOSED!")
	}
}