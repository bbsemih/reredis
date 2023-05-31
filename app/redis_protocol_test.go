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
		t.Errorf("expected SimpleString, got: %s", value.typ)
	}
	if value.String() != "rob" {
		t.Errorf("expected 'robpike', got '%s'", value.String())
	}
}
