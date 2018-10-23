package log

import (
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	ErrLimit   = errors.New("The length of the written data exceeds the limit")
	ErrNotDoIt = errors.New("Don't do it")
)

const defaultPrefix = "default"

const (
	KB = 1024
	MB = 1024 * 1024
	GB = 1024 * 1024 * 1024
)

const defaultMaxSize = 100 * MB

type CompressType int

const (
	Gzip CompressType = iota
	NotCompress
)

type File struct {
	compress   CompressType
	maxSize    int
	maxArchive int

	close       bool
	dir         string
	prefix      string
	defaultName string
	filename    chan string
	quitDel     chan struct{}
	maybeDel    chan struct{}
	errs        chan error
	fd          *os.File
	sync.RWMutex
	sync.WaitGroup
}

func genFileSuffix(compressType CompressType) (suffix string) {
	switch compressType {
	case Gzip:
		suffix = ".gz"
	}

	return
}

func NewFile(prefix string, dir string, compress CompressType, maxSize, maxArchive int) (f *File) {

	if prefix == "" {
		prefix = defaultPrefix
	}

	if maxSize == 0 {
		maxSize = defaultMaxSize
	}

	name := ""

	if len(dir) > 0 && (dir[len(dir)-1] != '/' && dir[len(dir)-1] != '\\') {
		name = filepath.Base(dir)
		if name == "." {
			name = ""
		} else if !strings.HasSuffix(name, ".log") {
			name += ".log"
		}

	}

	dir = filepath.Dir(dir)
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		os.MkdirAll(dir, 0755)
	}

	f = &File{
		dir:         dir,
		defaultName: strings.TrimSpace(name),
		compress:    compress,
		prefix:      prefix,
		maxSize:     maxSize,
		maxArchive:  maxArchive,
		filename:    make(chan string, 1000),
		maybeDel:    make(chan struct{}, 1000),
		quitDel:     make(chan struct{}),
		errs:        make(chan error, 10),
	}

	f.Add(2)
	go f.compressLoop()
	go f.sortAndDel()

	return f
}

func (f *File) fileNameNew() string {
	return f.dir + "/" + f.prefix + genFileName()
}

func (f *File) defaultFileName() string {
	if f.defaultName != "" {
		return filepath.Clean(f.dir + "/" + f.prefix + f.defaultName)
	}

	return f.dir + "/" + f.prefix + ".log"
}

func (f *File) getTimeFormFile(name string) string {

	if strings.HasSuffix(name, ".log") {
		name = name[:len(name)-len(".log")]
	}

	if strings.HasPrefix(name, f.prefix) {
		name = name[len(f.prefix):]
	}

	return name
}

func (f *File) sortFile(files0 []os.FileInfo) []os.FileInfo {

	var files []os.FileInfo

	for _, v := range files0 {

		if !strings.HasPrefix(v.Name(), f.prefix) {
			continue
		}

		if f.dir+"/"+v.Name() == f.defaultFileName() {
			continue
		}

		if strings.HasSuffix(v.Name(), ".log") ||
			strings.HasSuffix(v.Name(), ".gz") {
			files = append(files, v)
		}

	}

	sort.Slice(files, func(i, j int) bool {

		pathA := f.getTimeFormFile(files[i].Name())
		pathB := f.getTimeFormFile(files[j].Name())

		time1, err1 := strconv.Atoi(pathA)
		time2, err2 := strconv.Atoi(pathB)

		if err1 != nil || err2 != nil {
			return pathA < pathB
		}

		return time1 < time2
	})

	return files
}

func printFile(prefix string, files []os.FileInfo) {
	for _, v := range files {
		fmt.Printf("%s:%s\n", prefix, v.Name())
	}
}

func (f *File) recvError() error {
	select {
	case err, _ := <-f.errs:
		return err
	default:
	}

	return nil
}

func (f *File) sendError(err error) {
	for {
		select {
		case f.errs <- err:
			return
		default:
			select {
			case <-f.errs:
			default:
			}
		}
	}
}

func (f *File) getOldFile() []os.FileInfo {
	files0, err := ioutil.ReadDir(f.dir)
	if err != nil {
		f.sendError(err)
		return nil
	}

	//printFile("sort before", files0)

	files := f.sortFile(files0)

	//printFile("sort after", files)

	if len(files) <= f.maxArchive {
		return nil
	}

	if len(files) > 0 {
		return files[0 : len(files)-f.maxArchive]
	}

	return nil
}

func (f *File) sortAndDel() {

	defer f.Done()

	for {
		select {
		case <-f.maybeDel:
			files := f.getOldFile()
			for _, v := range files {
				if err := os.Remove(f.dir + "/" + v.Name()); err != nil {
					f.sendError(err)
				}
			}
		case _, ok := <-f.quitDel:
			if !ok {
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

	if err = os.Remove(name); err != nil {
		fmt.Printf("%s\n", err)
	}

	return nil
}

func (f *File) compressLoop() {

	defer f.Done()
	for v := range f.filename {

		err := f.gzipFile(v)
		if err != nil {
			f.sendError(err)
			continue
		}

		f.sendDelAsync()
	}
}

func genFileName() (newName string) {

	now := time.Now()

	year, month, day := now.Date()
	hour, minute, second := now.Clock()
	newName =
		fmt.Sprintf("%d%02d%02d%02d%02d%02d",
			year,
			month,
			day,
			hour,
			minute,
			second)

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

func (f *File) sendDelAsync() {
	select {
	case f.maybeDel <- struct{}{}:
	default:
	}
}

func (f *File) checkSize(b []byte) (err error) {
	if len(b) > f.maxSize {
		return ErrLimit
	}

	if f.fd == nil {
		if err = f.openNew(); err != nil {
			return err
		}
		return
	}

	sb, err := f.fd.Stat()
	if err != nil {
		return err
	}

	if int(sb.Size())+len(b) > f.maxSize {

		newName := f.fileNameNew()
		err = os.Rename(f.defaultFileName(), newName)
		if err != nil {
			return
		}

		select {
		case f.filename <- newName:
		default:
			f.sendDelAsync()
		}

		err = f.openNew()
		if err != nil {
			return
		}
	}

	return nil
}

func (f *File) openNew() (err error) {
	if f.fd != nil {
		f.fd.Close()
		f.fd = nil
	}

	f.fd, err = os.OpenFile(f.defaultFileName(), os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (f *File) Write(b []byte) (n int, err error) {
	f.Lock()
	defer f.Unlock()

	err = f.recvError()
	if err != nil {
		return
	}

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

	if f.fd != nil {
		f.fd.Close()
	}
	f.close = true
	close(f.quitDel)
	close(f.filename)
	f.Wait()
}
