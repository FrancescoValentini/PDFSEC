package cli

import (
	"fmt"
	"io"
	"os"
	"time"
)

// Logger wraps output with timestamps and log levels.
type Logger struct {
	out     io.Writer
	verbose bool
}

func NewLogger(verbose bool) *Logger {
	return &Logger{out: os.Stderr, verbose: verbose}
}

func (l *Logger) log(level, msg string) {
	if !l.verbose {
		return
	}
	ts := time.Now().Format("2006-01-02 15:04:05")
	fmt.Fprintf(l.out, "[%s] %s: %s\n", ts, level, msg)
}

func (l *Logger) Info(msg string)  { l.log("INFO", msg) }
func (l *Logger) Warn(msg string)  { l.log("WARN", msg) }
func (l *Logger) Error(msg string) { l.log("ERROR", msg) }
