package main

import (
	"fmt"
	"io"
	"net"
	"slices"
)

func main() {
	// boot.dev requested this port not me 😭
	tcpListener, err := net.Listen("tcp", "localhost:42069")
	if err != nil {
		panic("Error trying to listen on localhost:42069: " + err.Error())
	}

	fmt.Println("Wating for TCP connection...")
	tcpConn, err := tcpListener.Accept()
	if err != nil {
		panic("Error while waiting to accept tcp connection: " + err.Error())
	}
	fmt.Println("TCP connection established.")
	defer tcpListener.Close()

	lines := getLinesChannel(tcpConn)

	for line := range lines {
		fmt.Println(line)
	}

	fmt.Println("The connection has been closed.")
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
					ch <- currentLine
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
