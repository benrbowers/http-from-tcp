package request

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

type requestStatus int

const (
	RequestInitialized requestStatus = iota
	RequestDone
)

type Request struct {
	RequestLine RequestLine
	Status      requestStatus
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func (r *Request) parse(data []byte) (int, error) {
	requestLine, bytesParsed, err := parseRequestLine(data)

	if err != nil {
		return 0, err
	}

	if bytesParsed > 0 {
		r.RequestLine = requestLine
		r.Status = RequestDone
	}

	return bytesParsed, nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	requestBytes := []byte{}
	readBytes := make([]byte, 8)

	newRequest := &Request{}

	for newRequest.Status != RequestDone {
		readSize, err := reader.Read(readBytes)
		if err != nil {
			return newRequest, err
		}

		requestBytes = append(requestBytes, readBytes[0:readSize]...)

		_, err = newRequest.parse(requestBytes)
		// fmt.Println("Bytes parsed:", bytesParsed)
		if err != nil {
			return newRequest, err
		}
	}

	return newRequest, nil
}

// parseRequestLine parses an HTTP request line from a string of bytes.
// pareRequestLine returns the RequestLine, bytes consumed, and optional error.
func parseRequestLine(request []byte) (RequestLine, int, error) {
	requestText := string(request)

	lines := strings.Split(requestText, "\r\n")
	if len(lines) < 2 {
		return RequestLine{}, 0, nil // No CRLF, so need to read more.
	}

	requestLineText := lines[0]

	requestLineParts := strings.Split(requestLineText, " ")
	if len(requestLineParts) != 3 {
		return RequestLine{}, 0, errors.New("Invalid request line.")
	}

	method := requestLineParts[0]
	if !isCapitalOnly(method) {
		return RequestLine{}, 0, fmt.Errorf(
			`Invalid request method: "%s". Method may only contain captial letters.`,
			method,
		)
	}

	requestTarget := requestLineParts[1]

	httpVersion := requestLineParts[2]
	httpVersionParts := strings.Split(httpVersion, "/")
	if len(httpVersionParts) != 2 {
		return RequestLine{}, 0, fmt.Errorf(
			`Invalid HTTP version: "%s". Required format: HTTP-name "/" DIGIT "." DIGIT`,
			httpVersion,
		)
	}
	if httpVersionParts[0] != "HTTP" {
		return RequestLine{}, 0, fmt.Errorf(
			`Invalid HTTP version name: "%s". Only HTTP/1.1 is supported.`,
			httpVersionParts[0],
		)
	}
	if httpVersionParts[1] != "1.1" {
		return RequestLine{}, 0, fmt.Errorf(
			`Invalid HTTP version number: "%s". Only HTTP/1.1 is supported.`,
			httpVersionParts[1],
		)
	}

	return RequestLine{
			Method:        method,
			RequestTarget: requestTarget,
			HttpVersion:   httpVersionParts[1],
		},
		len(requestLineText) + 2, // + 2 for CRLF
		nil
}

func isCapitalOnly(text string) bool {
	textBytes := []byte(text)

	for _, char := range textBytes {
		if char < 65 || char > 90 {
			return false
		}
	}
	return true
}
