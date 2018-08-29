package log

import (
	"os"
	"sync"
)

const ErrLimit = errors.New("The length of the written data exceeds the limit")

const defaultPrefix = "default-"

type File struct {
	Compress  string
	MaxSize   int
	MaxArchve int

	fd *os.File
	sync.Mutex
}

func NewFile(prefix, compress string, maxSize, maxArchive int) *File {
	if compress == "" {
		compress = "gz"
	}

	if prefix == "" {
		prefix = defaultPrefix
	}

	return &File{
		Compress:   compress,
		prefix:     prefix,
		MaxSize:    maxSize,
		maxArchive: maxArchive,
	}
}

func (f *File) checkSize(b []byte) error {
	sb := f.fd.Stat()
	if sb.Size() + len(b) {
	}
}

func (f *File) openNew() {
}

func (f *File) Write(b []byte) (n int, err error) {
	f.Lock()
	defer f.Unlock()

	if len(b) > f.MaxSize {
		return 0, ErrLimit
	}
}

func (f *File) Close() {
}
