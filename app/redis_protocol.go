package main

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

//TODO DecodeResp
