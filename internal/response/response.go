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

func GetDefaultHeaders(contentLen int) *headers.Headers {
	h := headers.NewHeaders()

	h.Set("Content-Length", fmt.Sprintf("%d", contentLen))
	h.Set("Connection", "close")
	h.Set("Content-Type", "text/plain")
	slog.Info("GettingDefaultHeaders", "h", h)

	return h
}

// take a poiner to the headers map, iterates over it and writes them to the writer

// WRITER functions
type Writer struct {
	writer io.Writer
}

// can write a custom status line
func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
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
	_, err := w.writer.Write(statusLine)
	return err
}

// can write custom headers
func (w *Writer) WriteHeaders(h headers.Headers) error {
	var err error = nil
	b := []byte{}

	h.ForEach(func(n, v string) {
		b = fmt.Appendf(b, "%s: %s\r\n", n, v)
	})
	b = fmt.Appendf(b, "\r\n")

	slog.Info("WriteHeaders", "headers", b)
	w.writer.Write(b)

	return err
}

// can write a custom body
func (w *Writer) WriteBody(p []byte) (int, error) {
	slog.Info("WriteBody", "body", p)
	n, err := w.writer.Write(p)
	if err != nil {
	}

	// write header length
	return n, nil
}
