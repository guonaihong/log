package log

import (
	"os"
	"testing"
)

func TestDebugf(t *testing.T) {
	l := NewLog("debug", "test")

	l.AddWriter(os.Stdout)

	l.Debugf("hello world\n")
	l.Infof("hello world\n")
	l.Warnf("hello world\n")
	l.Errorf("hello world\n")
}
