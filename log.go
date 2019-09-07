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
	"unsafe"
)

type Log struct {
	procName string

	*sync.Mutex
	*sync.WaitGroup
	*sync.Pool

	buf *bytes.Buffer

	w []io.Writer

	addr      *Log
	level     int
	funcFrame int
	sessionID string
	cb        func(s string)
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

func New(level string, procName string, w ...io.Writer) *Log {
	l := &Log{
		procName:  procName,
		Mutex:     &sync.Mutex{},
		Pool:      &sync.Pool{},
		WaitGroup: &sync.WaitGroup{},
		level:     levelStr2N(level),
		buf:       bytes.NewBuffer(make([]byte, 512)),
		funcFrame: 3,
		w:         append([]io.Writer{nil}, w...),
	}

	l.addr = l
	return l
}

func (l *Log) newLog(reuse bool) *Log {

	var nlog *Log
	var ok bool

	if reuse {
		nlog, ok = l.Get().(*Log)

		if !ok {
			nlog = &Log{}
		}
	} else {
		nlog = &Log{}
		//tmp := *l
		//nlog = &tmp
	}

	nlog.init(l)
	return nlog
}

func (l *Log) releaseLog() {
	if l.addr != l {
		l.Put(l)
	}
}

func (l *Log) init(base *Log) {
	l.procName = base.procName
	l.Mutex = base.Mutex
	l.Pool = base.Pool
	l.WaitGroup = base.WaitGroup
	l.level = base.level
	l.buf = base.buf
	l.addr = base.addr
	l.w = base.w
	l.funcFrame = base.funcFrame
	l.sessionID = ""
}

func (l *Log) Cb(cb func(s string)) *Log {
	log := l.newLog(false)
	log.cb = cb
	return log
}

func (l *Log) F(frame int) *Log {
	log := l.newLog(false)
	log.funcFrame += frame
	return log
}

func (l *Log) ID(sessionID string) *Log {
	log := l.newLog(false)
	log.sessionID = sessionID
	return log
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
	l.buf.Reset()
	defer func() {
		l.Unlock()
		//l.releaseLog()
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

	b := l.buf.Bytes()
	if l.cb != nil {
		l.cb(*(*string)(unsafe.Pointer(&b)))
	}
	l.Wait()
}

func (l *Log) multWritef(caller bool, level string, format string, a ...interface{}) {
	l.Lock()
	l.buf.Reset()
	defer func() {
		l.Unlock()
		//l.releaseLog()
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

	b := l.buf.Bytes()
	if l.cb != nil {
		l.cb(*(*string)(unsafe.Pointer(&b)))
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
