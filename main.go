package main

import (
	"fmt"
	"io"
	"net"
	"strings"
)

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
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	for {
		connection, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		fmt.Println("Connection accepted on port 42069")

		line_channel := getLinesChannel(connection)

		for line := range line_channel {
			fmt.Println(line)
		}

		connection.Close()
		fmt.Println("Connection closed")
	}
}
