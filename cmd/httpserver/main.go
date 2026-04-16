package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/kn1ghtm0nster/http-from-tcp/internal/headers"
	"github.com/kn1ghtm0nster/http-from-tcp/internal/request"
	"github.com/kn1ghtm0nster/http-from-tcp/internal/response"
	"github.com/kn1ghtm0nster/http-from-tcp/internal/server"
)

const port = 42069

const responseBody400 = `<html>
  <head>
    <title>400 Bad Request</title>
  </head>
  <body>
    <h1>Bad Request</h1>
    <p>Your request honestly kinda sucked.</p>
  </body>
</html>
`
const responseBody500 = `<html>
  <head>
    <title>500 Internal Server Error</title>
  </head>
  <body>
    <h1>Internal Server Error</h1>
    <p>Okay, you know what? This one is on me.</p>
  </body>
</html>
`

const responseBody200 = `<html>
  <head>
    <title>200 OK</title>
  </head>
  <body>
    <h1>Success!</h1>
    <p>Your request was an absolute banger.</p>
  </body>
</html>
`

var handler = func(r *response.Writer, req *request.Request) {
	if strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin/") {
		httpBinPath := strings.TrimPrefix(req.RequestLine.RequestTarget, "/httpbin/")
		res, err := http.Get("https://httpbin.org/" + httpBinPath)
		if err != nil {
			log.Printf("Error fetching from httpbin: %v", err)
			res.StatusCode = http.StatusBadRequest
			return
		}
		defer res.Body.Close()
		r.WriteStatusLine(response.StatusOK)
		headers := headers.NewHeaders()
		headers["Content-Type"] = res.Header.Get("Content-Type")
		headers["Transfer-Encoding"] = "chunked"
		r.WriteHeaders(headers)
		readBuffer := make([]byte, 1024)

		for {
			num, err := res.Body.Read(readBuffer)
			if num > 0 {
				_, writerErr := r.WriteChunkedBody(readBuffer[:num])
				if writerErr != nil {
					log.Printf("Error writing chunked body: %v", writerErr)
					return
				}
			}
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Printf("Error reading from httpbin response body: %v", err)
				return
			}
		}

		_, err = r.WriteChunkedBodyDone()
		if err != nil {
			log.Printf("Error finishing chunked body: %v", err)
		}
		return
	}

	switch req.RequestLine.RequestTarget {
		case "/yourproblem":
			responseBody := responseBody400
			r.WriteStatusLine(response.BadRequest)
			responseHeaders := response.GetDefaultHeaders(len(responseBody))
			responseHeaders.Override("Content-Type", "text/html")
			r.WriteHeaders(responseHeaders)
			r.WriteBody([]byte(responseBody))
		case "/myproblem":
			responseBody := responseBody500
			r.WriteStatusLine(response.InternalServerError)
			responseHeaders := response.GetDefaultHeaders(len(responseBody))
			responseHeaders.Override("Content-Type", "text/html")
			r.WriteHeaders(responseHeaders)
			r.WriteBody([]byte(responseBody))
		default:
			responseBody := responseBody200
			r.WriteStatusLine(response.StatusOK)
			responseHeaders := response.GetDefaultHeaders(len(responseBody))
			responseHeaders.Override("Content-Type", "text/html")
			r.WriteHeaders(responseHeaders)
			r.WriteBody([]byte(responseBody))
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