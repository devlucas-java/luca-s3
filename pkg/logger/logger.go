package logger

import (
	"log"
	"os"
)

type Level int

const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
	TRACE
)

type Logger struct {
	level Level

	err   *log.Logger
	info  *log.Logger
	warn  *log.Logger
	debug *log.Logger
	trace *log.Logger
}

func NewLogger(level Level) *Logger {
	flags := log.Ldate | log.Ltime | log.Lshortfile

	return &Logger{
		level: level,
		err:   log.New(os.Stderr, "[ERROR] ", flags),
		info:  log.New(os.Stdout, "[INFO] ", flags),
		warn:  log.New(os.Stdout, "[WARN] ", flags),
		debug: log.New(os.Stdout, "[DEBUG] ", flags),
		trace: log.New(os.Stdout, "[TRACE] ", flags),
	}
}

func (l *Logger) Trace(v ...interface{}) {
	if l.level <= TRACE {
		l.err.Println(v...)
	}
}
func (l *Logger) Debug(v ...interface{}) {
	if l.level <= DEBUG {
		l.debug.Println(v...)
	}
}
func (l *Logger) Info(v ...interface{}) {
	if l.level <= INFO {
		l.info.Println(v...)
	}
}
func (l *Logger) Warn(v ...interface{}) {
	if l.level <= WARN {
		l.warn.Println(v...)
	}
}
func (l *Logger) Error(v ...interface{}) {
	if l.level <= ERROR {
		l.err.Println(v...)
	}
}

func (l *Logger) Tracef(format string, v ...interface{}) {
	if l.level <= TRACE {
		l.err.Printf(format, v...)
	}
}
func (l *Logger) Debugf(msg string, v ...interface{}) {
	if l.level <= DEBUG {
		l.debug.Printf(msg, v...)
	}
}
func (l *Logger) Infof(msg string, v ...interface{}) {
	if l.level <= INFO {
		l.info.Printf(msg, v...)
	}
}
func (l *Logger) Warnf(msg string, v ...interface{}) {
	if l.level <= WARN {
		l.warn.Printf(msg, v...)
	}
}
func (l *Logger) Errorf(msg string, v ...interface{}) {
	if l.level <= ERROR {
		l.err.Printf(msg, v...)
	}
}
