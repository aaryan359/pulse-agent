// pkg/logger/logger.go
package logger

import (
	"fmt"
	"log"
	"os"
	"time"
)

var (
	infoLogger  *log.Logger
	errorLogger *log.Logger
	warnLogger  *log.Logger
	debugLogger *log.Logger
	logLevel    string
)

func Init() {
	logLevel = os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}

	infoLogger = log.New(os.Stdout, "", 0)
	errorLogger = log.New(os.Stderr, "", 0)
	warnLogger = log.New(os.Stdout, "", 0)
	debugLogger = log.New(os.Stdout, "", 0)
}

func formatMessage(level, format string, v ...interface{}) string {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	message := fmt.Sprintf(format, v...)
	return fmt.Sprintf("[%s] [%s] %s", timestamp, level, message)
}

func Info(format string, v ...interface{}) {
	infoLogger.Println(formatMessage("INFO", format, v...))
}

func Error(format string, v ...interface{}) {
	errorLogger.Println(formatMessage("ERROR", format, v...))
}

func Warn(format string, v ...interface{}) {
	warnLogger.Println(formatMessage("WARN", format, v...))
}

func Debug(format string, v ...interface{}) {
	if logLevel == "debug" {
		debugLogger.Println(formatMessage("DEBUG", format, v...))
	}
}

func Fatal(format string, v ...interface{}) {
	errorLogger.Println(formatMessage("FATAL", format, v...))
	os.Exit(1)
}
