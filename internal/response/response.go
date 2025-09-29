package response

import (
	"fmt"
	"io"
	"ithink/internal/headers"
	"ithink/internal/request"
	"log/slog"
)

// start-line CRLF
// *(field-line CLRF)
// CRLF
// [ message body ]
type Response struct {
}

type StatusCode int

const (
	StatusOK                  StatusCode = 200
	StatusBadRequest          StatusCode = 400
	StatusInternalServerError StatusCode = 500
)

type HandlerError struct {
	StatusCode StatusCode
	Message    string
}
type Handler func(w io.Writer, req *request.Request) *HandlerError

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	statusLine := []byte{}
	switch statusCode {
	case StatusOK:
		statusLine = []byte("HTTP/1.1 200 OK\r\n")
	case StatusBadRequest:
		statusLine = []byte("HTTP/1.1 400 Bad Request\r\n")
	case StatusInternalServerError:
		statusLine = []byte("HTTP/1.1 500 Internal Server Error\r\n")
	default:
		return fmt.Errorf("unrecognized error code")
	}

	slog.Info("WritingStatusLine", "statusLine", statusLine)
	_, err := w.Write(statusLine)
	return err
}

func GetDefaultHeaders(contentLen int) *headers.Headers {
	h := headers.NewHeaders()

	h.Set("Content-Length", fmt.Sprintf("%d", contentLen))
	h.Set("Connection", "close")
	h.Set("Content-Type", "text/plain")
	slog.Info("GettingDefaultHeaders", "h", h)

	return h
}

// take a poiner to the headers map, iterates over it and writes them to the writer
func WriteHeaders(w io.Writer, headers *headers.Headers) error {
	var err error = nil
	b := []byte{}

	headers.ForEach(func(n, v string) {
		b = fmt.Appendf(b, "%s: %s\r\n", n, v)
	})
	b = fmt.Appendf(b, "\r\n")
	slog.Info("WritingHeaders", "headers", b)
	w.Write(b)

	return err
}
