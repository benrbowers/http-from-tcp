package server

import (
	"log"
	"net"
	"sync/atomic"
)

const address = "localhost:42069"

// Contains the state of the server
type Server struct {
	listener net.Listener
	closed   atomic.Bool
}

// Creates a net.Listener and returns a new
// Server instance. Starts listening for
// requests inside a goroutine.
func Serve(port int) (*Server, error) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return nil, err
	}

	newServer := &Server{
		listener: listener,
		closed:   atomic.Bool{},
	}

	go func() {
		newServer.listen()
	}()

	return newServer, nil
}

// Closes the listener and the server
func (s *Server) Close() error {
	err := s.listener.Close()
	if err != nil {
		return err
	}
	s.closed.Store(true)
	return nil
}

// Uses a loop to .Accept new connections as
// they come in, and handles each one in a new
// goroutine. I used an atomic.Bool to track
// whether the server is closed or not so that
// I can ignore connection errors after the
// server is closed.
func (s *Server) listen() {
	for !s.closed.Load() {
		log.Println("Waiting for request at", address)
		tcpConn, err := s.listener.Accept()
		if err != nil {
			log.Fatalf("Error accepting TCP connection: %v", err)
		}
		s.handle(tcpConn)
	}
}

const response = "HTTP/1.1 200 OK\r\n" +
	"Content-Type: text/plain\r\n" +
	"Content-Length: 13\r\n" +
	"\r\n" +
	"Hello World!\n"

// Handles a single connection by writing the following response and then closing the connection:
// For now, no matter what request is sent, the response will always be the same.
func (s *Server) handle(conn net.Conn) {
	_, err := conn.Write([]byte(response))
	if err != nil {
		log.Fatalf("Error writing response to connection: %v", err)
	}
	err = conn.Close()
	if err != nil {
		log.Fatalf("Error trying to close connection: %v", err)
	}
	log.Println("Connection received and response sent successfully.")
}
