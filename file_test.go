package log

import (
	"testing"
)

var file *File

func TestMain(t *testing.M) {
	file = NewFile("test-", "gz", 100*MB, 10)
}

func BenchmarkWrite(b *testing.B) {
	for i := 0; i < b.N; i++ {
		file.Write("hello world")
	}
}
