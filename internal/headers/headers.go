package headers

import (
	"bytes"
	"fmt"
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

	if data[nameColon-1] == ' ' {
		return 0, false, fmt.Errorf("No whitespace is allowed between the field name and colon.")
	}

	fieldName := bytes.TrimSpace(data[0:nameColon])
	fieldValue := bytes.TrimSpace(data[nameColon+1 : crlf])

	h[string(fieldName)] = string(fieldValue)
	return crlf + 2, false, nil
}
