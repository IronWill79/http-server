package request

import (
	"errors"
	"fmt"
	"io"
	"os"
	"slices"
	"strings"
)

const bufferSize = 8
const crlf = "\r\n"

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type RequestStatus int

const (
	Initialized RequestStatus = iota
	Done
)

var statusName = map[RequestStatus]string{
	Initialized: "initialized",
	Done:        "done",
}

func (rs RequestStatus) String() string {
	return statusName[rs]
}

type Request struct {
	RequestLine RequestLine
	Status      RequestStatus
}

func parseRequestLine(req_line *RequestLine, text string) (int, error) {
	if !strings.Contains(text, crlf) {
		return 0, nil
	}
	length := len(text)
	valid_methods := []string{"GET", "POST"}
	lines := strings.Split(text, crlf)
	parts := strings.Split(lines[0], " ")
	if len(parts) != 3 {
		err := errors.New("invalid request line - not 3 parts")
		fmt.Fprintf(os.Stderr, "%v: %v\n", err, parts)
		return length, err
	}
	method := parts[0]
	if method != strings.ToUpper(method) || !slices.Contains(valid_methods, method) {
		err := errors.New("invalid method")
		fmt.Fprintf(os.Stderr, "%v: %v\n", err, method)
		return length, err
	}
	target := parts[1]
	http_line := parts[2]
	if http_line != "HTTP/1.1" {
		err := errors.New("invalid HTTP version")
		fmt.Fprintf(os.Stderr, "%v: %v\n", err, http_line)
		return length, err
	}
	http_version := strings.Split(http_line, "/")[1]
	req_line.Method = method
	req_line.RequestTarget = target
	req_line.HttpVersion = http_version
	return length, nil
}

func (r *Request) parse(data []byte) (int, error) {
	switch r.Status {
	case Initialized:
		length, err := parseRequestLine(&r.RequestLine, string(data))
		if err != nil {
			fmt.Printf("Request.parse: error parsing: %v\n", err)
			return 0, err
		} else if length == 0 {
			return 0, nil
		}
		r.Status = Done
		return length, nil
	case Done:
		return 0, errors.New("error: trying to read data in a done state")
	default:
		return 0, errors.New("error: unknown state")
	}
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	buf := make([]byte, bufferSize)
	readToIndex := 0
	request := &Request{Status: Initialized}
	for request.Status != Done {
		if readToIndex == len(buf) {
			new_buf := make([]byte, len(buf)*2)
			_ = copy(new_buf, buf)
			buf = new_buf
		}

		// read into the buffer
		numBytesRead, err := reader.Read(buf[readToIndex:])
		if err != nil {
			if err != io.EOF {
				return nil, err
			} else {
				request.Status = Done
				break
			}
		}
		readToIndex += numBytesRead

		// parse from the buffer
		numBytesParsed, err := request.parse(buf)
		if err != nil {
			return nil, err
		}
		if numBytesParsed > 0 {
			new_buf := make([]byte, len(buf))
			_ = copy(new_buf, buf[numBytesParsed:])
			buf = new_buf
			readToIndex -= numBytesParsed
		}
	}
	return request, nil
}
