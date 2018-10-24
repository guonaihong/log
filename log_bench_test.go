package log

import (
	"os"
	"testing"
)

func Benchmark2Output(b *testing.B) {
	//b.SetParallelism(100)

	nullf, err := os.OpenFile("/dev/null", os.O_WRONLY, 0666)
	if err != nil {
		b.Fatalf("%v", err)
	}
	defer nullf.Close()

	l := New("test", "test", nullf)
	for i := 0; i < b.N; i++ {
		id := func() *Log {
			var nlog *Log
			if true {
				tmp := *l
				nlog = &tmp
			} else {
				nlog = &Log{}
			}
			nlog.init(l)
			return nlog
		}

		id().Infof("%p\n", l)
	}
}

func BenchmarkOutput(b *testing.B) {
	//b.SetParallelism(100)

	nullf, err := os.OpenFile("/dev/null", os.O_WRONLY, 0666)
	if err != nil {
		b.Fatalf("%v", err)
	}
	defer nullf.Close()
	log := New("debug", "test", nullf)

	for i := 0; i < b.N; i++ {
		//log.F(1).Infof("this is a dummy log\n")
		log.ID("sessionID").Infof("this is a dummy log\n")
	}
}
