package udpsender

import (
	"bytes"
	"fmt"
	"io"
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
	protocol := "udp"
	//listener, err := net.ResolveUDPAddr(protocol, addr)
	addr := net.UDPAddr{Port: 42069}

	conn, err := net.ListenUDP(protocol, &addr)
	if err != nil {
		fmt.Printf("couldn't establish connection")
	}
	defer conn.Close()
	defer fmt.Printf("connection closed")

	buf := make([]byte, 1024)

	for {
		n, client, _ := conn.ReadFromUDP(buf)
		conn.WriteToUDP(buf[:n], client)

	}

	//for line := range getLinesChannel(conn) {
	//	fmt.Printf("%s", line)
	//}
}
