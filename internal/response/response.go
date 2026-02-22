package response

import (
	"app/internal/headers"
	"strconv"
)

type StatusCode int

const StatusOK StatusCode = 200
const StatusBadRequest StatusCode = 400
const StatusInternalError StatusCode = 500

func GetDefaultHeaders(contentLen int) headers.Headers {
	defaultHeaders := headers.Headers{}
	defaultHeaders.Set("Content-Length", strconv.Itoa(contentLen))
	defaultHeaders.Set("Connection", "close")
	defaultHeaders.Set("Content-Type", "text/plain")

	return defaultHeaders
}

func GetChunkedHeaders() headers.Headers {
	chunkedHeaders := headers.Headers{}
	chunkedHeaders.Set("Transfer-Encoding", "chunked")
	chunkedHeaders.Set("Connection", "close")

	return chunkedHeaders
}
