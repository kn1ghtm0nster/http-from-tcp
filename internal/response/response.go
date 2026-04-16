package response

import (
	"errors"
	"fmt"
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

func (w *Writer) WriteChunkedBody(p []byte) (int, error) {
	if w.state != WriterStateBody {
		return 0, errors.New("cannot write chunked body in current state")
	}

	hexNumber, err := fmt.Fprintf(w.w, "%x\r\n", len(p))
	if err != nil {
		return 0, err
	}

	pBytes, err := w.w.Write(p)
	if err != nil {
		return 0, err
	}

	endBytes, err := w.w.Write([]byte("\r\n"))
	if err != nil {
		return 0, err
	}

	return hexNumber + pBytes + endBytes, nil
	
}

func (w *Writer) WriteChunkedBodyDone() (int, error) {
	if w.state != WriterStateBody {
		return 0, errors.New("cannot write chunked body done in current state")
	}

	endHexNumber, err := w.w.Write([]byte("0\r\n"))
	if err != nil {
		return 0, err
	}
	w.state = WriterStateTrailers
	return endHexNumber, nil
}

func (w *Writer) WriteTrailers(h headers.Headers) error {
	if w.state != WriterStateTrailers {
		return errors.New("cannot write trailers in current state")
	}

	for k, v := range h {
		_, err := w.w.Write([]byte(k + ": " + v + "\r\n"))
		if err != nil {
			return err
		}
	}
	_, err := w.w.Write([]byte("\r\n"))
	if err != nil {
		return err
	}
	return nil
}