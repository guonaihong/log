package log

type Log struct {
	level int
}

func NewLog(level int) *Log {
}

func (l *Log) Debugf(format string, a ...interface{}) {
}

func (l *Log) Infof(format string, a ...interface{}) {
}

func (l *Log) Warnf(format string, a ...interface{}) {
}

func (l *Log) Errorf(format string, a ...interface{}) {
}
