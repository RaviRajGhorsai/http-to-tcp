package request

import (
	"fmt"
	"io"
	"strings"
	"unicode"
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type Request struct {
	RequestLine RequestLine
}

var ERROR_BAD_START_LINE = fmt.Errorf("bad start line")
var SEPARATOR = "\r\n"

func validateStartLine(data []string) ([]string, error) {

	// validate http part
	httpParts := strings.Split(data[2], "/")

	if len(httpParts) != 2 || httpParts[0] != "HTTP" || httpParts[1] != "1.1" {
		return nil, ERROR_BAD_START_LINE
	}

	// validate method part
	for _, r := range data[0] {
		if !unicode.IsUpper(r) {
			return nil, ERROR_BAD_START_LINE
		}

	}

	// validate target part
	requestTarget := data[1]

	if len(requestTarget) == 0 || requestTarget[0] != '/' {
		return nil, ERROR_BAD_START_LINE
	}

	return data, nil

}

func parseRequestLine(b string) (*RequestLine, string, error) {

	// finds the index where the SEPARATOR is present
	before, after, ok := strings.Cut(b, SEPARATOR)
	if !ok {

		return nil, b, nil
	}

	// save the starting line (GET / HTTP/1.1)
	startLine := before
	restOfMsg := after

	// again speparate startline by space separator
	parts := strings.Split(startLine, " ")

	if len(parts) != 3 {

		return nil, restOfMsg, ERROR_BAD_START_LINE

	}

	var err error
	parts, err = validateStartLine(parts)

	if err != nil {
		return nil, restOfMsg, err

	}

	httppart := strings.Split(parts[2], "/")

	return &RequestLine{
		Method:        parts[0],
		RequestTarget: parts[1],
		HttpVersion:   httppart[1],
	}, restOfMsg, nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {

	data, err := io.ReadAll(reader)

	if err != nil {

		return nil, fmt.Errorf("unable to io.ReadAll: %w", err)

	}

	str := string(data)
	requestLine, _, err := parseRequestLine(str)
	if err != nil {
		return nil, err
	}
	if requestLine == nil {
		return nil, fmt.Errorf("incomplete request line")
	}

	return &Request{
		RequestLine: *requestLine,
	}, err
}
