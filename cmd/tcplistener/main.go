package tcplistener

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
)

func getLinesChannel(f io.ReadCloser) <-chan string {
	out := make(chan string, 1)

	go func() {
		defer f.Close()
		defer close(out)

		str := ""
		for {
			buffer := make([]byte, 8)
			n, err := f.Read(buffer) // reads into buffer and returns integer of bytes read
			if err != nil {
				if err == io.EOF {
					break
				}
				panic(err)
			}

			data := buffer[:n]
			// if we find a newline index
			if i := bytes.IndexByte(data, '\n'); i != -1 {
				// only add data until the newline
				str += string(data[:i])
				data = data[i+1:]
				out <- str // return line
				str = ""   // reset
			}

			// add to line
			str += string(data)
		}
		if len(str) != 0 {
			out <- str // return line
		}

	}()

	return out
}

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

		for line := range getLinesChannel(conn) {
			fmt.Printf("%s", line)
		}
	}
}
