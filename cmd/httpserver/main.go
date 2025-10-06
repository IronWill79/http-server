package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/IronWill79/http-server/internal/headers"
	"github.com/IronWill79/http-server/internal/request"
	"github.com/IronWill79/http-server/internal/response"
	"github.com/IronWill79/http-server/internal/server"
)

const chunkedBufferSize = 1024

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
	case "/video":
		video, err := os.ReadFile("./assets/vim.mp4")
		if err != nil {
			w.WriteStatusLine(response.StatusCodeBadRequest)
			h := headers.GetDefaultHeaders(len(responseBadRequest))
			h.Set("content-type", "text/html")
			w.WriteHeaders(h)
			w.WriteBody([]byte(responseBadRequest))
		}
		w.WriteStatusLine(response.StatusCodeSuccess)
		h := headers.NewHeaders()
		h.Set("Content-Type", "video/mp4")
		w.WriteHeaders(h)
		w.WriteBody(video)
	default:
		if strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin/") {
			suffix := strings.TrimPrefix(req.RequestLine.RequestTarget, "/httpbin/")
			h := headers.NewHeaders()
			h.Set("Content-Type", "text/plain")
			h.Set("Transfer-Encoding", "chunked")
			h.Set("Trailer", "X-Content-Length, X-Content-SHA256")
			w.WriteStatusLine(response.StatusCodeSuccess)
			w.WriteHeaders(h)
			resp, err := http.Get(fmt.Sprintf("https://httpbin.org/%s", suffix))
			if err != nil {
				log.Fatalf("http request failed: %v", err)
			}
			var responseBody []byte
			buf := make([]byte, chunkedBufferSize)
			defer resp.Body.Close()
			for {
				n, err := resp.Body.Read(buf)
				if n > 0 {
					log.Printf("%d bytes read", n)
					bytesWritten, err := w.WriteChunkedBody(buf[:n])
					if err != nil {
						log.Fatalf("WriteChunkedBody failed: %v", err)
					}
					responseBody = append(responseBody, buf[:n]...)
					log.Printf("%d bytes written", bytesWritten)
				}
				if err == io.EOF {
					break
				}
				if err != nil {
					log.Fatalf("response read failed: %v", err)
				}
			}
			_, err = w.WriteChunkedBodyDone()
			if err != nil {
				log.Fatalf("WriteChunkedBodyDone failed: %v", err)
			}
			log.Printf("%d bytes written. ChunkedBody done", len(responseBody))
			t := headers.NewHeaders()
			t.Set("X-Content-Length", fmt.Sprintf("%d", len(responseBody)))
			t.Set("X-Content-SHA256", fmt.Sprintf("%x", sha256.Sum256(responseBody)))
			err = w.WriteTrailers(t)
			if err != nil {
				log.Fatalf("WriteTrailers failed: %v", err)
			}
		} else {
			w.WriteStatusLine(response.StatusCodeSuccess)
			h := headers.GetDefaultHeaders(len(responseSuccess))
			h.Set("content-type", "text/html")
			w.WriteHeaders(h)
			w.WriteBody([]byte(responseSuccess))
		}
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
