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

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	_, err := w.writer.Write([]byte(statusCode.String()))
	if err != nil {
		return err
	}
	return nil
}

func (w *Writer) WriteHeaders(h headers.Headers) error {
	for k, v := range h {
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

func (w *Writer) WriteBody(p []byte) (int, error) {
	bytesWritten, err := w.writer.Write(p)
	if err != nil {
		return 0, err
	}
	return bytesWritten, nil
}

func (w *Writer) WriteChunkedBody(p []byte) (int, error) {
	hexLength := fmt.Sprintf("%x\r\n", len(p))
	hexBytesWritten, err := w.writer.Write([]byte(hexLength))
	if err != nil {
		return 0, err
	}
	bytesWritten, err := w.writer.Write(p)
	if err != nil {
		return hexBytesWritten, err
	}
	crlfWritten, err := w.writer.Write([]byte{13, 10})
	if err != nil {
		return hexBytesWritten + bytesWritten, err
	}
	return hexBytesWritten + bytesWritten + crlfWritten, nil
}

func (w *Writer) WriteChunkedBodyDone() (int, error) {
	bytesWritten, err := w.WriteChunkedBody([]byte{13, 10})
	if err != nil {
		return 0, err
	}
	return bytesWritten, nil
}

func (w *Writer) WriteTrailers(h headers.Headers) error {
	fmt.Fprint(w.writer, "0\r\n")
	for k, v := range h {
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
