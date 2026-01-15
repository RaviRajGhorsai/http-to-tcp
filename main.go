package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

func getLinesChannel(f io.ReadCloser) <-chan string {

	out := make(chan string, 1)

	go func() {

		defer f.Close()

		defer close(out)

		current_line := ""

		for {
			
			// reads messages in 8 bytes chunk
			data := make([]byte, 8)

			n, err := f.Read(data)

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
				if err != io.EOF {
					fmt.Println(err) 
				}
				break
			}

		}
	}()

	return out

}

func main() {
	
	// open file 
	file, err := os.Open("messages.txt")
	if err != nil {
		fmt.Println(err)
	}
	
	
	lines := getLinesChannel(file)

	// prints the messages that is read as 8 byte chunk
	for line := range lines {
		fmt.Printf("read: %s\n", line)

	}

}
