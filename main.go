package main

import (
	"fmt"
	"os"
)


func main () {
	byteData, err := os.Open("messages.txt")
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}

	lines := getLinesChannel(byteData)
	for line := range lines {
		fmt.Printf("read: %s\n", line)
	}
}	