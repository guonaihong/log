package log

import (
	"testing"
)

func TestParseSocket(t *testing.T) {
	w, err := ParseSocket("udp://127.0.0.1:1234")
	if err != nil {
		t.Fatalf("parse socket fail:%s\n", err)
	}

	_, err = w.Write([]byte("hello world"))
	if err != nil {
		t.Fatalf("write data fail:%s\n", err)
	}
}
