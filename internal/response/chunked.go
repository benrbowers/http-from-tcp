package response

import (
	"app/internal/headers"
	"errors"
	"fmt"
	"io"
)

func (w *Writer) WriteChunkedBody(p []byte) (int, error) {
	if w.state != writingBody {
		return 0, fmt.Errorf("Tried to write chunked body with invalid Writer state: %d", w.state)
	}

	if len(p) == 0 {
		return 0, nil
	}

	bytesWritten := 0

	chunkSize := len(p)
	n, err := fmt.Fprintf(w.writer, "%x\r\n", chunkSize)
	if err != nil {
		return bytesWritten, fmt.Errorf("Error writing data-size hex: %w", err)
	}
	bytesWritten += n

	p = append(p, '\r', '\n')
	n, err = w.writer.Write(p)
	if err != nil {
		return bytesWritten, fmt.Errorf("Error writing chunked body: %w", err)
	}
	bytesWritten += n

	return bytesWritten, nil
}

func (w *Writer) WriteChunkedBodyDone() (int, error) {
	if w.state != writingBody {
		return 0, fmt.Errorf("Tried to end chunked body with invalid Writer state: %d", w.state)
	}

	n, err := w.writer.Write([]byte("0\r\n"))

	if err == nil {
		w.state = writingTrailers
	}

	return n, err
}

func (w *Writer) WriteTrailers(trailers headers.Headers) error {
	if w.state != writingTrailers {
		return fmt.Errorf("Tried writing trailers with invalid Writer state: %d", w.state)
	}

	for trailer, trailerVal := range trailers {
		trailerLine := fmt.Sprintf("%s: %s\r\n", trailer, trailerVal)
		err := w.write([]byte(trailerLine))
		if err != nil {
			return fmt.Errorf("Error writing trailer: %w", err)
		}
	}

	err := w.write([]byte{'\r', '\n'})
	if err != nil {
		return err
	}

	w.state = writingDone
	return nil
}

func (w *Writer) WriteChunkedBodyFromReader(r io.Reader) (int, error) {
	bytesWritten := 0

	readBuffer := make([]byte, 1024)
	for !w.Done() {
		n, err := r.Read(readBuffer)
		if n > 0 {
			_, err := w.WriteChunkedBody(readBuffer[0:n])
			if err != nil {
				return bytesWritten, err
			}
			bytesWritten += n
		}
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return bytesWritten, err
		}
	}

	_, err := w.WriteChunkedBodyDone()
	return bytesWritten, err
}
