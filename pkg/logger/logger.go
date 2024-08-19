package logger

import (
	"fmt"
	"io"
	"log"
	"os"
)

type Logger struct {
	INFO  *log.Logger
	WARN  *log.Logger
	ERROR *log.Logger
	DEBUG *log.Logger
	TRACE *log.Logger
}

func NewLogger(basepath, path string) *Logger {
	l := &Logger{}

	// Ensure the base path exists
	err := os.MkdirAll(basepath, 0755)
	if err != nil {
		fmt.Fprintf(os.Stderr, "logger - NewLogger - os.MkdirAll: %v\n", err)
		return l
	}

	fullpath := basepath + "/" + path
	file, err := os.OpenFile(fullpath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "logger - NewLogger - os.OpenFile: %v\n", err)
		return l
	}

	// Use a multi-writer to write to both file and stdout
	multiWriter := io.MultiWriter(file, os.Stdout)

	l.INFO = log.New(multiWriter, "[INFO]  ", log.Lshortfile|log.LstdFlags)
	l.WARN = log.New(multiWriter, "[WARN]  ", log.Lshortfile|log.LstdFlags)
	l.ERROR = log.New(multiWriter, "[ERROR] ", log.Lshortfile|log.LstdFlags)
	l.DEBUG = log.New(multiWriter, "[DEBUG] ", log.Lshortfile|log.LstdFlags)
	l.TRACE = log.New(multiWriter, "[TRACE] ", log.Lshortfile|log.LstdFlags)

	return l
}
