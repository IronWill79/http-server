package main

import (
	"fmt"
	"os"
	"strings"
)

const filename = "messages.txt"

func main() {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	buffer := make([]byte, 8)
	line := ""
	for {
		buffer_length, err := file.Read(buffer)
		if err != nil {
			break
		}
		sections := strings.Split(string(buffer[:buffer_length]), "\n")
		line += sections[0]
		if len(sections) > 1 {
			fmt.Printf("read: %s\n", line)
			line = sections[1]
		}
	}
	if line != "" {
		fmt.Printf("read: %s\n", line)
	}
}
