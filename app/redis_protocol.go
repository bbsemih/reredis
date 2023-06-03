package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

type Type byte

const (
	//OK -->"+OK\r\n"
	SimpleString Type = '+'
	//hello--> "$5\r\nhello\r\n"
	BulkString Type = '$'
	//"*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n"
	Array Type = '*'
)

// Value represents the data of a valid RESP type.
type Value struct {
	typ   Type
	bytes []byte
	array []Value
}

// converts value into string
func (v Value) String() string {
	if v.typ == BulkString || v.typ == SimpleString {
		return string(v.bytes)
	}
	return ""
}

// converts value into array
func (v Value) Array() []Value {
	if v.typ == Array {
		return v.array
	}
	return []Value{}
}

// This parses a RESP message and returns a RedisValue
func DecodeRESP(byteStream *bufio.Reader) (Value, error) {
	dataTypeByte, err := byteStream.ReadByte()
	if err != nil {
		return Value{}, err
	}

	switch string(dataTypeByte) {
	case "+":
		return decodeSimpleString(byteStream)
	case "$":
		return decodeBulkString(byteStream)
	case "*":
		return decodeArray(byteStream)
	}
	return Value{}, fmt.Errorf("invalid RESP data type byte: %s", string(dataTypeByte))
}

func decodeSimpleString(byteStream *bufio.Reader) (Value, error) {
	readBytes, err := readUntilCRLF(byteStream)
	if err != nil {
		return Value{}, err
	}
	return Value{
		typ:   SimpleString,
		bytes: readBytes,
	}, nil
}

// BulkString --> $5\r\nhello\r\n
func decodeBulkString(byteStream *bufio.Reader) (Value, error) {
	readBytesForCount, err := readUntilCRLF(byteStream)
	if err != nil {
		return Value{}, fmt.Errorf("failed to bulk string length: %s", err)
	}
	//Atoi = to integer - Itoa = to string
	count, err := strconv.Atoi(string(readBytesForCount))
	if err != nil {
		return Value{}, fmt.Errorf("failed to bulk string length: %s", err)
	}
	//Create a byte slice with the appropriate length to store the bulk string contents
	readBytes := make([]byte, count+2)

	if _, err := io.ReadFull(byteStream, readBytes); err != nil {
		return Value{}, fmt.Errorf("failed to read bulk string contents: %s", err)
	}
	return Value{
		typ:   BulkString,
		bytes: readBytes[:count],
	}, nil
}

// "*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n"
func decodeArray(byteStream *bufio.Reader) (Value, error) {
	readBytesForCount, err := readUntilCRLF(byteStream)
	if err != nil {
		return Value{}, fmt.Errorf("failed to bulk string length: %s", err)
	}
	count, err := strconv.Atoi(string(readBytesForCount))
	if err != nil {
		return Value{}, fmt.Errorf("failed to parse bulk string length: %s", err)
	}
	array := []Value{}

	for i := 1; i < count; i++ {
		value, err := DecodeRESP(byteStream)
		if err != nil {
			return Value{}, err
		}
		array = append(array, value)
	}
	return Value{
		typ:   Array,
		array: array,
	}, nil
}

// CRLF --> CR (Carriage Return(\r)) LF (Line Feed\l)
// Result of input "Hello\nWorld\nThis is a test\r\n" would be
// Hello World This is a test
func readUntilCRLF(byteStream *bufio.Reader) ([]byte, error) {
	readBytes := []byte{}
	for {
		//returning a slice containing the data up to and including the delimiter
		b, err := byteStream.ReadBytes('\n')
		if err != nil {
			return nil, err
		}
		readBytes = append(readBytes, b...)
		if len(readBytes) >= 2 && readBytes[len(readBytes)-2] == '\r' {
			break
		}
	}
	return readBytes[:len(readBytes)-2], nil
}
