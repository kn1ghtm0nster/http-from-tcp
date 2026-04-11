package request

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"unicode"

	"github.com/kn1ghtm0nster/http-from-tcp/internal/headers"
)

const bufferSize = 8

const (
	requestStateInitialized requestState = iota
	requestStateDone
	requestStateParsingHeaders
)

type requestState int

type Request struct {
	RequestLine RequestLine
	Headers 	headers.Headers
	state       requestState
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	buf := make([]byte, bufferSize)
	readToIndex := 0
	req := &Request{
		state: requestStateInitialized,
		Headers: headers.NewHeaders(),
	}

	for req.state != requestStateDone {
		// if buffer is full, create buffer twice the size and copy old ddata into new slice
		if readToIndex == len(buf) {
			newBuf := make([]byte, len(buf)*2)
			copy(newBuf, buf)
			buf = newBuf
		}
		// read from the reader into the buffer starting at the current read index
		nextIdx, err := reader.Read(buf[readToIndex:])
		if err != nil {
			if err == io.EOF && req.state == requestStateDone {
				return req, nil
			}
			if err == io.EOF && req.state != requestStateDone {
				return nil, errors.New("incomplete request")
			}
			return nil, err
		}
		// update the index of how many bytes have been read into the buffer
		// and then attempt to parse the data that has been read so far
		readToIndex += nextIdx
		parsedBytes, err := req.parse(buf[:readToIndex])
		if err != nil {
			return nil, err
		}
		// remove the data that was parsed from the buffer and copy to a new slice
		copy(buf, buf[parsedBytes:readToIndex])
		readToIndex -= parsedBytes
	}
	return req, nil
}

func parseRequestLine(data []byte) (RequestLine, int, error) {
	idx := bytes.Index(data, []byte("\r\n"))
	// discard everything after the first line
	if idx == -1 {
		return RequestLine{}, 0, nil
	}

	// separate out the request-line into parts by spaces
	requestLineParts := strings.Split(string(data[:idx]), " ")
	if len(requestLineParts) != 3 {
		return RequestLine{}, 0, errors.New("invalid request format")
	}

	method := requestLineParts[0]
	if method == "" {
		return RequestLine{}, 0, errors.New("no method specified in request")
	}
	// ensure that method is composed entirely of uppercase letters
	for _, c := range method {
		if !unicode.IsUpper(c) {
			return RequestLine{}, 0, errors.New("invalid method format")
		}
	}

	// ensure that the HTTP version is strictly HTTP/1.1
	httpVersion := requestLineParts[2]
	if httpVersion != "HTTP/1.1" {
		return RequestLine{}, 0, errors.New("unsupported HTTP version")
	}

	return RequestLine{
		Method:        requestLineParts[0],
		RequestTarget: requestLineParts[1],
		HttpVersion:   strings.TrimPrefix(requestLineParts[2], "HTTP/"),
	}, idx + 2, nil
}

func (r *Request) parse(data []byte) (int, error) {
	totalBytesParsed := 0

	// continue parsing the headers until the request is completely parsed
	for r.state != requestStateDone {
		n, err := r.parseSingle(data[totalBytesParsed:])
		if err != nil {
			return totalBytesParsed, err
		}
		if n == 0 {
			break
		}
		totalBytesParsed += n
	}
	return totalBytesParsed, nil
}

func (r *Request) parseSingle(data []byte) (int, error) {
	switch r.state {
		case requestStateDone:
			return 0, errors.New("error: trying to read data in a done state")
		case requestStateInitialized:
			requestLine, bytesRead, err := parseRequestLine(data)
			if err != nil {
				return 0, err
			}
			if bytesRead == 0 {
				return 0, nil
			}
			r.RequestLine = requestLine
			r.state = requestStateParsingHeaders
			return bytesRead, nil
		case requestStateParsingHeaders:
			n, done, err := r.Headers.Parse(data)
			if err != nil {
				return 0, err
			}
			if done {
				r.state = requestStateDone
				return n, nil
			}
			return n, nil
		default:
			return 0, errors.New("error: unknown state")
	}
}