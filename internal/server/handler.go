package server

import (
	"io"

	"github.com/IronWill79/http-server/internal/request"
	"github.com/IronWill79/http-server/internal/response"
)

type Handler func(w io.Writer, req *request.Request) *HandlerError

type HandlerError struct {
	Status  response.StatusCode
	Message string
}

func (e *HandlerError) Write(w io.Writer) {
	response.WriteStatusLine(w, e.Status)
	h := response.GetDefaultHeaders(len(e.Message))
	response.WriteHeaders(w, h)
	response.Write(w, []byte(e.Message))
}
