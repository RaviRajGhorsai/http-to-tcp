package request

import (
	"bytes"
	"fmt"
	"io"
)

type parserState string

const (
	StateInit parserState = "init"
	StateDone parserState = "done"
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type Request struct {
	RequestLine RequestLine
	state       parserState
}

var ERROR_BAD_START_LINE = fmt.Errorf("bad start line")
var SEPARATOR = []byte("\r\n")

func newRequest() *Request {

	return &Request{
		state: StateInit,
	}

}

func validateStartLine(data [][]byte) error {

	// validate http part

	httpParts := bytes.Split(data[2], []byte("/"))

	if len(httpParts) != 2 || !bytes.Equal(httpParts[0], []byte("HTTP")) || !bytes.Equal(httpParts[1], []byte("1.1")) {
		return ERROR_BAD_START_LINE
	}

	// validate method part
	for _, b := range data[0] {
		if b < 'A' || b > 'Z' {
			return ERROR_BAD_START_LINE
		}
	}

	// validate target part
	requestTarget := data[1]

	if len(requestTarget) == 0 || requestTarget[0] != '/' {
		return ERROR_BAD_START_LINE
	}

	return nil

}

func parseRequestLine(b []byte) (*RequestLine, int, error) {

	// finds the index where the SEPARATOR is present
	index := bytes.Index(b, SEPARATOR)
	if index == -1 {
		return nil, 0, nil
	}

	startLine := b[:index]
	read := index + len(SEPARATOR)

	// again speparate startline by space separator
	parts := bytes.Split(startLine, []byte(" "))

	if len(parts) != 3 {

		return nil, 0, ERROR_BAD_START_LINE

	}

	err := validateStartLine(parts)

	if err != nil {
		return nil, 0, err

	}

	httppart := bytes.Split(parts[2], []byte("/"))

	return &RequestLine{
		Method:        string(parts[0]),
		RequestTarget: string(parts[1]),
		HttpVersion:   string(httppart[1]),
	}, read, nil
}

func (r *Request) parse(data []byte) (int, error) {

	read := 0
outer:
	for {
		switch r.state {
		case StateInit:

			// we parse the mesasge that we got and return as requestLine, no. of byte read
			requestLine, n, err := parseRequestLine(data[read:])
			if err != nil {
				return 0, err
			}

			if n == 0 {
				break outer
			}
			
			// we save the successfully read line
			r.RequestLine = *requestLine
			
			// then increase read buffer/size
			read += n
			
			r.state = StateDone

		case StateDone:
			break outer

		}
	}
	return read, nil
}

func (r *Request) done() bool {
	return r.state == StateDone

}

func RequestFromReader(reader io.Reader) (*Request, error) {

	request := newRequest()
	
	// initialize msg byte length as 0
	buf := make([]byte, 1024)
	bufLen := 0

	for !request.done() {
		
		// n == total byte size of message
		n, err := reader.Read(buf[bufLen:])
		if err != nil {

			return nil, err

		}

		// set buffer length as total message length
		bufLen += n
		
		// eandN is the no. of bytes read by parser
		readN, err := request.parse(buf[:bufLen])
		if err != nil {
			return nil, err

		}

		// this copies the rest not read message and ommits the read message
		copy(buf, buf[readN:bufLen])
		bufLen -= readN

	}

	return request, nil
}
