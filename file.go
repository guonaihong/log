package log

import (
	"os"
	"sync"
	"time"
)

const ErrLimit = errors.New("The length of the written data exceeds the limit")

const defaultPrefix = "default-"

const GB = 1024 * 1024 * 1024
const MB = 1024 * 1024
const KB = 1024

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

func (f *File) checkSize(b []byte) (err error) {
	if len(b) > f.MaxSize {
		return ErrLimit
	}

	sb := f.fd.Stat()
	if sb.Size()+len(b) > f.MaxSize {
		//todo compress
		now := time.Now()
		year, month, day := now.Date()

		newName := fmt.Sprintf("%d%02d%02d%02d%02d%02d",
			year, month, day, now.Hour(), now.Minute(), now.Second())

		os.Rename(sb.Name(), newName+".log")

		err = f.openNew()
		if err != nil {
			return
		}
	}
}

func (f *File) openNew() (err error) {
	f.fd.Close()
	f.fd, err = os.Create(defaultPrefix + "0.log")
	if err != nil {
		return err
	}
}

func (f *File) Write(b []byte) (n int, err error) {
	f.Lock()
	defer f.Unlock()

	err = f.checkSize(b)
	if err != nil {
		return
	}

	return f.fd.Write(b)
}

func (f *File) Close() {
	f.Lock()
	defer f.Unlock()
	f.Close()
}
