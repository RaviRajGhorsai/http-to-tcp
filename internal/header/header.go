package headers

import (
	"bytes"
	"fmt"
)

type Headers map[string]string

var crlf = []byte("\r\n")

func NewHeaders() Headers {
	return map[string]string{}
}

func parseHeaderLine(fieldLine []byte) (string, string, error) {

	// split name and value field separated by colon
	parts := bytes.SplitN(fieldLine, []byte(":"), 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("Invalid Header")
	}

	name := parts[0]
	value := bytes.TrimSpace(parts[1]) // Trims all white space in value field

	// check if name field has whitespace if yes it is invalid
	if bytes.HasSuffix(name, []byte(" ")) {

		return "", "", fmt.Errorf("Invalid Header (invalid name field)\n")
	}

	return string(name), string(value), nil
}

func (h Headers) Parse(data []byte) (int, bool, error) {

	read := 0

	done := false

	for {

		index := bytes.Index(data[read:], crlf)
		if index == -1 {

			break
		}

		// empty header
		if index == 0 {

			done = true
			read += len(crlf)
			break
		}

		name, value, err := parseHeaderLine(data[read : read+index])
		if err != nil {

			return 0, false, err
		}

		read += index + len(crlf)

		h[name] = value

		

	}

	return read, done, nil

}
