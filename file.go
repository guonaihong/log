package log

import (
	"compress/gzip"
	"io/ioutil"
	"os"
	"strconv"
	"sync"
	"time"
)

//TODO return err
const (
	ErrLimit   = errors.New("The length of the written data exceeds the limit")
	ErrNotDoIt = errors.New("Don't do it")
)

const defaultPrefix = "default-"

const (
	KB = 1024
	MB = 1024 * 1024
	GB = 1024 * 1024 * 1024
)

const defaultMaxSize = 100 * MB

const (
	Gzip = itoa
	NotCompress
)

type File struct {
	compress  string
	maxSize   int
	maxArchve int

	close    bool
	dir      string
	filename chan string
	quitDel  chan struct{}
	maybeDel chan struct{}
	fd       *os.File
	sync.RWMutex
}

func genFileSuffix(compressType int) (suffix string) {
	switch compressType {
	case Gzip:
		suffix = ".gz"
	}

	return
}

func NewFile(prefix, compress int, maxSize, maxArchive int) (f *File) {

	if prefix == "" {
		prefix = defaultPrefix
	}

	if maxSize == 0 {
		maxSize = defaultMaxSize
	}

	f = &File{
		compress:   compress,
		prefix:     prefix,
		maxSize:    maxSize,
		maxArchive: maxArchive,
		filename:   make(chan string, 1000),
		maybeDel:   make(chan struct{}, 1000),
	}

	go f.compressLoop()
	go f.sortAndDel()

	return f
}

func (f *File) sortFile() []os.FileInfo {
	files, err := ioutil.ReadDir(f.dir)
	if err != nil {
		return
	}
}

func (f *File) sortAndDel() {
	for {
		select {
		case <-f.maybeDel:
			//todo sort and del the compress file
		case _, ok := <-f.quitDel:
			if ok {
				return
			}
		}
	}
}

func (f *File) gzipFile(name string) error {
	suffix := genFileSuffix(f.compress)
	if suffix == "" {
		return ErrNotDoIt
	}

	inFd, err := os.Open(name)
	if err != nil {
		return err
	}
	defer inFd.Close()

	outFd, err := os.Create(name + suffix)
	if err != nil {
		return err
	}

	defer outFd.Close()

	zw := gzip.NewWriter(outFd)
	zw.Name = name
	zw.ModTime = time.Now()
	io.Copy(zw, inFd)
	zw.Close()
}

func (f *File) compressLoop() {

	for v := range f.filename {

		err := f.gzipFile(v)
		if err != nil {
			continue
		}

		f.sendDelAsync()
	}
}

func genFileName() (newName string) {

	now := time.Now()

	year, month, day := now.Date()

	newName =
		fmt.Sprintf("%d%02d%02d%02d%02d%02d",
			year,
			month,
			day,
			now.Hour(),
			now.Minute(),
			now.Second())

	newName0 := ""
	i := 0

	for ; ; i++ {
		if i == 0 {
			newName0 = fmt.Sprintf("%s.log", newName)
		} else {
			newName0 = fmt.Sprintf("%s%d.log", newName, i)
		}

		_, err := os.Stat(newName0)
		if os.IsNotExist(err) {
			break
		}
	}

	newName = newName0
	return
}

func (f *File) fileNameNew() string {
	return f.prefix + genFileName() + ".log"
}

func (f *File) sendDelAsync() {
	select {
	case f.maybeDel <- struct{}{}:
	default:
	}
}

func (f *File) checkSize(b []byte) (err error) {
	if len(b) > f.MaxSize {
		return ErrLimit
	}

	if f.fd == nil {
		if err = f.openNew(); err != nil {
			return err
		}
		return
	}

	sb := f.fd.Stat()
	if sb.Size()+len(b) > f.MaxSize {

		os.Rename(sb.Name(), f.fileNameNew())

		select {
		case f.filename <- sb.Name():
		default:
			f.sendDelAsync()
		}

		err = f.openNew()
		if err != nil {
			return
		}
	}
}

func (f *File) openNew() (err error) {
	if f.fd != nil {
		f.fd.Close()
		f.fd = nil
	}

	f.fd, err = os.Create(f.prefix + ".log")
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
	if f.close {
		return
	}

	f.Lock()
	defer f.Unlock()
	if f.close {
		return
	}

	f.Close()
	f.close = true
	close(f.quitDel)
	close(f.filename)
}
