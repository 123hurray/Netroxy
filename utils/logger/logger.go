package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

var logChan = make(chan func(), 10)
var logFd *os.File
var logLevel int

const LOG_LEVEL_DEBUG = 0
const LOG_LEVEL_INFO = 1
const LOG_LEVEL_WARN = 2
const LOG_LEVEL_ERROR = 3
const LOG_LEVEL_FATAL = 4
const LOG_LEVEL_QUIET = 5

func Start(level int, logFilePath string) {
	var err error
	logLevel = level
	if logFilePath != "" {
		logFd, err = os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			fmt.Println("Failed to open log file ", logFilePath, ".Error:", err)
		}
	}
	go func() {
		if logFd != nil {
			defer logFd.Close()
		}
		for {
			select {
			case run := <-logChan:
				run()
			}
		}
	}()
}

func Debug(v ...interface{}) {
	if logLevel == LOG_LEVEL_DEBUG {
		printf("DEBUG", v...)
	}
}
func Info(v ...interface{}) {
	if logLevel <= LOG_LEVEL_INFO {
		printf("INFO", v...)
	}
}
func Warn(v ...interface{}) {
	if logLevel <= LOG_LEVEL_WARN {
		printf("WARN", v...)
	}
}
func Error(v ...interface{}) {
	if logLevel <= LOG_LEVEL_ERROR {
		printf("ERROR", v...)
	}
}
func Fatal(v ...interface{}) {
	if logLevel <= LOG_LEVEL_FATAL {
		printf("FATAL", v...)
	}
}
func printf(level string, args ...interface{}) {
	_, file, line, ok := runtime.Caller(2)
	if ok == true {
		_, fileName := filepath.Split(file)
		msg := fmt.Sprintf("[%s][%s][%s:%d]:", time.Now().Format("15:04:05.000"), level, fileName, line)
		msg += fmt.Sprintln(args)

		logChan <- func() {
			fmt.Print(msg)
			if logFd != nil {
				logFd.WriteString(msg)
			}
			if level == "FATAL" {
				os.Exit(1)
			}
		}
	}
}
