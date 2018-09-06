package log

import (
	"os"
	"testing"
)

func TestDebugf(t *testing.T) {
	l := NewLog("debug", "test")

	l.AddWriter(os.Stdout)

	l.Debugf("debug hello world\n")
	l.Infof("info hello world\n")
	l.Warnf("warn hello world\n")
	l.Errorf("error hello world\n")
}
