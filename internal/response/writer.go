package response

import (
	"errors"
	"io"

	"github.com/kn1ghtm0nster/http-from-tcp/internal/headers"
)

type writerState int

const (
	WriterStateReady writerState = iota
	WriterStateHeaders
	WriterStateBody
)

type Writer struct {
	w 	io.Writer
	state 	writerState
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{
		w:     w,
		state: WriterStateReady,
	}
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	if w.state != WriterStateReady {
		return errors.New("cannot write status line in current state")
	}

	_, err := w.w.Write(getStatusLine(statusCode))
	if err != nil {
		return err
	}
	w.state = WriterStateHeaders
	return nil
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	if w.state != WriterStateHeaders {
		return errors.New("cannot write headers in current state")
	}

	for k, v := range headers {
		_, err := w.w.Write([]byte(k + ": " + v + "\r\n"))
		if err != nil {
			return err
		}
	}
	_, err := w.w.Write([]byte("\r\n"))
	if err != nil {
		return err
	}

	w.state = WriterStateBody
	return nil
}

func (w *Writer) WriteBody(body []byte) (int, error) {
	if w.state != WriterStateBody {
		return 0, errors.New("cannot write body in current state")
	}

	n, err := w.w.Write(body)
	if err != nil {
		return n, err
	}
	return n, nil
}