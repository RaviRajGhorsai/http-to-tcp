package main

import (
	"bytes"
	"fmt"
	"net"
)

func getLinesChannel(conn net.Conn) <-chan string {

	out := make(chan string, 1)

	go func() {

		defer conn.Close()

		defer close(out)

		current_line := ""

		for {

			// reads messages in 8 bytes chunk
			data := make([]byte, 8)

			n, err := conn.Read(data)

			if n != 0 {
				data = data[:n]

				// check if the 8 byte chunk has \n or end of line or next line, they are converted into parts
				if i := bytes.IndexByte(data, '\n'); i != -1 {
					current_line += string(data[:i])
					data = data[i+1:]
					out <- current_line
					current_line = ""
				}
				// the last part is set to current_line because that part may be in complete
				current_line += string(data)
			}

			if err != nil {

				break
			}

		}
	}()

	return out

}

func main() {

	// listener on tcp port 42069
	listener, err := net.Listen("tcp", ":42069")

	if err != nil {

		fmt.Println("error listening on port 42069...")
		return
	}

	// defer will make sure when program exits it closes listener
	defer listener.Close()

	fmt.Println("Server Listening in port 42069")

	for {

		conn, err := listener.Accept()

		if err != nil {

			fmt.Println("Error on establishing connection....")
			continue
		}

		fmt.Println("Connection has been established...")

		// prints the messages that is read as 8 byte chunk
		for line := range getLinesChannel(conn) {

			fmt.Printf("%s\n", line)

		}
	}

}
