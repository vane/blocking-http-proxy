package proxy

import (
	"fmt"
	"log"
	"os"
)

type Logger struct{

}

type FileLogger struct {
	Logger
	filename string
}

func NewLogger() *Logger {
	return &Logger{}
}

func NewFileLogger(filename string) *FileLogger {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
	}
	return &FileLogger{filename: filename}
}

func (l *FileLogger) WriteLine(msg string, print bool, v ...interface{}) {
	f, err := os.OpenFile(l.filename, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}
	line := fmt.Sprintf(msg, v...)
	if print {
		log.Printf(msg, v...)
	}
	f.WriteString(line)
	defer f.Close()
}

func (l *Logger) Printf(msg string, v ...interface{}) {
	log.Printf(msg, v...)
}
