package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFieldLineParse(t *testing.T) {
	// Test: Valid single header
	headers := Headers{}
	data := []byte("Host: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["Host"])
	assert.Equal(t, 23, n)
	assert.Equal(t, 1, len(headers))
	assert.False(t, done)

	// Test: Valid single header with extra whitespace
	headers = Headers{}
	data = []byte("   Host:    localhost:42069   \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	assert.Equal(t, "localhost:42069", headers["Host"])
	assert.Equal(t, 32, n)
	assert.Equal(t, 1, len(headers))
	assert.False(t, done)

	// Test: Valid 2 headers with existing headers
	headers = Headers{}
	headers["Accept-Language"] = "en-US"
	headers["Connection"] = "keep-alive"
	assert.Equal(t, 2, len(headers))
	data = []byte("Host: localhost:42069\r\nUser-Agent: Mozilla/5.0\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	assert.Equal(t, "localhost:42069", headers["Host"])
	assert.Equal(t, 23, n)
	assert.Equal(t, 3, len(headers))
	assert.False(t, done)
	data = data[n:]
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	assert.Equal(t, "Mozilla/5.0", headers["User-Agent"])
	assert.Equal(t, 25, n)
	assert.Equal(t, 4, len(headers))
	assert.False(t, done)
	data = data[n:]
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	assert.Equal(t, 2, n)
	assert.Equal(t, 4, len(headers))
	assert.True(t, done)

	// Test: Valid done
	headers = Headers{}
	data = []byte("\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	assert.Equal(t, 2, n)
	assert.Equal(t, 0, len(headers))
	assert.True(t, done)

	// Test: Invalid spacing header
	headers = Headers{}
	data = []byte("       Host : localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.Equal(t, 0, len(headers))
	assert.False(t, done)
}
