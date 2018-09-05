package log

import (
	"testing"
)

func TestWrite(t *testing.T) {
	var file *File
	//file = NewFile("test-", "./", Gzip, 1*MB, 3)
	file = NewFile("test-", "./test-my.log", Gzip, 1*MB, 3)
	defer file.Close()
	for i := 0; i < 100000; i++ {
		file.Write([]byte("hello world"))
	}
}

func BenchmarkWrite(b *testing.B) {
	var file *File
	file = NewFile("test-", ".", Gzip, 1*MB, 3)
	defer file.Close()
	for i := 0; i < b.N; i++ {
		file.Write([]byte("hello world"))
	}
}
