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

	defer byteData.Close()

	// while there is data to read from the file, read it and print it to the console
	// 8 bytes at a time
	fileBuffer := make([]byte, 8)
	for {
		n, err := byteData.Read(fileBuffer)
		if err != nil {
			if err.Error() == "EOF" {
				return
			}
			fmt.Printf("Error reading file: %v\n", err)
			return
		}
		fmt.Printf("read: %s\n", string(fileBuffer[:n]))
	}
}