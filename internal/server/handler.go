package server

import (
	"io"

	"github.com/IronWill79/http-server/internal/headers"
	"github.com/IronWill79/http-server/internal/request"
	"github.com/IronWill79/http-server/internal/response"
)

type Handler func(w *response.Writer, req *request.Request)

type HandlerError struct {
	Status  response.StatusCode
	Message string
}

func (e *HandlerError) Write(w io.Writer) {
	writer := response.NewWriter(w)
	writer.WriteStatusLine(e.Status)
	h := headers.GetDefaultHeaders(len(e.Message))
	writer.WriteHeaders(h)
	writer.WriteBody([]byte(e.Message))
}
