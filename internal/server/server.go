package server

import (
	"fmt"
	"io"
	"log/slog"
	"net"
	"time"
)

type Server struct {
	state  string
	closed bool
}

func runConnection(s *Server, conn io.ReadWriteCloser) {
	defer conn.Close()
	slog.Info("runConnection#")
	out := []byte("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: 13\r\n\r\nHello World!")
	_, err := conn.Write(out)
	if err != nil {
		slog.Error("write error", "error", err)
	}
	slog.Info("Return", "Res:", string(out))
	time.Sleep(10 * time.Millisecond)
}

func runServer(s *Server, listener net.Listener) { // go routine doesn't need error
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

		go runConnection(s, conn) // go routine to
	}

}

func Serve(port uint16) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}

	server := &Server{closed: false}
	go runServer(server, listener) // go routine

	return server, nil
}

func (s *Server) Close() error {
	s.closed = true
	return nil
}
