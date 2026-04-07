package main

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

func getLinesChannel(f io.ReadCloser) <-chan string {
	lines := make(chan string)
	go func(f io.ReadCloser) {
		defer f.Close()
		defer close(lines)

		// while there is data to read from the file, read it and print it to the console
		// 8 bytes at a time
		fileBuffer := make([]byte, 8)

		currentLine := ""

		for {

			n, err := f.Read(fileBuffer)
			if err != nil {
				if errors.Is(err, io.EOF) {
					if len(currentLine) > 0 {
						lines <- currentLine
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
					lines <- currentLine + part
					currentLine = ""
				}
			}
		}
	}(f)
	return lines
}