package response

import (
	"fmt"
	"io"
	"strconv"

	"github.com/IronWill79/http-server/internal/headers"
)

type StatusCode int

const (
	StatusCodeSuccess     StatusCode = 200
	StatusCodeBadRequest  StatusCode = 400
	StatusCodeServerError StatusCode = 500
)

var statusCodeName = map[StatusCode]string{
	StatusCodeSuccess:     "HTTP/1.1 200 OK\r\n",
	StatusCodeBadRequest:  "HTTP/1.1 400 Bad Request\r\n",
	StatusCodeServerError: "HTTP/1.1 500 Internal Server Error\r\n",
}

func (code StatusCode) String() string {
	return statusCodeName[code]
}

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	_, err := w.Write([]byte(statusCode.String()))
	if err != nil {
		return err
	}
	return nil
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	h := headers.NewHeaders()
	h["content-length"] = strconv.Itoa(contentLen)
	h["connection"] = "close"
	h["content-type"] = "text/plain"
	return h
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	for k, v := range headers {
		_, err := fmt.Fprintf(w, "%s: %s\r\n", k, v)
		if err != nil {
			return err
		}
	}
	_, err := w.Write([]byte("\r\n"))
	if err != nil {
		return err
	}
	return nil
}

func Write(w io.Writer, body []byte) error {
	_, err := w.Write(body)
	if err != nil {
		return err
	}
	return nil
}
