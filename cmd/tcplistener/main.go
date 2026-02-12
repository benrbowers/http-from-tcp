package main

import (
	"app/internal/request"
	"fmt"
	"net"
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

	req, err := request.RequestFromReader(tcpConn)
	if err != nil {
		panic("Error parsing request from TCP connection: " + err.Error())
	}

	fmt.Println("Request line:")
	fmt.Println("- Method:", req.RequestLine.Method)
	fmt.Println("- Target:", req.RequestLine.RequestTarget)
	fmt.Println("- Version:", req.RequestLine.HttpVersion)
	fmt.Println("Headers:")
	for fieldName, fieldValue := range req.Headers {
		fmt.Printf("- %s: %s\n", fieldName, fieldValue)
	}
}
