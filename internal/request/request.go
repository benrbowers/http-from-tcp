package request

import (
	"app/internal/headers"
	"errors"
	"fmt"
	"io"
	"strings"
)

type requestStatus int

const (
	requestInitialized requestStatus = iota
	requestParsingHeaders
	requestDone
)

type Request struct {
	RequestLine RequestLine
	Status      requestStatus
	Headers     headers.Headers
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func (r *Request) parse(data []byte) (int, error) {
	totalBytesParsed := 0
	for r.Status != requestDone {
		n, err := r.parseSingle(data[totalBytesParsed:])
		if err != nil {
			return 0, err
		}
		if n == 0 {
			break
		}
		totalBytesParsed += n
	}

	return totalBytesParsed, nil
}

func (r *Request) parseSingle(data []byte) (int, error) {
	switch r.Status {
	case requestInitialized:
		requestLine, bytesParsed, err := parseRequestLine(data)

		if err != nil {
			return 0, err
		}

		if bytesParsed > 0 {
			r.RequestLine = *requestLine
			r.Status = requestParsingHeaders
		}

		return bytesParsed, nil
	case requestParsingHeaders:
		bytesParsed, done, err := r.Headers.Parse(data)

		if err != nil {
			return 0, err
		}

		if done {
			r.Status = requestDone
		}

		return bytesParsed, nil
	case requestDone:
		return 0, fmt.Errorf("error: trying to read data in a done state.")
	default:
		return 0, fmt.Errorf("Unknown request status.")
	}
}

const bufferSize int = 8

func RequestFromReader(reader io.Reader) (*Request, error) {
	buf := make([]byte, bufferSize)
	readToIndex := 0

	newRequest := &Request{
		Headers: headers.Headers{},
	}

	for newRequest.Status != requestDone {
		if readToIndex >= len(buf) {
			newBuf := make([]byte, len(buf)*2)
			copy(newBuf, buf)
			buf = newBuf
		}

		// Reader can read to a SUBSLICE, very cool
		readSize, err := reader.Read(buf[readToIndex:])
		if err != nil {
			if errors.Is(err, io.EOF) {
				if newRequest.Status != requestDone {
					return nil, fmt.Errorf("Incomplete request, in state: %d, read n bytes on EOF: %d", newRequest.Status, readSize)
				}
				break
			}
			return nil, err
		}

		readToIndex += readSize

		numBytesParsed, err := newRequest.parse(buf[:readToIndex])
		if err != nil {
			return nil, err
		}

		// Shifting data out to reuse buffer, in two
		// simple lines. Also very cool.
		copy(buf, buf[numBytesParsed:])
		readToIndex -= numBytesParsed
	}

	return newRequest, nil
}

// parseRequestLine parses an HTTP request line from a string of bytes.
// pareRequestLine returns the RequestLine, bytes consumed, and optional error.
func parseRequestLine(request []byte) (*RequestLine, int, error) {
	requestText := string(request)

	crlfIndex := strings.Index(requestText, "\r\n")
	if crlfIndex == -1 {
		return nil, 0, nil // No CRLF, so need to read more.
	}

	requestText = requestText[:crlfIndex]

	requestParts := strings.Split(requestText, " ")
	if len(requestParts) != 3 {
		return nil, 0, errors.New("Invalid request line.")
	}

	method := requestParts[0]
	if !isCapitalOnly(method) {
		return nil, 0, fmt.Errorf(
			`Invalid request method: "%s". Method may only contain captial letters.`,
			method,
		)
	}

	requestTarget := requestParts[1]

	httpVersion := requestParts[2]
	httpVersionParts := strings.Split(httpVersion, "/")
	if len(httpVersionParts) != 2 {
		return nil, 0, fmt.Errorf(
			`Invalid HTTP version: "%s". Required format: HTTP-name "/" DIGIT "." DIGIT`,
			httpVersion,
		)
	}
	if httpVersionParts[0] != "HTTP" {
		return nil, 0, fmt.Errorf(
			`Invalid HTTP version name: "%s". Only HTTP/1.1 is supported.`,
			httpVersionParts[0],
		)
	}
	if httpVersionParts[1] != "1.1" {
		return nil, 0, fmt.Errorf(
			`Invalid HTTP version number: "%s". Only HTTP/1.1 is supported.`,
			httpVersionParts[1],
		)
	}

	return &RequestLine{
			Method:        method,
			RequestTarget: requestTarget,
			HttpVersion:   httpVersionParts[1],
		},
		len(requestText) + 2, // + 2 for CRLF
		nil
}

func isCapitalOnly(text string) bool {
	for _, char := range text {
		if char < 'A' || char > 'Z' {
			return false
		}
	}
	return true
}
