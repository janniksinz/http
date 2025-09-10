package headers

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestHeaderParser(t *testing.T) {
	// Test: Valid single header
	headers := NewHeaders()
	data := []byte("Host: localhost:42069\r\nFooFoo:    barbar      \r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers.Get("Host"))
	assert.Equal(t, "barbar", headers.Get("FooFoo"))
	assert.Equal(t, 50, n)
	assert.True(t, done)

	// Test: Valid single header
	headers = NewHeaders()
	data = []byte("Host: localhost:42069\r\nFooFoo:barbar\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers.Get("Host"))
	assert.Equal(t, "barbar", headers.Get("FooFoo"))
	assert.Equal(t, "", headers.Get("MissingKey"))
	assert.Equal(t, 40, n)
	assert.True(t, done)

	// Test: Invalid Header characters
	headers = NewHeaders()
	data = []byte("H@st: localhost:42069\r\nFooFoo:barbar\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Invalid spacing header
	headers = NewHeaders()
	data = []byte("       Host : localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)
}
