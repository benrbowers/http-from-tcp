package request

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	requestBytes, err := io.ReadAll(reader)
	if err != nil {
		return &Request{}, err
	}

	requestLine, err := parseRequestLine(requestBytes)
	if err != nil {
		return &Request{}, err
	}

	return &Request{
		RequestLine: requestLine,
	}, nil
}

func parseRequestLine(request []byte) (RequestLine, error) {
	requestText := string(request)

	lines := strings.Split(requestText, "\r\n")
	if len(lines) < 2 {
		return RequestLine{}, errors.New("Invalid request format.")
	}

	requestLineText := lines[0]

	requestLineParts := strings.Split(requestLineText, " ")
	if len(requestLineParts) != 3 {
		return RequestLine{}, errors.New("Invalid request line.")
	}

	method := requestLineParts[0]
	if !isCapitalOnly(method) {
		return RequestLine{}, fmt.Errorf(
			`Invalid request method: "%s". Method may only contain captial letters.`,
			method,
		)
	}

	requestTarget := requestLineParts[1]

	httpVersion := requestLineParts[2]
	httpVersionParts := strings.Split(httpVersion, "/")
	if len(httpVersionParts) != 2 {
		return RequestLine{}, fmt.Errorf(
			`Invalid HTTP version: "%s". Required format: HTTP-name "/" DIGIT "." DIGIT`,
			httpVersion,
		)
	}
	if httpVersionParts[0] != "HTTP" {
		return RequestLine{}, fmt.Errorf(
			`Invalid HTTP version name: "%s". Only HTTP/1.1 is supported.`,
			httpVersionParts[0],
		)
	}
	if httpVersionParts[1] != "1.1" {
		return RequestLine{}, fmt.Errorf(
			`Invalid HTTP version number: "%s". Only HTTP/1.1 is supported.`,
			httpVersionParts[1],
		)
	}

	return RequestLine{
		Method:        method,
		RequestTarget: requestTarget,
		HttpVersion:   httpVersionParts[1],
	}, nil
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
