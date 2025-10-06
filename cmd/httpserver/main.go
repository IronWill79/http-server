package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/IronWill79/http-server/internal/headers"
	"github.com/IronWill79/http-server/internal/request"
	"github.com/IronWill79/http-server/internal/response"
	"github.com/IronWill79/http-server/internal/server"
)

const port = 42069

const responseBadRequest = `<html>
  <head>
    <title>400 Bad Request</title>
  </head>
  <body>
    <h1>Bad Request</h1>
    <p>Your request honestly kinda sucked.</p>
  </body>
</html>`

const responseServerError = `<html>
  <head>
    <title>500 Internal Server Error</title>
  </head>
  <body>
    <h1>Internal Server Error</h1>
    <p>Okay, you know what? This one is on me.</p>
  </body>
</html>`

const responseSuccess = `<html>
  <head>
    <title>200 OK</title>
  </head>
  <body>
    <h1>Success!</h1>
    <p>Your request was an absolute banger.</p>
  </body>
</html>`

func handler(w *response.Writer, req *request.Request) {
	switch req.RequestLine.RequestTarget {
	case "/yourproblem":
		w.WriteStatusLine(response.StatusCodeBadRequest)
		h := headers.GetDefaultHeaders(len(responseBadRequest))
		h.Set("content-type", "text/html")
		w.WriteHeaders(h)
		w.WriteBody([]byte(responseBadRequest))
	case "/myproblem":
		w.WriteStatusLine(response.StatusCodeServerError)
		h := headers.GetDefaultHeaders(len(responseServerError))
		h.Set("content-type", "text/html")
		w.WriteHeaders(h)
		w.WriteBody([]byte(responseServerError))
	default:
		w.WriteStatusLine(response.StatusCodeSuccess)
		h := headers.GetDefaultHeaders(len(responseSuccess))
		h.Set("content-type", "text/html")
		w.WriteHeaders(h)
		w.WriteBody([]byte(responseSuccess))
	}
}

func main() {
	server, err := server.Serve(port, handler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}
