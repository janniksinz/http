package request

import (
	"bytes"
	"fmt"
	"io"
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type parserState string

const (
	StateInit  parserState = "init"
	StateDone  parserState = "done"
	StateError parserState = "error"
)

type Request struct {
	RequestLine RequestLine
	state       parserState
}

func newRequest() *Request {
	return &Request{
		state: StateInit,
	}
}

var ERROR_MALFORMED_REQUEST_LINE = fmt.Errorf("malformed request-line")
var ERROR_UNSUPPORTED_HTTP_VERSION = fmt.Errorf("unsupported HTTP version")
var ERROR_REQUEST_IN_ERROR_STATE = fmt.Errorf("request in error state")
var SEPARATOR = []byte("\r\n")

func RequestFromReader(reader io.Reader) (*Request, error) {
	request := newRequest()

	// NOTE: buffer could be too small
	buff := make([]byte, 1024)
	bufLen := 0
	for !request.done() {
		n, err := reader.Read(buff[bufLen:])
		// TODO: what to do with error?
		if err != nil {
			return nil, err
		}

		// parse buffer
		bufLen += n // discard already parsed data
		readN, err := request.parse(buff[:bufLen])
		if err != nil {
			return nil, err
		}

		// move remaining data to the beginning
		copy(buff, buff[readN:bufLen])
		bufLen -= readN

	}

	return request, nil

}

/*
returns:

	RequestLine - the fully parsed RequestLine
	int - the amounts of bytes already parsed
*/
func parseRequestLine(b []byte) (*RequestLine, int, error) {
	idx := bytes.Index(b, SEPARATOR)
	if idx == -1 {
		return nil, 0, nil
	}

	startLine := b[:idx]
	read := idx + len(SEPARATOR)

	parts := bytes.Split(startLine, []byte(" "))
	if len(parts) != 3 {
		return nil, 0, ERROR_MALFORMED_REQUEST_LINE
	}

	httpParts := bytes.Split(parts[2], []byte("/"))
	if len(httpParts) != 2 || string(httpParts[0]) != "HTTP" || string(httpParts[1]) != "1.1" {
		return nil, 0, ERROR_MALFORMED_REQUEST_LINE
	}

	rl := &RequestLine{
		Method:        string(parts[0]),
		RequestTarget: string(parts[1]),
		HttpVersion:   string(httpParts[1]),
	}

	return rl, read, nil
}

/*
returns:

	int - how much of the buffer we parsed
	error
*/
func (r *Request) parse(data []byte) (int, error) {
	read := 0
outer:
	for {
		switch r.state {
		case StateError:
			return 0, ERROR_REQUEST_IN_ERROR_STATE
		case StateInit:
			// parse data starting from read index
			rl, n, err := parseRequestLine(data[read:])
			if err != nil {
				r.state = StateError
				return 0, err
			}
			if n == 0 {
				break outer
			}
			r.RequestLine = *rl
			read += n
			r.state = StateDone
		case StateDone:
			break outer
		}
	}
	return 0, nil
}

func (r *Request) done() bool {
	return r.state == StateDone || r.state == StateError
}
