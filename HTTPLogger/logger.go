package HTTPLogger

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

const (
	LOG_ERRORS = 1 << iota
	LOG_INFO
	LOG_WARNINGS
	LOG_DEBUG
)

// Logger is a generalized logger to save some boilerplate. It is made using the standard library.
type Logger struct {
	Flags int // Determines which log streams (err/info/debug/warn) are enabled
	err   *log.Logger
	info  *log.Logger
	debug *log.Logger
	warn  *log.Logger
}

// Error logs application errors and includes file location/line number in output.
// It is not meant for logging HTTP errors or other errors that arise from the client.
func (l *Logger) Error(appErr error) {
	if l.Flags&LOG_ERRORS != 0 {
		err := l.err.Output(2, "| "+appErr.Error())
		if err != nil {
			fmt.Printf("LOGGING ERROR: %v\n", err)
		}
	}
}

// Info logs with prefix INFO
func (l *Logger) Info(msg string) {
	if l.Flags&LOG_INFO != 0 {
		err := l.info.Output(2, "| "+msg)
		if err != nil {
			fmt.Printf("LOGGING ERROR: %v\n", err)
		}
	}
}

// Debug logs with prefix DEBUG
func (l *Logger) Debug(msg string) {
	if l.Flags&LOG_DEBUG != 0 {
		err := l.debug.Output(2, "| "+msg)
		if err != nil {
			fmt.Printf("LOGGING ERROR: %v\n", err)
		}
	}
}

// Warn logs with prefix WARNING
func (l *Logger) Warn(msg string) {
	if l.Flags&LOG_WARNINGS != 0 {
		err := l.warn.Output(2, "| "+msg)
		if err != nil {
			fmt.Printf("LOGGING ERROR: %v\n", err)
		}
	}
}

// InfoHTTP prints an INFO log message with the http method, the route, and the HTTP Status
func (l *Logger) InfoHTTP(r *http.Request, status int) {
	l.Info(fmt.Sprintf("%v %v %v", r.Method, r.RequestURI, status))
}

// ApplicationError makes two calls: It first calls HTTPLogger with error 500, it then calls Error to log the error
func (l *Logger) ApplicationError(r *http.Request, err error) {
	l.InfoHTTP(r, http.StatusInternalServerError)
	l.Error(err)
}

// New creates a new logger. The flags arguement specifies which kinds of data are written to.
func New(out io.Writer, flags int) *Logger {
	logger := &Logger{}

	logger.err = log.New(out, "ERROR | ", log.LstdFlags|log.Llongfile)
	logger.info = log.New(out, "INFO | ", log.LstdFlags)
	logger.debug = log.New(out, "DEBUG | ", log.LstdFlags)
	logger.warn = log.New(out, "WARNING | ", log.LstdFlags)
	logger.Flags = flags

	return logger
}
