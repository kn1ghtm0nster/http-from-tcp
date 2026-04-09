package request

import (
	"errors"
	"io"
	"strings"
	"unicode"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	byteData, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	// parsed data here via custom parseRequestLine function
	requestLine, err := parseRequestLine(byteData)
	if err != nil {
		return nil, err
	}

	return &Request{
		RequestLine: requestLine,
	}, nil
}

func parseRequestLine(data []byte) (RequestLine, error) {
	parts := strings.Split(string(data), "\r\n")
	// discard everything after the first line
	if len(parts) == 0 {
		return RequestLine{}, errors.New("empty request")
	}

	// separate out the request-line into parts by spaces
	requestLineParts := strings.Split(parts[0], " ")
	if len(requestLineParts) != 3 {
		return RequestLine{}, errors.New("invalid request format")
	}

	method := requestLineParts[0]
	if method == "" {
		return RequestLine{}, errors.New("no method specified in request")
	}
	// ensure that method is composed entirely of uppercase letters
	for _, c := range method {
		if !unicode.IsUpper(c) {
			return RequestLine{}, errors.New("invalid method format")
		}
	}

	// ensure that the HTTP version is strictly HTTP/1.1
	httpVersion := requestLineParts[2]
	if httpVersion != "HTTP/1.1" {
		return RequestLine{}, errors.New("unsupported HTTP version")
	}

	return RequestLine{
		Method:        requestLineParts[0],
		RequestTarget: requestLineParts[1],
		HttpVersion:   strings.TrimPrefix(requestLineParts[2], "HTTP/"),
	}, nil
}