package response

import (
	"app/internal/headers"
	"fmt"
	"io"
	"strconv"
)

type StatusCode int

const StatusOK StatusCode = 200
const StatusBadRequest StatusCode = 400
const StatusInternalError StatusCode = 500

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	switch statusCode {
	case StatusOK:
		_, err := w.Write([]byte("HTTP/1.1 200 OK\r\n"))
		if err != nil {
			return err
		}
	case StatusBadRequest:
		_, err := w.Write([]byte("HTTP/1.1 400 Bad Request\r\n"))
		if err != nil {
			return err
		}
	case StatusInternalError:
		_, err := w.Write([]byte("HTTP/1.1 500 Internal Server Error\r\n"))
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("Unknown status code: %d", statusCode)
	}

	return nil
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	defaultHeaders := headers.Headers{}
	defaultHeaders.Set("Content-Length", strconv.Itoa(contentLen))
	defaultHeaders.Set("Connection", "close")
	defaultHeaders.Set("Content-Type", "text/plain")

	return defaultHeaders
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	for name, value := range headers {
		_, err := w.Write([]byte(name + ": " + value + "\r\n"))
		if err != nil {
			return err
		}
	}

	return nil
}

func WriteCRLF(w io.Writer) error {
	_, err := w.Write([]byte{'\r', '\n'})
	return err
}
