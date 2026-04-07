package main

import (
	"fmt"
	"log"
	"net"
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
		linesCh := getLinesChannel(connection)
		for line := range linesCh {
			fmt.Println(line)
		}
		fmt.Println("CONNECTION CLOSED!")
	}
}