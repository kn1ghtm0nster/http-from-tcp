package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main () {
	addr, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		log.Printf("ERROR RESOLVING UDP ADDR: %v\n", err)
		return
	}
	fmt.Println("NEW UDP ADDRESS RESOLVED")

	conn, err := net.DialUDP(addr.Network(), nil, addr)
	if err != nil {
		log.Printf("ERROR DIALING UDP: %v\n", err)
		return
	}
	fmt.Println("UDP CONNECTION ESTABLISHED")
	defer conn.Close()

	// create a new buffer reader that reads from os.stdin
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		inputStr, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("ERROR READING FROM STDIN: %v\n", err)
		}

		_, err = conn.Write([]byte(inputStr))
		if err != nil {
			log.Printf("ERROR WRITING TO UDP CONNECTION: %v\n", err)
		}
	}
}