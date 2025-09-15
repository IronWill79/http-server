package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	address, err := net.ResolveUDPAddr("udp", "127.0.0.1:42069")
	if err != nil {
		panic(err)
	}
	connection, err := net.DialUDP("udp", nil, address)
	if err != nil {
		panic(err)
	}
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading string: %v\n", err)
		}
		_, err = connection.Write([]byte(line))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error writing to UDP: %v\n", err)
		}
	}
}
