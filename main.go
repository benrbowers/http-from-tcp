package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	input, err := os.Open("messages.txt")
	if err != nil {
		panic("Could not open messages.txt:\n" + err.Error())
	}

	word := make([]byte, 8)
	for {
		n, err := input.Read(word)

		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}

		if n > 0 {
			fmt.Printf("read: %s\n", string(word[0:n]))
		}
	}
}
