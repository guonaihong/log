package log

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"sync"
	"time"
)

type Log struct {
	procName string

	sync.Mutex
	sync.WaitGroup

	buf *bytes.Buffer

	level int

	w []io.Writer
}

const (
	lDebug = iota
	lInfo
	lWarn
	lError
)

func levelStr2N(level string) int {

	switch strings.ToLower(level) {

	case "debug":
		return lDebug
	case "info":
		return lInfo
	case "warn":
		return lWarn
	case "error":
		return lError
	}
	return lDebug
}

func NewLog(level string, procName string, w ...io.Writer) *Log {
	return &Log{
		procName: procName,
		level:    levelStr2N(level),
		buf:      bytes.NewBuffer(make([]byte, 512)),
		w:        append([]io.Writer{nil}, w...),
	}
}

func (l *Log) formatHeader(level string) {

	now := time.Now()

	year, month, day := now.Date()
	hour, min, sec := now.Clock()

	padding := ""
	if len(level) == 4 {
		padding = " "
	}

	fmt.Fprintf(
		l.buf,
		"[%s] ",
		l.procName,
	)

	fmt.Fprintf(
		l.buf,
		"[%4d-%02d-%02d %02d:%02d:%02d.%06d] ",
		year,
		month,
		day,
		hour,
		min,
		sec,
		now.Nanosecond()/1e3)

	fmt.Fprintf(
		l.buf,
		"[%s%s]",
		level,
		padding,
	)
}

func (l *Log) AddWriter(w ...io.Writer) {
	l.w = append(l.w, w...)
}

func (l *Log) multWrite(level string, format string, a ...interface{}) {
	l.Lock()
	defer func() {
		l.buf.Reset()
		l.Unlock()
	}()

	l.formatHeader(level)
	fmt.Fprintf(l.buf, format, a...)

	l.Add(len(l.w))
	for _, w := range l.w {
		go func(w io.Writer) {
			defer l.Done()
			if w, ok := w.(io.Writer); ok {
				w.Write(l.buf.Bytes())
			}
		}(w)
	}

	l.Wait()
}

func (l *Log) Debugf(format string, a ...interface{}) {
	if l.level > lDebug {
		return
	}

	defer l.buf.Reset()

	l.multWrite("debug", format, a...)
}

func (l *Log) Infof(format string, a ...interface{}) {
	if l.level > lInfo {
		return
	}

	l.multWrite("info", format, a...)
}

func (l *Log) Warnf(format string, a ...interface{}) {
	if l.level > lWarn {
		return
	}

	l.multWrite("warn", format, a...)
}

func (l *Log) Errorf(format string, a ...interface{}) {
	if l.level > lError {
		return
	}

	l.multWrite("error", format, a...)
}
