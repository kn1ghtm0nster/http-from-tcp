package headers

import (
	"errors"
	"strings"
)

type Headers map[string]string

func NewHeaders() Headers {
	return make(Headers)
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	parsedData := string(data)

	idx := strings.Index(parsedData, "\r\n")
	if idx == -1 {
		return 0, false, nil
	}
	if idx == 0 {
		return idx + 2, true, nil
	}

	headerLine := parsedData[:idx]
	colonIdx := strings.Index(headerLine, ":")
	if colonIdx == -1 {
		return 0, false, nil
	}

	key := headerLine[:colonIdx]
	value := strings.TrimSpace(headerLine[colonIdx+1:])

	// key must not have leading or trailing whitespace
	if key != strings.TrimSpace(key) {
		return 0, false, errors.New("invalid key format")
	}

	// Key cannot contain special characters
	if strings.ContainsAny(key, " !@#$%^&*()[]{}<>?/\\|`~;\"'") {
		return 0, false, errors.New("invalid character in header key")
	}

	// normalize the header key to lowercase
	key = strings.ToLower(key)
	// if the header key already exists, append the new value to the existing value separated by a comma
	if existingValue, ok := h[key]; ok {
		h[key] = existingValue + ", " + value
	} else {
		h[key] = value
	}
	return idx + 2, false, nil
}

func (h Headers) Get(key string) (string, bool) {
	value, ok := h[strings.ToLower(key)]
	return value, ok
}

func (h Headers) Override(key, value string) {
	h[key] = value
}