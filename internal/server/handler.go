package server

import (
	"app/internal/request"
	"app/internal/response"
)

type Handler func(w *response.Writer, req *request.Request)
