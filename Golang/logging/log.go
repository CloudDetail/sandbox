package logging

import (
	"log"
	"os"
	"sync"
	"time"

	"fmt"
)

type Level int

const (
	LevelInfo Level = iota
	LevelWarn
	LevelError
)

var (
	mu       sync.Mutex
	logger   = log.New(os.Stdout, "", 0)
	logLevel = LevelInfo
)

func SetLevel(l Level) {
	mu.Lock()
	defer mu.Unlock()
	logLevel = l
}

func formatLog(level Level, msg string) string {
	t := time.Now().Format("2006-01-02 15:04:05")
	levelStr := ""
	switch level {
	case LevelInfo:
		levelStr = "INFO"
	case LevelWarn:
		levelStr = "WARN"
	case LevelError:
		levelStr = "ERROR"
	}
	return "[" + t + "][" + levelStr + "] " + msg
}

func Info(format string, args ...interface{}) {
	mu.Lock()
	defer mu.Unlock()
	if logLevel <= LevelInfo {
		logger.Println(formatLog(LevelInfo, fmt.Sprintf(format, args...)))
	}
}

func Warn(format string, args ...interface{}) {
	mu.Lock()
	defer mu.Unlock()
	if logLevel <= LevelWarn {
		logger.Println(formatLog(LevelWarn, fmt.Sprintf(format, args...)))
	}
}

func Error(format string, args ...interface{}) {
	mu.Lock()
	defer mu.Unlock()
	if logLevel <= LevelError {
		logger.Println(formatLog(LevelError, fmt.Sprintf(format, args...)))
	}
}
