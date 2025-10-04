package main

import (
	"fmt"
	"net"

	"github.com/IronWill79/http-server/internal/request"
)

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

		req, err := request.RequestFromReader(connection)
		if err != nil {
			panic(err)
		}
		fmt.Println("Request line:")
		fmt.Printf("- Method: %s\n", req.RequestLine.Method)
		fmt.Printf("- Target: %s\n", req.RequestLine.RequestTarget)
		fmt.Printf("- Version: %s\n", req.RequestLine.HttpVersion)
		fmt.Println("Headers:")
		for k, v := range req.Headers {
			fmt.Printf("- %s: %s\n", k, v)
		}
		fmt.Printf("Body:\n%s\n", string(req.Body))

		connection.Close()
		fmt.Println("Connection closed")
	}
}
