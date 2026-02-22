package response

import (
	"fmt"
)

func (w *Writer) WriteChunkedBody(p []byte) (int, error) {
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
	n, err := w.writer.Write([]byte("0\r\n\r\n"))
	return n, err
}
