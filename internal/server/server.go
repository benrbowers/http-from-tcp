package server

import (
	"app/internal/request"
	"app/internal/response"
	"fmt"
	"log"
	"net"
	"sync/atomic"
)

const address = "localhost:42069"

// Contains the state of the server
type Server struct {
	listener net.Listener
	handler  Handler
	closed   atomic.Bool
}

// Creates a net.Listener and returns a new
// Server instance. Starts listening for
// requests inside a goroutine.
func Serve(port int, handler Handler) (*Server, error) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return nil, err
	}

	newServer := &Server{
		listener: listener,
		handler:  handler,
		closed:   atomic.Bool{},
	}

	go newServer.listen()

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

// Handles a single connection by writing the response and then closing the connection:
func (s *Server) handle(conn net.Conn) {
	rWriter := response.NewWriter(conn)

	req, err := request.RequestFromReader(conn)
	if err != nil {
		body := fmt.Appendf(nil, "Error parsing request: %v", err)
		rWriter.WriteStatusLine(response.StatusBadRequest)
		headers := response.GetDefaultHeaders(len(body))
		rWriter.WriteHeaders(headers)
		rWriter.WriteBody(body)

		return
	}

	s.handler(rWriter, req)

	err = conn.Close()
	if err != nil {
		log.Fatalf("Error trying to close connection: %v", err)
	}
	log.Println("Connection received and response sent successfully.")
}
