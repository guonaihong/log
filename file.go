package log

type File struct {
	MaxSize  int
	Compress string
}

func (f *File) Write(b []byte) (n int, err error) {
}
