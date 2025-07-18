package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	udpAddr, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		panic("Error resolving UDP address: " + err.Error())
	}

	fmt.Println("Dialing UDP address...")
	udpConn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		panic("Error dialing UDP address: " + err.Error())
	}
	defer udpConn.Close()

	stdin := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")

		line, err := stdin.ReadBytes(byte('\n'))
		if err != nil {
			panic("Error reading line from Stdin: " + err.Error())
		}

		_, err = udpConn.Write(line)
		if err != nil {
			panic("Error writing line to UDP: " + err.Error())
		}
	}
}
