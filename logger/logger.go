package logger

import (
	"fmt"
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

var log *logrus.Logger
var logFile *os.File

type customFormatter struct {
	logrus.TextFormatter
}

func init() {
	log = logrus.New()
	formatter := &Formatter{
		TimestampFormat: "2006-01-02 15:04:05",
		LogFormat:       "%time% [%lvl%] - %msg%\n",
	}
	log.Formatter = formatter
}

// SetOutput sets the output writer for the logger.
func SetOutput(writer io.Writer) {
	log.SetOutput(writer)
}

// OpenLogFile opens a log file for writing.
func OpenLogFile(filename string) {
	logFile, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("Error opening log file:", err)
	}
	ConfigureFileAndConsoleOutput(logFile)
}

// OpenLogFile opens a log file for writing.
func CloseLogFile() {
	logFile.Close()
}

// ConfigureFileAndConsoleOutput configures the logger to write to both a file and console.
func ConfigureFileAndConsoleOutput(logFile *os.File) {
	writer := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(writer)
}

// Debug logs a message at the debug level.
func Debug(args ...interface{}) {
	log.Debug(args...)
}

// Info logs a message at the info level.
func Info(args ...interface{}) {
	log.Info(args...)
}

// Warn logs a message at the warn level.
func Warn(args ...interface{}) {
	log.Warn(args...)
}

// Error logs a message at the error level.
func Error(args ...interface{}) {
	log.Error(args...)
}

// Fatal logs a message at the fatal level, then exits the program.
func Fatal(args ...interface{}) {
	//log.Fatal(args...)
	var formattedArgs []interface{}
	for _, arg := range args {
		formattedArgs = append(formattedArgs, fmt.Sprintf("%v", arg))
	}
	// Convert each string to interface{} before passing
	log.Fatal(append([]interface{}{}, formattedArgs...)...)
}
