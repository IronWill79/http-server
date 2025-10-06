package response

import (
	"fmt"
	"io"

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

type Writer struct {
	writer io.Writer
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{
		writer: w,
	}
}

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	_, err := w.Write([]byte(statusCode.String()))
	if err != nil {
		return err
	}
	return nil
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	_, err := w.writer.Write([]byte(statusCode.String()))
	if err != nil {
		return err
	}
	return nil
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

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	for k, v := range headers {
		_, err := fmt.Fprintf(w.writer, "%s: %s\r\n", k, v)
		if err != nil {
			return err
		}
	}
	_, err := w.writer.Write([]byte("\r\n"))
	if err != nil {
		return err
	}
	return nil
}

func WriteBody(w io.Writer, body []byte) error {
	_, err := w.Write(body)
	if err != nil {
		return err
	}
	return nil
}
func (w *Writer) WriteBody(body []byte) (int, error) {
	bytesWritten, err := w.writer.Write(body)
	if err != nil {
		return 0, err
	}
	return bytesWritten, nil
}
