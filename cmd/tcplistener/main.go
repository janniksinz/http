package main

import (
	"fmt"
	"ithink/internal/request"
	"log"
	"net"
)

func main() {
	protocol := "tcp"
	addr := ":42069"
	listener, err := net.Listen(protocol, addr)
	if err != nil {
		log.Fatal("error", "error", err)
	}
	defer listener.Close()
	defer fmt.Printf("connection closed")

	for {
		// accept connection
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("error", "error", err)
		} else {
			//fmt.Printf("connection has been accepted: %s\n", conn)
		}

		request, err := request.RequestFromReader(conn)
		if err != nil {
			log.Fatal("error", "error", err)
		}

		fmt.Printf("Request line:\n")
		fmt.Printf("- Method: %s\n", request.RequestLine.Method)
		fmt.Printf("- Target: %s\n", request.RequestLine.RequestTarget)
		fmt.Printf("- Version: %s\n", request.RequestLine.HttpVersion)

	}
}
