package server

import (
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

type Handler func(w *response.Writer, req *request.Request)

// the handle function
func (s *Server) handle(conn io.ReadWriteCloser) {
	defer conn.Close()
	slog.Info("runConnection#")

	responseWriter := response.NewWriter(conn)
	// get request
	req, RequestFromReaderError := request.RequestFromReader(conn) // read the request
	if RequestFromReaderError != nil {
		slog.Info("RequestFromReaderError", "error", RequestFromReaderError)
		responseWriter.WriteStatusLine(response.StatusBadRequest)
		responseWriter.WriteHeaders(*response.GetDefaultHeaders(0))
		return
	}

	// handle request
	//writer := bytes.NewBuffer([]byte{}) // create a buffer to write to (becomes a writer)
	s.handler(responseWriter, req) // call the handler function
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
