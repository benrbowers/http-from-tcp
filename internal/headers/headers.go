package headers

import (
	"bytes"
	"fmt"
	"slices"
	"strings"
)

type Headers map[string]string

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	crlf := bytes.Index(data, []byte{'\r', '\n'})
	if crlf == -1 {
		// Not enough data to parse
		return 0, false, nil
	}
	if crlf == 0 {
		// End of field lines
		return 2, true, nil
	}

	nameColon := bytes.Index(data, []byte{':'})

	if data[nameColon-1] == ' ' || data[nameColon-1] == '\t' {
		return 0, false, fmt.Errorf("No whitespace is allowed between the field name and colon.")
	}

	fieldName := bytes.TrimSpace(data[0:nameColon])
	fieldValue := bytes.TrimSpace(data[nameColon+1 : crlf])

	if !isValidFieldName(fieldName) {
		return 0, false, fmt.Errorf("Field name contains invalid characters: %s", fieldName)
	}

	h.Set(string(fieldName), string(fieldValue))

	return crlf + 2, false, nil
}

func (h Headers) Set(key, value string) {
	key = strings.ToLower(key)
	currentVal, alreadyExists := h[key]
	if alreadyExists {
		h[key] = currentVal + ", " + value
	} else {
		h[key] = value
	}
}

func (h Headers) Get(key string) (string, bool) {
	key = strings.ToLower(key)
	value, ok := h[key]
	return value, ok
}

func (h Headers) Replace(key, value string) {
	key = strings.ToLower(key)
	h[key] = value
}

func (h Headers) Remove(key string) {
	key = strings.ToLower(key)
	delete(h, key)
}

func (h Headers) Merge(h2 Headers) {
	for key, val := range h2 {
		h.Set(key, val)
	}
}

var tokenChars = []byte{'!', '#', '$', '%', '&', '\'', '*', '+', '-', '.', '^', '_', '`', '|', '~'}

func isValidFieldName(name []byte) bool {
	if len(name) < 1 {
		return false
	}

	for _, char := range name {
		if char >= '0' && char <= '9' {
			continue
		}
		if char >= 'A' && char <= 'Z' {
			continue
		}
		if char >= 'a' && char <= 'z' {
			continue
		}
		if slices.Contains(tokenChars, char) {
			continue
		}
		return false
	}

	return true
}
