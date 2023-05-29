package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

func main() {
	fmt.Println("Starting server...")
	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error occured while accepting connection: ", err.Error())
		os.Exit(1)
	}

	defer conn.Close()
	//Respond to multiple PINGs from client
	for {
		buf := make([]byte, 1024)
		_, err := conn.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("error reading from client:", err.Error())
		}
		conn.Write([]byte("+PONG\r\n"))
	}
}
