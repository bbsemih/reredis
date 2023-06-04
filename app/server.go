package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"time"
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
		command := value.Array()[0].String()
		args := value.Array()[1:]

		switch command {
		case "ping":
			conn.Write([]byte("+PONG\r\n"))
			//ECHO format in resp is: *2\r\n$4\r\nECHO\r\n$6\r\nfoobar\r\n
			//args: ["ECHO","ONE","TWO"]
			//----------TODO-------------
			//del? lists?
		case "echo":
			conn.Write([]byte(fmt.Sprintf("$%d\r\n%s\r\n", len(args[0].String()), args[0].String())))
		case "get":
			value, found := storage.Get(args[0].String())
			if found {
				conn.Write([]byte(fmt.Sprintf("$%d\r\n%s\r\n", len(value), value)))
			} else {
				conn.Write([]byte("$-1\r\n"))
			}
		//SET <key> <value> px <expirelength>
		case "set":
			if len(args) > 2 {
				if args[2].String() == "px" {
					expiryStr := args[3].String()
					expiryInMlSeconds, err := strconv.Atoi(expiryStr)
					if err != nil {
						conn.Write([]byte(fmt.Sprint("-ERR PX value is not an integer\r\n", expiryStr)))
						break
					}
					storage.SetWithExpiry(args[0].String(), args[1].String(), time.Duration(expiryInMlSeconds)*time.Millisecond)
				} else {
					conn.Write([]byte(fmt.Sprintf("-ERR unknown option for set: %s\r\n", args[2].String())))
				}
			} else {
				storage.Set(args[0].String(), args[1].String())
			}
			conn.Write([]byte("+OK\r\n"))
		default:
			conn.Write([]byte("-ERR unknown command '" + command + "'\r\n"))
		}
	}
}
