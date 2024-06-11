package darelog

import (
	"fmt"
	"io"
	"os"
	"log"
	"sync"
)

const (
	INFO = 0x1
	WARN = 0x2
	DEBUG = 0x3
)


type LOG struct {
	mux sync.Mutex
	LogFile *os.File
	LVL int
}

func NewLogger(lvl int) *LOG {
	return &LOG{
		LVL: lvl,
	}
}

func GetEnvLOGLEVEL() int {
	// export LOGLEVEL=[INFO|WARN|DEBUG]
	return GetLOGLEVEL(os.Getenv("LOGLEVEL"))
}

func GetLOGLEVEL(loglvl string) (retval int) {
	switch loglvl {
		case "INFO":
			retval = INFO
		case "WARN":
			retval = WARN
		case "DEBUG":
			retval = DEBUG
	}
	if retval > 0 {
		log.Printf("LOGLEVEL='%s'", loglvl)
	}
	return
}

func (l *LOG) SetLOGLEVEL(lvl int) {
	l.LVL = lvl
}

// SetOutput sets the output writer for the logger.
func (l *LOG) SetOutput(writer io.Writer) {
	log.SetOutput(writer)
}

// OpenLogFile opens a log file for writing.
func (l *LOG) OpenLogFile(filename string) {
	l.mux.Lock()
	defer l.mux.Unlock()
	if l.LogFile != nil {
		return
	}
	logFile, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("Error opening log file:", err)
	}
	l.LogFile = logFile
	l.ConfigureFileAndConsoleOutput()
}

// OpenLogFile opens a log file for writing.
func (l *LOG) CloseLogFile() {
	l.mux.Lock()
	defer l.mux.Unlock()
	if l.LogFile != nil {
		l.LogFile.Close()
	}
}

// ConfigureFileAndConsoleOutput configures the logger to write to both a file and console.
func (l *LOG) ConfigureFileAndConsoleOutput() {
	if l.LogFile == nil {
		return
	}
	writer := io.MultiWriter(os.Stdout, l.LogFile)
	log.SetOutput(writer)
}

// Info logs a message at the info level.
func (l *LOG) Info(format string, args ...any) {
	if l.LVL >= INFO || l.LVL == DEBUG {
		log.Printf("[INFO] "+format, args...)
	}
}

// Warn logs a message at the warn level.
func (l *LOG) Warnx(format string, args ...any) {
	// always print warnings
	log.Printf("[WARN] "+format, args...)
}

// Error logs a message at the error level.
func (l *LOG) Error(format string, args ...any) {
	// always print errors
	log.Printf("[ERROR] "+format, args...)
}

// Debug logs a message at the info level.
func (l *LOG) Debug(format string, args ...any) {
	if l.LVL == DEBUG {
		log.Printf("[DEBUG] "+format, args...)
	}
}

// Fatal logs a message at the fatal level, then exits the program.
func (l *LOG) Fatal(format string, args ...any) {
	log.Fatalf("[FATAL] "+format, args...)
}
