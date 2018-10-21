package log

import (
	"bytes"
	"fmt"
	"io"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

type Log struct {
	procName string

	*sync.Mutex
	*sync.WaitGroup

	buf *bytes.Buffer

	w []io.Writer

	level     int
	funcFrame int
	sessionID string
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
		procName:  procName,
		Mutex:     &sync.Mutex{},
		WaitGroup: &sync.WaitGroup{},
		level:     levelStr2N(level),
		buf:       bytes.NewBuffer(make([]byte, 512)),
		funcFrame: 3,
		w:         append([]io.Writer{nil}, w...),
	}
}

func (l *Log) F(frame int) *Log {
	log := *l
	log.funcFrame += frame
	return &log
}

func (l *Log) ID(sessionID string) *Log {
	log := *l
	log.sessionID = sessionID
	return &log
}

func (l *Log) formatHeader(caller bool, level string) {

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
		"[%s%s] ",
		level,
		padding,
	)

	if caller {
		_, file, line, ok := runtime.Caller(l.funcFrame)
		if !ok {
			file = "???"
			line = 0
		}
		fmt.Fprintf(
			l.buf,
			`[%s:%d] `,
			filepath.Base(file),
			line,
		)
	}

	if len(l.sessionID) > 0 {
		fmt.Fprintf(
			l.buf,
			"<sid:%s> ",
			l.sessionID,
		)
	}
}

func (l *Log) AddWriter(w ...io.Writer) {
	l.w = append(l.w, w...)
}

func (l *Log) multWrite(caller bool, level string, a ...interface{}) {
	l.Lock()
	defer func() {
		l.buf.Reset()
		l.Unlock()
	}()

	l.formatHeader(caller, level)
	fmt.Fprint(l.buf, a...)

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

func (l *Log) multWritef(caller bool, level string, format string, a ...interface{}) {
	l.Lock()
	defer func() {
		l.buf.Reset()
		l.Unlock()
	}()

	l.formatHeader(caller, level)
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

	l.multWritef(false, "debug", format, a...)
}

func (l *Log) Debug(a ...interface{}) {
	if l.level > lDebug {
		return
	}

	l.multWrite(false, "debug", a...)
}

func (l *Log) Infof(format string, a ...interface{}) {
	if l.level > lInfo {
		return
	}

	l.multWritef(false, "info", format, a...)
}

func (l *Log) Info(a ...interface{}) {
	if l.level > lInfo {
		return
	}

	l.multWrite(false, "info", a...)
}

func (l *Log) Warnf(format string, a ...interface{}) {
	if l.level > lWarn {
		return
	}

	l.multWritef(true, "warn", format, a...)
}

func (l *Log) Warn(a ...interface{}) {
	if l.level > lWarn {
		return
	}

	l.multWrite(true, "warn", a...)
}

func (l *Log) Errorf(format string, a ...interface{}) {
	if l.level > lError {
		return
	}

	l.multWritef(true, "error", format, a...)
}

func (l *Log) Error(a ...interface{}) {
	if l.level > lError {
		return
	}

	l.multWrite(true, "error", a...)
}
