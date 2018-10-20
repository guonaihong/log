package log

import (
	"os"
	"testing"
)

func error2(log *Log, a ...interface{}) {
	log.F(1).Error(a...)
}

func TestDebugf(t *testing.T) {
	l := NewLog("debug", "test")

	l.AddWriter(os.Stdout)

	l.Debugf("hello world\n")
	l.Infof("hello world\n")
	l.Info("hello", "world\n")
	l.Warnf("hello world\n")
	l.Warn("hello", "world\n")
	l.Errorf("hello world\n")
	l.Error("hello", "world\n")

	error2(l, "hello2 world2\n")
}
