package response

import (
	"strconv"

	"github.com/kn1ghtm0nster/http-from-tcp/internal/headers"
)

type StatusCode int

const (
	StatusOK StatusCode = 200
	BadRequest StatusCode = 400
	InternalServerError StatusCode = 500
)

func getStatusLine(statusCode StatusCode) []byte {
	// maps the given status code to the correct reason phrase
	var reasonPhrase string
	switch statusCode {
	case StatusOK:
		reasonPhrase = "OK"
	case BadRequest:
		reasonPhrase = "Bad Request"
	case InternalServerError:
		reasonPhrase = "Internal Server Error"
	default:
		reasonPhrase = ""
	}

	return []byte("HTTP/1.1 " + strconv.Itoa(int(statusCode)) + " " + reasonPhrase + "\r\n")
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	return headers.Headers{
		"Content-Length": strconv.Itoa(contentLen),
		"Connection": "close",
		"Content-Type": "text/plain",
	}
}