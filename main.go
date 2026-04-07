package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)


func main () {
	byteData, err := os.Open("messages.txt")
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}

	defer byteData.Close()

	// variable to hold contents of current line 
	currentLine := ""

	// while there is data to read from the file, read it and print it to the console
	// 8 bytes at a time
	fileBuffer := make([]byte, 8)
	for {
		n, err := byteData.Read(fileBuffer)
		if err != nil {
			if errors.Is(err, io.EOF) {
				if len(currentLine) > 0 {
					fmt.Printf("read: %s\n", currentLine)
					currentLine = ""
				}
				return
			}
			fmt.Printf("Error reading file: %v\n", err)
			return
		}

		str := string(fileBuffer[:n])

		parts := strings.Split(str, "\n")
		for i, part := range parts {
			if i == len(parts)-1 {
				currentLine += part
			} else {
				fmt.Printf("read: %s\n", currentLine + part)
				currentLine = ""
			}
		}
	}
}	