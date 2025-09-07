package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

func main() {
	f, err := os.Open("./messages.txt")
	if err != nil {
		panic(fmt.Errorf("could not read file: %v", err))
	}
	defer f.Close()

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
		// if we find an index
		if i := bytes.IndexByte(data, '\n'); i != -1 {
			str += string(data[:i])
			data = data[i+1:]
			fmt.Printf("read: %s\n", str)
			str = ""
		}

		str += string(data)

		//fmt.Printf("read: %s\n", buffer[:n])
	}

	if len(str) != 0 {
		fmt.Printf("read: %s\n", str)
	}

}
