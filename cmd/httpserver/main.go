package main

import (
	"app/internal/request"
	"app/internal/response"
	"app/internal/server"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const port = 42069

func main() {
	server, err := server.Serve(port, handler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}

var handler server.Handler = func(w *response.Writer, req *request.Request) {
	switch req.RequestLine.RequestTarget {
	case "/yourproblem":
		w.WriteStatusLine(response.StatusBadRequest)
		headers := response.GetDefaultHeaders(len(badRequestHtml))
		headers.Replace("Content-Type", "text/html")
		w.WriteHeaders(headers)
		w.WriteBody(badRequestHtml)
	case "/myproblem":
		w.WriteStatusLine(response.StatusInternalError)
		headers := response.GetDefaultHeaders(len(internalErrorHtml))
		headers.Replace("Content-Type", "text/html")
		w.WriteHeaders(headers)
		w.WriteBody(internalErrorHtml)
	default:
		w.WriteStatusLine(response.StatusOK)
		headers := response.GetDefaultHeaders(len(okHtml))
		headers.Replace("Content-Type", "text/html")
		w.WriteHeaders(headers)
		w.WriteBody(okHtml)
	}
}

var badRequestHtml = []byte(`<html>
  <head>
    <title>400 Bad Request</title>
  </head>
  <body>
    <h1>Bad Request</h1>
    <p>Your request honestly kinda sucked.</p>
  </body>
</html>`)
var internalErrorHtml = []byte(`<html>
  <head>
    <title>500 Internal Server Error</title>
  </head>
  <body>
    <h1>Internal Server Error</h1>
    <p>Okay, you know what? This one is on me.</p>
  </body>
</html>`)
var okHtml = []byte(`<html>
  <head>
    <title>200 OK</title>
  </head>
  <body>
    <h1>Success!</h1>
    <p>Your request was an absolute banger.</p>
  </body>
</html>`)
