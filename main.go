package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

const filename = "messages.txt"

func getLinesChannel(f io.ReadCloser) <-chan string {
	ch := make(chan string)

	go func() {
		buffer := make([]byte, 8)
		line := ""
		for {
			buffer_length, err := f.Read(buffer)
			if err != nil {
				if line != "" {
					ch <- line
				}
				close(ch)
				return
			}
			sections := strings.Split(string(buffer[:buffer_length]), "\n")
			line += sections[0]
			if len(sections) > 1 {
				ch <- line
				line = sections[1]
			}
		}
	}()

	return ch
}

func main() {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	line_channel := getLinesChannel(file)

	for line := range line_channel {
		fmt.Printf("read: %s\n", line)
	}
}
