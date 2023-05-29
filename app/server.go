package main

import (
	"bufio"
	"errors"
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
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error occured while accepting connection: ", err.Error())
			os.Exit(1)
		}
		go handleConnection(conn)
	}
}

// redis uses resp as its protocol
func handleConnection(conn net.Conn) {
	defer conn.Close()
	for {
		value, err := DecodeRESP(bufio.NewReader(conn))
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			fmt.Println("Error decoding RESP: ", err.Error())
			return
		}
		command, err := value.Array()[0].String()
		args := value.Array()[1:]

		switch command {
		case "PING":
			conn.Write([]byte("+PONG\r\n"))
		//ECHO format in resp is: *2\r\n$4\r\nECHO\r\n$6\r\nfoobar\r\n
		//args: ["ECHO","ONE","TWO"]
		case "ECHO":
			conn.Write([]byte(fmt.Sprintf("$%\r\n%s\r\n", len(args[0].String()), args[0].String())))
		default:
			conn.Write([]byte("-ERR unknown command '" + command + "'\r\n"))
		}
	}
}
