package main

import (
	h "app/internal/headers"
	"app/internal/request"
	"app/internal/response"
	"app/internal/server"
	"crypto/sha256"
	"fmt"
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
	if strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin/stream") {
		httpbinStreamHandler(w, req)
		return
	}
	if strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin/html") {
		httpbinHtmlHandler(w, req)
		return
	}
}
var httpbinStreamHandler server.Handler = func(w *response.Writer, req *request.Request) {
	target := strings.TrimPrefix(req.RequestLine.RequestTarget, "/httpbin/")
	resp, err := fetchHttpbin(target)
	if err != nil {
		handle500(w, req)
		return
	}
	defer resp.Body.Close()

	w.WriteStatusLine(response.StatusOK)
	headers := response.GetChunkedHeaders()
	headers.Set("Content-Type", "application/json")
	w.WriteHeaders(headers)

	_, err = w.WriteChunkedBodyFromReader(resp.Body)
	if err != nil {
		log.Fatalf("Error writing chunked body from reader: %v", err)
		return
	}

	err = w.WriteTrailers(h.Headers{})
	if err != nil {
		log.Fatalf("Error writing trailers: %v", err)
		return
	}
}
var httpbinHtmlHandler server.Handler = func(w *response.Writer, req *request.Request) {
	target := strings.TrimPrefix(req.RequestLine.RequestTarget, "/httpbin/")
	resp, err := fetchHttpbin(target)
	if err != nil {
		handle500(w, req)
		return
	}
	defer resp.Body.Close()

	w.WriteStatusLine(response.StatusOK)
	headers := response.GetChunkedHeaders()
	headers.Set("Content-Type", "text/html")
	headers.Set("Trailer", "X-Content-SHA256")
	headers.Set("Trailer", "X-Content-Length")
	w.WriteHeaders(headers)

	hasher := sha256.New()
	n, err := w.WriteChunkedBodyFromReader(io.TeeReader(resp.Body, hasher))
	if err != nil {
		log.Fatalf("Error writing chunked body from reader: %v", err)
		return
	}

	trailers := h.Headers{}
	trailers.Set("X-Content-Length", fmt.Sprintf("%d", n))
	trailers.Set("X-Content-SHA256", fmt.Sprintf("%x", hasher.Sum(nil)))

	err = w.WriteTrailers(trailers)
	if err != nil {
		log.Fatalf("Error writing trailers: %v", err)
		return
	}
}

func fetchHttpbin(target string) (*http.Response, error) {
	resp, err := http.Get("https://httpbin.org/" + target)
	if err != nil {
		log.Printf(
			"Error fetching %s: %v",
			"https://httpbin.org/"+target,
			err,
		)
	}
	return resp, err
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
