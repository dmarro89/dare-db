package logger

import (
	"fmt"
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

type Logger interface {
	Start(filename string)
	Close()
	Info(args ...interface{})
	Warn(args ...interface{})
	Debug(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
}

type DareLogger struct {
	logger *logrus.Logger
	file   *os.File
}

func NewDareLogger() Logger {
	log := logrus.New()
	log.SetFormatter(&Formatter{
		TimestampFormat: "2006-01-02 15:04:05",
		LogFormat:       "%time% [%lvl%] - %msg%\n",
	})
	return &DareLogger{logger: log}
}

// Start opens a log file for writing and config output
func (dareLogger *DareLogger) Start(filename string) {
	logFile, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("Error opening log file:", err)
	}
	dareLogger.file = logFile
	dareLogger.logger.SetOutput(io.MultiWriter(os.Stdout, logFile))
}

// Close close the log file.
func (dareLogger *DareLogger) Close() {
	if dareLogger.file != nil {
		err := dareLogger.file.Close()
		if err != nil {
			fmt.Println("error closing log file: %w", err)
		}
		dareLogger.file = nil
	}
}

// Debug logs a message at the debug level.
func (dareLogger *DareLogger) Debug(args ...interface{}) {
	dareLogger.logger.Debug(args...)
}

// Info logs a message at the info level.
func (dareLogger *DareLogger) Info(args ...interface{}) {
	dareLogger.logger.Info(args...)
}

// Warn logs a message at the warn level.
func (dareLogger *DareLogger) Warn(args ...interface{}) {
	dareLogger.logger.Warn(args...)
}

// Error logs a message at the error level.
func (dareLogger *DareLogger) Error(args ...interface{}) {
	dareLogger.logger.Error(args...)
}

// Fa tal logs a message at the fatal level, then exits the program.
func (dareLogger *DareLogger) Fatal(args ...interface{}) {
	dareLogger.logger.Fatal(args...)
}
