package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHeadersParse(t *testing.T) {
	// Test: Valid single header
	headers := NewHeaders()
	data := []byte("Host: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 23, n)
	assert.False(t, done)

	// Test: Valid single header with extra white space
	headers = NewHeaders()
	data = []byte("Host:    localhost:42069   \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 29, n)
	assert.False(t, done)

	// Test: Valid 2 headers with existing headers
	headers = NewHeaders()
	headers["user-agent"] = "Go-http-client/1.1"
	data = []byte("Host: localhost:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, "Go-http-client/1.1", headers["user-agent"])
	assert.Equal(t, 23, n)
	assert.False(t, done)

	// Test: Valid done
	headers = NewHeaders()
	data = []byte("\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, 2, n)
	assert.True(t, done)

	// Test: Invalid spacing header
	headers = NewHeaders()
	data = []byte("       Host : localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Mixed letter casing on both header name and value normalizes key only
	headers = NewHeaders()
	data = []byte("hOsT: LoCaLhOsT:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "LoCaLhOsT:42069", headers["host"])
	assert.False(t, done)

	// Test: Header with digits in the value
	headers = NewHeaders()
	data = []byte("Content-Length: 12345\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "12345", headers["content-length"])
	assert.False(t, done)

	// Test: Header with special characters in the value
	headers = NewHeaders()
	data = []byte("X-Custom-Header: !@#$%^&*()\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "!@#$%^&*()", headers["x-custom-header"])
	assert.False(t, done)

	// Test: Header key with mixed casing is normalized to lowercase
	headers = NewHeaders()
	data = []byte("X-CuStOm-HeAdEr: some-custom-value\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "some-custom-value", headers["x-custom-header"])
	assert.False(t, done)


	// Test: Header key with special character is invalid 
	headers = NewHeaders()
	data = []byte("X-CustomHe@der: some-value\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Same header key repeated is normalized into a single string value 
	headers = NewHeaders()
	headers["x-custom-header"] = "first-value"
	data = []byte("X-Custom-Header: second-value\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "first-value, second-value", headers["x-custom-header"])
	assert.False(t, done)
}