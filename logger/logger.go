package ilog

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"
)

const (
	INFO    = 0x1
	WARN    = 0x2
	DEBUG   = 0x3
	MAXSIZE = 100 * 1024 * 1024 // 100 MB
)

type LOG struct {
	mux     sync.RWMutex
	LogFile *os.File
	LVL     int
	wrote   int // counts bytes
}

func NewLogger(lvl int) *LOG {
	return &LOG{
		LVL: lvl,
	}
}

func GetEnvLOGLEVEL() int {
	if logstr, ok := os.LookupEnv("LOGLEVEL"); !ok {
		return -1
	}
	logstr := os.Getenv("LOGLEVEL")
	// export LOGLEVEL=[INFO|WARN|DEBUG]
	return GetLOGLEVEL(logstr)
}

func GetLOGLEVEL(loglvl string) (retval int) {
	switch loglvl {
	case "INFO":
		retval = INFO
	case "WARN":
		retval = WARN
	case "DEBUG":
		retval = DEBUG
	default:
		retval = -1
	}
	if retval > 0 {
		log.Printf("LOGLEVEL='%s'", loglvl)
	}
	return
}

// ilog.IfDebug returns true if LOGLEVEL is DEBUG
func (l *LOG) IfDebug() bool {
	l.mux.RLock()
	defer l.mux.RUnlock()
	if l.LVL == DEBUG {
		return true
	}
	return false
}

func (l *LOG) SetLOGLEVEL(lvl int) {
	l.mux.Lock()
	defer l.mux.Unlock()
	l.LVL = lvl
	l.Info("LOGLEVEL=%d", lvl)
	log.Printf("SetLOGLEVEL = %d", lvl)
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
		log.SetOutput(os.Stdout)
		log.Printf("Could not ConfigureFileAndConsoleOutput: LogFile is nil")
		return
	}
	writer := io.MultiWriter(os.Stdout, l.LogFile)
	log.SetOutput(writer)
}

// Info logs a message at the info level.
func (l *LOG) Info(format string, args ...any) {
	l.mux.RLock()
	defer l.mux.RUnlock()
	if l.LVL >= INFO || l.LVL == DEBUG {
		log.Printf("[INFO] "+format, args...)
	}
}

// Warn logs a message at the warn level.
func (l *LOG) Warn(format string, args ...any) {
	l.mux.RLock()
	defer l.mux.RUnlock()
	// always print warnings
	log.Printf("[WARN] "+format, args...)
}

// Error logs a message at the error level.
func (l *LOG) Error(format string, args ...any) {
	l.mux.RLock()
	defer l.mux.RUnlock()
	// always print errors
	log.Printf("[ERROR] "+format, args...)
}

// Debug logs a message at the info level.
func (l *LOG) Debug(format string, args ...any) {
	l.mux.RLock()
	defer l.mux.RUnlock()
	if l.LVL == DEBUG {
		log.Printf("[DEBUG] "+format, args...)
	}
}

// Fatal logs a message at the fatal level, then exits the program.
func (l *LOG) Fatal(format string, args ...any) {
	log.Fatalf("[FATAL] "+format, args...)
}
