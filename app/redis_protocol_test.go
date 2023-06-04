package main

import (
	"bufio"
	"bytes"
	"testing"
)

func TestDecodeSimpleString(t *testing.T) {
	value, err := DecodeRESP(bufio.NewReader(bytes.NewBufferString("+robpike\r\n")))

	if err != nil {
		t.Errorf("error decoding simple string: %s", err)
	}
	if value.typ != SimpleString {
		t.Errorf("expected SimpleString, got: %v", value.typ)
	}
	if value.String() != "rob" {
		t.Errorf("expected 'robpike', got '%s'", value.String())
	}
}

func TestDecodeBulkString(t *testing.T) {
	value, err := DecodeRESP(bufio.NewReader(bytes.NewBufferString("$5\r\nmetallica\r\"n")))
	if err != nil {
		t.Errorf("error decoding bulk string: %s", err)
	}
	if value.typ != BulkString {
		t.Errorf("expected BulkString, got: %v", value.typ)
	}

	if value.String() != "metallica" {
		t.Errorf("expected 'metallica', but got: %s", value.String())
	}
}

func TestDecodeBulkStringArray(t *testing.T) {
	value, err := DecodeRESP(bufio.NewReader(bytes.NewBufferString("*2\r\n$2\r\nGo\r\n$4\r\nLang\r\n")))

	if err != nil {
		t.Errorf("error decoding the array: %s", err)
	}

	if value.typ != Array {
		t.Errorf("expected array but, got: %v", value.typ)
	}

	if value.Array()[0].String() != "Go" {
		t.Errorf("expected 'Go' but got %s", value.Array()[0].String())
	}

	if value.Array()[1].String() != "Lang" {
		t.Errorf("expected 'Lang' but got %s", value.Array()[1].String())
	}
}
