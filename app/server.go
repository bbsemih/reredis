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

	storage := NewStorage()
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error occured while accepting connection: ", err.Error())
			os.Exit(1)
		}
		go handleConnection(conn, storage)
	}
}

// redis uses RESP as its protocol
func handleConnection(conn net.Conn, storage *Storage) {
	defer conn.Close()
	for {
		// This parses a RESP message and returns a RedisValue
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
		case "ping":
			conn.Write([]byte("+PONG\r\n"))
		//ECHO format in resp is: *2\r\n$4\r\nECHO\r\n$6\r\nfoobar\r\n
		//args: ["ECHO","ONE","TWO"]
		case "echo":
			conn.Write([]byte(fmt.Sprintf("$%d\r\n%s\r\n", len(args[0].String()), args[0].String())))
		//SET <key> <value>
		case "get":
			conn.Write([]byte(fmt.Sprintf("+%s\r\n", storage.Get(args[0].String()))))
		case "set":
			storage.Set(args[0].String(), args[1].String())
			conn.Write([]byte("+OK\r\n"))
		default:
			conn.Write([]byte("-ERR unknown command '" + command + "'\r\n"))
		}
	}
}
