package main

import (
	"app/internal/request"
	"app/internal/response"
	"app/internal/server"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
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

func handler(w *response.Writer, req *request.Request) {
	if strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin") {
		httpbinProxyHandler(w, req)
		return
	}
	if req.RequestLine.RequestTarget == "/yourproblem" {
		handle200(w, req)
		return
	}
	if req.RequestLine.RequestTarget == "/myproblem" {
		handle500(w, req)
		return
	}

	handle200(w, req)
}

var httpbinProxyHandler server.Handler = func(w *response.Writer, req *request.Request) {
	proxyPath := "/httpbin/"
	if strings.HasPrefix(req.RequestLine.RequestTarget, proxyPath) {
		target := req.RequestLine.RequestTarget[len(proxyPath):]
		resp, err := http.Get("https://httpbin.org/" + target)
		if err != nil {
			log.Printf(
				"Error fetching %s: %v",
				"https://httpbin.org/"+target,
				err,
			)
			handle500(w, req)
			return
		}

		w.WriteStatusLine(response.StatusOK)
		headers := response.GetChunkedHeaders()
		headers.Set("Content-Type", "application/json")
		w.WriteHeaders(headers)

		readBuffer := make([]byte, 1024)
		for !w.Done() {
			n, err := resp.Body.Read(readBuffer)
			if err != nil {
				if errors.Is(err, io.EOF) {
					w.WriteChunkedBodyDone()
					break
				}
				log.Printf("Error reading body from httpbin.org: %v", err)
				break
			}

			if n == 0 {
				w.WriteChunkedBodyDone()
				break
			}
			w.WriteChunkedBody(readBuffer[0:n])
		}
	}
}

var handle200 server.Handler = func(w *response.Writer, _ *request.Request) {
	w.WriteStatusLine(response.StatusOK)
	headers := response.GetDefaultHeaders(len(okHtml))
	headers.Replace("Content-Type", "text/html")
	w.WriteHeaders(headers)
	w.WriteBody(okHtml)
}
var handle400 server.Handler = func(w *response.Writer, _ *request.Request) {
	w.WriteStatusLine(response.StatusBadRequest)
	headers := response.GetDefaultHeaders(len(badRequestHtml))
	headers.Replace("Content-Type", "text/html")
	w.WriteHeaders(headers)
	w.WriteBody(badRequestHtml)
}
var handle500 server.Handler = func(w *response.Writer, _ *request.Request) {
	w.WriteStatusLine(response.StatusInternalError)
	headers := response.GetDefaultHeaders(len(internalErrorHtml))
	headers.Replace("Content-Type", "text/html")
	w.WriteHeaders(headers)
	w.WriteBody(internalErrorHtml)
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
