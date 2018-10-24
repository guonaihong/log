package log

import (
	"os"
	"sync"
	"testing"
)

func TestDebugf(t *testing.T) {
	error2 := func(log *Log, a ...interface{}) {
		log.F(1).Error(a...)
	}

	l := New("debug", "test")

	l.AddWriter(os.Stdout)

	l.Debugf("hello world\n")
	l.Infof("hello world\n")
	l.Info("hello", " world\n")
	l.Warnf("hello world\n")
	l.Warn("hello", " world\n")
	l.Errorf("hello world\n")
	l.Error("hello", " world\n")

	error2(l /* *Log */, "hello2 world2\n")

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		l.ID("1").Debugf("hello world.1\n")
		l.ID("1").Infof("hello world.1\n")
		l.ID("1").Info("hello", " world.1\n")
		l.ID("1").Warnf("hello world.1\n")
		l.ID("1").Warn("hello", " world.1\n")
		l.ID("1").Errorf("hello world.1\n")
		l.ID("1").Error("hello", " world.1\n")
	}()

	l.ID("2").Debugf("hello world.2\n")
	l.ID("2").Infof("hello world.2\n")
	l.ID("2").Info("hello", " world.2\n")
	l.ID("2").Warnf("hello world.2\n")
	l.ID("2").Warn("hello", " world.2\n")
	l.ID("2").Errorf("hello world.2\n")
	l.ID("2").Error("hello", " world.2\n")
	wg.Wait()
}
