package server

import (
	"app/internal/request"
	"app/internal/response"
	"io"
)

type HandlerError struct {
	error
	StatusCode response.StatusCode
	Message    string
}

type Handler func(w io.Writer, req *request.Request) *HandlerError

func (h HandlerError) Write(w io.Writer) {
	response.WriteStatusLine(w, h.StatusCode)
	headers := response.GetDefaultHeaders(len(h.Message))
	response.WriteHeaders(w, headers)
	response.WriteCRLF(w)
	w.Write([]byte(h.Message))
}
