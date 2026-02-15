package server

import (
	"app/internal/response"
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

// Handles a single connection by writing the following response and then closing the connection:
// For now, no matter what request is sent, the response will always be the same.
func (s *Server) handle(conn net.Conn) {
	err := response.WriteStatusLine(conn, response.StatusOK)
	if err != nil {
		log.Fatalf("Error writing status-line to connection: %v", err)
	}
	defaultHeaders := response.GetDefaultHeaders(0)
	err = response.WriteHeaders(conn, defaultHeaders)
	if err != nil {
		log.Fatalf("Error writing headers to connection: %v", err)
	}
	err = conn.Close()
	if err != nil {
		log.Fatalf("Error trying to close connection: %v", err)
	}
	log.Println("Connection received and response sent successfully.")
}
