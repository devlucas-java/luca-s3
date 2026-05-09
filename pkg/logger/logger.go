package logger

import (
	"log"
	"os"
	"sync"
	"sync/atomic"
)

type Level int32

const (
	TRACE Level = iota
	DEBUG
	INFO
	WARN
	ERROR
	FATAL
)

type Logger struct {
	level atomic.Int32
	mu    sync.Mutex

	fatal *log.Logger
	err   *log.Logger
	info  *log.Logger
	warn  *log.Logger
	debug *log.Logger
	trace *log.Logger
}

var instance = newLogger(DEBUG)

func Instance() *Logger {
	return instance
}

func SetLogLevel(level Level) {
	instance.level.Store(int32(level))
}

func newLogger(level Level) *Logger {
	flags := log.Ldate | log.Ltime | log.Lshortfile

	l := &Logger{
		fatal: log.New(os.Stderr, "-----[FATAL] ", flags),
		err:   log.New(os.Stderr, "----[ERROR] ", flags),
		warn:  log.New(os.Stdout, "----[WARN] ", flags),
		info:  log.New(os.Stdout, "---[INFO] ", flags),
		debug: log.New(os.Stdout, "--[DEBUG] ", flags),
		trace: log.New(os.Stdout, "-[TRACE] ", flags),
	}

	l.level.Store(int32(level))

	return l
}

func (l *Logger) canLog(level Level) bool {
	return level >= Level(l.level.Load())
}

func (l *Logger) Trace(v ...interface{}) {
	if l.canLog(TRACE) {
		l.mu.Lock()
		defer l.mu.Unlock()

		l.trace.Println(v...)
	}
}

func (l *Logger) Debug(v ...interface{}) {
	if l.canLog(DEBUG) {
		l.mu.Lock()
		defer l.mu.Unlock()

		l.debug.Println(v...)
	}
}

func (l *Logger) Info(v ...interface{}) {
	if l.canLog(INFO) {
		l.mu.Lock()
		defer l.mu.Unlock()

		l.info.Println(v...)
	}
}

func (l *Logger) Warn(v ...interface{}) {
	if l.canLog(WARN) {
		l.mu.Lock()
		defer l.mu.Unlock()

		l.warn.Println(v...)
	}
}

func (l *Logger) Error(v ...interface{}) {
	if l.canLog(ERROR) {
		l.mu.Lock()
		defer l.mu.Unlock()

		l.err.Println(v...)
	}
}

func (l *Logger) Fatal(v ...interface{}) {
	if l.canLog(FATAL) {
		l.mu.Lock()
		defer l.mu.Unlock()

		l.fatal.Fatalln(v...)
	}
}

func (l *Logger) Tracef(format string, v ...interface{}) {
	if l.canLog(TRACE) {
		l.mu.Lock()
		defer l.mu.Unlock()

		l.trace.Printf(format, v...)
	}
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	if l.canLog(DEBUG) {
		l.mu.Lock()
		defer l.mu.Unlock()

		l.debug.Printf(format, v...)
	}
}

func (l *Logger) Infof(format string, v ...interface{}) {
	if l.canLog(INFO) {
		l.mu.Lock()
		defer l.mu.Unlock()

		l.info.Printf(format, v...)
	}
}

func (l *Logger) Warnf(format string, v ...interface{}) {
	if l.canLog(WARN) {
		l.mu.Lock()
		defer l.mu.Unlock()

		l.warn.Printf(format, v...)
	}
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	if l.canLog(ERROR) {
		l.mu.Lock()
		defer l.mu.Unlock()

		l.err.Printf(format, v...)
	}
}

func (l *Logger) Fatalf(format string, v ...interface{}) {
	if l.canLog(FATAL) {
		l.mu.Lock()
		defer l.mu.Unlock()

		l.fatal.Fatalf(format, v...)
	}
}
