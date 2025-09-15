package main

import (
	"fmt"
	"os"
)

const filename = "messages.txt"

func main() {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	buffer := make([]byte, 8)
	for {
		_, err := file.Read(buffer)
		if err != nil {
			break
		}
		fmt.Printf("read: %s\n", string(buffer))
	}
}
