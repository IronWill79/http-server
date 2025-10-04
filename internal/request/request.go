package request

import (
	"errors"
	"fmt"
	"io"
	"slices"
	"strconv"
	"strings"

	"github.com/IronWill79/http-server/internal/headers"
)

const bufferSize = 8
const crlf = "\r\n"

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type RequestState int

const (
	requestStateInitialized RequestState = iota
	requestStateParsingHeaders
	requestStateParsingBody
	requestStateDone
)

var statusName = map[RequestState]string{
	requestStateInitialized:    "initialized",
	requestStateParsingHeaders: "parsing-headers",
	requestStateParsingBody:    "parsing-body",
	requestStateDone:           "done",
}

func (rs RequestState) String() string {
	return statusName[rs]
}

type Request struct {
	Body        []byte
	Headers     headers.Headers
	RequestLine RequestLine
	state       RequestState
}

func (r *Request) parseRequestLine(data []byte) (int, error) {
	text := string(data)
	if !strings.Contains(text, crlf) {
		return 0, nil
	}
	length := strings.Index(text, crlf) + 2
	valid_methods := []string{"GET", "POST"}
	lines := strings.Split(text, crlf)
	parts := strings.Split(lines[0], " ")
	if len(parts) != 3 {
		err := errors.New("invalid request line - not 3 parts")
		return length, err
	}
	method := parts[0]
	if method != strings.ToUpper(method) || !slices.Contains(valid_methods, method) {
		err := errors.New("invalid method")
		return length, err
	}
	target := parts[1]
	http_line := parts[2]
	if http_line != "HTTP/1.1" {
		err := errors.New("invalid HTTP version")
		return length, err
	}
	http_version := strings.Split(http_line, "/")[1]
	r.RequestLine.Method = method
	r.RequestLine.RequestTarget = target
	r.RequestLine.HttpVersion = http_version
	return length, nil
}

func (r *Request) parseSingle(data []byte) (int, error) {
	switch r.state {
	case requestStateInitialized:
		length, err := r.parseRequestLine(data)
		if err != nil {
			return 0, err
		}
		if length == 0 {
			return 0, nil
		}
		r.state = requestStateParsingHeaders
		return length, nil
	case requestStateParsingHeaders:
		length, done, err := r.Headers.Parse(data)
		if err != nil {
			return 0, err
		}
		if length == 0 {
			return 0, nil
		}
		if done {
			r.state = requestStateParsingBody
		}
		return length, nil
	case requestStateParsingBody:
		contentLengthString, err := r.Headers.Get("content-length")
		if err != nil {
			r.state = requestStateDone
			return 0, nil
		}
		contentLength, err := strconv.Atoi(contentLengthString)
		if err != nil {
			return 0, fmt.Errorf("error: content-length not a number %s", contentLengthString)
		}
		r.Body = append(r.Body, data...)
		if len(r.Body) > contentLength {
			return 0, fmt.Errorf("error: body (%d) longer than content-length (%d)\nBody: `%v`", len(string(r.Body)), contentLength, r.Body)
		}
		if len(r.Body) == contentLength {
			r.state = requestStateDone
		}
		return len(data), nil
	case requestStateDone:
		return 0, errors.New("error: trying to read data in a done state")
	default:
		return 0, errors.New("error: unknown state")
	}
}

func (r *Request) parse(data []byte) (int, error) {
	totalBytesParsed := 0
	for r.state != requestStateDone {
		n, err := r.parseSingle(data[totalBytesParsed:])
		if err != nil {
			return 0, err
		} else if n == 0 {
			return totalBytesParsed, nil
		}
		totalBytesParsed += n
	}
	return totalBytesParsed, nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	buf := make([]byte, bufferSize)
	readToIndex := 0
	request := &Request{state: requestStateInitialized, Headers: headers.NewHeaders()}
	for request.state != requestStateDone {
		if readToIndex == len(buf) {
			new_buf := make([]byte, len(buf)*2)
			_ = copy(new_buf, buf)
			buf = new_buf
		}

		// read into the buffer
		numBytesRead, err := reader.Read(buf[readToIndex:])
		if err != nil {
			if errors.Is(err, io.EOF) {
				request.state = requestStateDone
				break
			}
			return nil, err
		}
		readToIndex += numBytesRead

		// parse from the buffer
		numBytesParsed, err := request.parse(buf[:readToIndex])
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
	contentLengthString, err := request.Headers.Get("content-length")
	if err != nil {
		return request, nil
	}
	contentLength, err := strconv.Atoi(contentLengthString)
	if err != nil {
		return nil, fmt.Errorf("error: content-length not a number %s", contentLengthString)
	}
	if len(request.Body) < contentLength {
		return nil, fmt.Errorf("body shorter than content-length")
	}
	return request, nil
}
