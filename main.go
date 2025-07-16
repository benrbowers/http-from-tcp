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

	lines := getLinesChannel(input)

	for line := range lines {
		fmt.Printf("read: %s\n", line)
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	ch := make(chan string)

	go func() {
		defer close(ch)
		defer f.Close()

		currentLine := ""
		word := make([]byte, 8)
		for {
			n, err := f.Read(word)

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
					ch <- currentLine

					currentLine = ""
					currentLine += string(word[newLine+1:])
				}
			}
		}
	}()

	return ch
}
