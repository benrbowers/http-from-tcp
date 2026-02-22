package response

import (
	"app/internal/headers"
	"fmt"
	"io"
)

type writerState int

const (
	writingStatusLine writerState = iota
	writingHeaders
	writingBody
	writingDone
)

type Writer struct {
	writer io.Writer
	state  writerState
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{writer: w}
}

func (w *Writer) write(p []byte) error {
	_, err := w.writer.Write(p)
	return err
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	if w.state != writingStatusLine {
		return fmt.Errorf("Tried to write status-line with invalid Writer state: %d", w.state)
	}

	switch statusCode {
	case StatusOK:
		err := w.write([]byte("HTTP/1.1 200 OK\r\n"))
		if err != nil {
			return err
		}
	case StatusBadRequest:
		err := w.write([]byte("HTTP/1.1 400 Bad Request\r\n"))
		if err != nil {
			return err
		}
	case StatusInternalError:
		err := w.write([]byte("HTTP/1.1 500 Internal Server Error\r\n"))
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("Unknown status code: %d", statusCode)
	}

	w.state = writingHeaders
	return nil
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	if w.state != writingHeaders {
		return fmt.Errorf("Tried to write headers with invalid Writer state: %d", w.state)
	}

	for name, value := range headers {
		err := w.write([]byte(name + ": " + value + "\r\n"))
		if err != nil {
			return err
		}
	}

	err := w.write([]byte{'\r', '\n'})
	if err != nil {
		return err
	}

	w.state = writingBody
	return nil
}

func (w *Writer) WriteBody(data []byte) (int, error) {
	if w.state != writingBody {
		return 0, fmt.Errorf("Tried to write body with invalid Writer state: %d", w.state)
	}

	n, err := w.writer.Write(data)
	if err != nil {
		return n, err
	}

	w.state = writingDone
	return n, nil
}

func (w *Writer) Done() bool {
	return w.state == writingDone
}
