package main

import (
	"bufio"
	"fmt"
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
}

func decodeBulkString(byteStream *bufio.Reader) (Value, error) {
}

func decodeArray(byteStream *bufio.Reader) (Value, error) {

}
