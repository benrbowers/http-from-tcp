package main

import (
	"fmt"
	"io"
	"os"
	"slices"
)

func main() {
	input, err := os.Open("messages.txt")
	if err != nil {
		panic("Could not open messages.txt:\n" + err.Error())
	}

	currentLine := ""
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
			newLine := slices.Index(word, byte('\n'))
			if newLine == -1 {
				currentLine += string(word[0:n])
			} else {
				currentLine += string(word[0:min(newLine, n)])
				fmt.Printf("read: %s\n", currentLine)

				currentLine = ""
				currentLine += string(word[newLine+1:])
			}
		}
	}
}
