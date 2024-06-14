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

type ILOG interface {
	LogStart(filename string)
	LogClose()
	Error(format string, args ...any)
	Debug(format string, args ...any)
	Fatal(format string, args ...any)
	Info(format string, args ...any)
	Warn(format string, args ...any)
	GetLOGLEVEL() int
	SetLOGLEVEL(int)
	IfDebug() bool
}

type LOG struct {
	mux      sync.RWMutex
	LogFile  *os.File
	LOGLEVEL int
	wrote    int // counts bytes
}

func NewLogger(lvl int, logfile string) ILOG {
	logs := &LOG{
		LOGLEVEL: lvl,
	}
	if logfile != "" {
		logs.LogStart(logfile)
	}
	return logs
}

func GetEnvLOGLEVEL() int {
	// export LOGLEVEL=[INFO|WARN|DEBUG]
	if logstr, ok := os.LookupEnv("LOGLEVEL"); ok {
		return GetLOGLEVEL(logstr)
	}
	return INFO //default to INFO
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
	if l.LOGLEVEL == DEBUG {
		return true
	}
	return false
}

func (l *LOG) GetLOGLEVEL() int {
	l.mux.RLock()
	l.Info("LOGLEVEL=%d", l.LOGLEVEL)
	l.mux.RUnlock()
	return l.LOGLEVEL
}

func (l *LOG) SetLOGLEVEL(lvl int) {
	l.mux.Lock()
	l.LOGLEVEL = lvl
	l.mux.Unlock()
	l.Info("SetLOGLEVEL LOGLEVEL=%d", lvl)
	log.Printf("SetLOGLEVEL = %d", lvl)
}

// SetOutput sets the output writer for the logs.
func (l *LOG) SetOutput(writer io.Writer) {
	log.SetOutput(writer)
}

// LogStart opens a log file for writing.
func (l *LOG) LogStart(filename string) {
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

// LogStart opens a log file for writing.
func (l *LOG) LogClose() {
	l.mux.Lock()
	defer l.mux.Unlock()
	if l.LogFile != nil {
		l.LogFile.Close()
		l.LogFile = nil
	}
}

// ConfigureFileAndConsoleOutput configures the logs to write to both a file and console.
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
	if l.LOGLEVEL >= INFO || l.LOGLEVEL == DEBUG {
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
	if l.LOGLEVEL == DEBUG {
		log.Printf("[DEBUG] "+format, args...)
	}
}

// Fatal logs a message at the fatal level, then exits the program.
func (l *LOG) Fatal(format string, args ...any) {
	log.Fatalf("[FATAL] "+format, args...)
}
