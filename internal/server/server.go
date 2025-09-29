package server

import (
	"bytes"
	"fmt"
	"io"
	"ithink/internal/request"
	"ithink/internal/response"
	"log/slog"
	"net"
)

type Server struct {
	closed  bool
	handler Handler
}

type HandlerError struct {
	StatusCode response.StatusCode
	Message    string
}

type Handler func(w io.Writer, req *request.Request) *HandlerError

// the handle function
func (s *Server) handle(conn io.ReadWriteCloser) {
	defer conn.Close()
	slog.Info("runConnection#")

	// get request
	headers := response.GetDefaultHeaders(0)                       // get default headers
	req, RequestFromReaderError := request.RequestFromReader(conn) // read the request
	if RequestFromReaderError != nil {
		slog.Info("RequestFromReaderError", "error", RequestFromReaderError)
		response.WriteStatusLine(conn, response.StatusBadRequest)
		response.WriteHeaders(conn, headers)
		return
	}

	// handle request
	writer := bytes.NewBuffer([]byte{})    // create a buffer to write to (becomes a writer)
	handlerError := s.handler(writer, req) // call the handler function

	var body []byte = nil
	var status response.StatusCode = response.StatusOK
	if handlerError != nil {
		slog.Info("HandlerError", "error", handlerError)
		status = handlerError.StatusCode
		body = []byte(handlerError.Message)
	} else {
		body = writer.Bytes()
	}

	headers.Replace("Content-Length", fmt.Sprintf("%d", len(body))) // replace the content-length
	response.WriteStatusLine(conn, status)                          // write status line
	response.WriteHeaders(conn, headers)                            // write headers
	conn.Write(body)                                                // write the body
}

func listen(s *Server, listener net.Listener) { // go routine doesn't need error
	slog.Info("runServer#")

	// listener
	for {
		conn, err := listener.Accept()
		if s.closed {
			return
		}
		if err != nil {
			return
		}

		go s.handle(conn) // go routine to
	}

}

func Serve(port uint16, handler Handler) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}

	server := &Server{
		closed:  false,
		handler: handler,
	}
	go listen(server, listener) // go routine

	return server, nil
}

func (s *Server) Close() error {
	s.closed = true
	return nil
}
