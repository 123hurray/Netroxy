/*
 * The MIT License (MIT)
 *
 * Copyright (c) 2016 Ray Zhang
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

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
