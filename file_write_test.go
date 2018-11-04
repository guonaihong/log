package log

import (
	"testing"
)

func TestWrite0(t *testing.T) {
	var file *File
	file = NewFile("test-", "./", Gzip, 1*MB, 3)
	defer file.Close()
	for i := 0; i < 1024*1024/len("hello world"); i++ {
		_, err := file.Write([]byte("hello world"))
		if err != nil {
			t.Fatalf("err:%s\n", err)
		}
	}
}

func TestWrite1(t *testing.T) {
	var file *File
	file = NewFile("test-", "./test-my.log", Gzip, 1*MB, 3)
	defer file.Close()
	for i := 0; i < 1024*1024/len("hello world"); i++ {
		_, err := file.Write([]byte("hello world"))
		if err != nil {
			t.Fatalf("err:%s\n", err)
		}
	}
}

func TestWrite2(t *testing.T) {
	var file *File
	file = NewFile("test-", "./mylog/test-my.log", Gzip, 1*MB, 3)
	defer file.Close()
	for i := 0; i < 1024*1024/len("hello world"); i++ {
		_, err := file.Write([]byte("hello world"))
		if err != nil {
			t.Fatalf("err:%s\n", err)
		}
	}
}

func TestWrite3(t *testing.T) {
	var file *File
	file = NewFile("test-", "/tmp/", Gzip, 1*MB, 3)
	defer file.Close()
	for i := 0; i < 1024*1024/len("hello world"); i++ {
		_, err := file.Write([]byte("hello world"))
		if err != nil {
			t.Fatalf("err:%s\n", err)
		}
	}
}

func TestWrite4(t *testing.T) {
	var file *File
	file = NewFile("test-", "/tmp/log/access.log", Gzip, 1*MB, 3)
	defer file.Close()
	for i := 0; i < 1024*1024/len("hello world"); i++ {
		_, err := file.Write([]byte("hello world"))
		if err != nil {
			t.Fatalf("err:%s\n", err)
		}
	}
}

func TestWrite5(t *testing.T) {
	var file *File
	file = NewFile("vpr-log", "./test-log/", Gzip, 1*MB, 3)
	defer file.Close()
	for i := 0; i < 1024*1024/len("hello world"); i++ {
		_, err := file.Write([]byte("hello world"))
		if err != nil {
			t.Fatalf("err:%s\n", err)
		}
	}
}

func TestWrite6(t *testing.T) {
	var file *File
	file = NewFile("", "./test-log/", Gzip, 1*MB, 3)
	defer file.Close()
	for i := 0; i < 1024*1024/len("hello world"); i++ {
		_, err := file.Write([]byte("hello world"))
		if err != nil {
			t.Fatalf("err:%s\n", err)
		}
	}
}

func TestWriteBig(t *testing.T) {
	var file *File
	file = NewFile("test", "./test/", Gzip, 100*MB, 10)
	defer file.Close()
	for i := 0; i < 1024*1024*1024/len("hello world"); i++ {
		_, err := file.Write([]byte("hello world"))
		if err != nil {
			t.Fatalf("err:%s\n", err)
		}
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
