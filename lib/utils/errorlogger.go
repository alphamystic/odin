package utils

import (
	"os"
	"io"
	"log"
	"fmt"
	"time"
	"sync"
	"path/filepath"
)

type ErrorLogger struct {
	Dir      string
	Perm     os.FileMode
	loggers  map[string]*log.Logger
	files    map[string]*os.File
	mu       sync.Mutex
}

type Logger struct {
	Name string
	Text string
}
// NewErrorLogger initializes loggers for the specified file names
func NewErrorLogger(dir string, perm os.FileMode, fileNames []string) *ErrorLogger {
	el := &ErrorLogger{
		Dir:     dir,
		Perm:    perm,
		loggers: make(map[string]*log.Logger),
		files:   make(map[string]*os.File),
	}

	for _, name := range fileNames {
		el.openLogFile(name)
	}

	return el
}

// openLogFile opens or creates a log file for the given name
func (el *ErrorLogger) openLogFile(name string) {
	currentDate := time.Now().Format("2006-01-02")
	fileName := fmt.Sprintf("%s_%s.log", name, currentDate)
	filePath := filepath.Join(el.Dir, fileName)

	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, el.Perm)
	if err != nil {
		log.Printf("[ERRLOG-ERROR] Error opening log file %s: %v", name, err)
		return
	}

	el.files[name] = file
	el.loggers[name] = log.New(io.MultiWriter(file, os.Stdout), fmt.Sprintf("[%s] ", name), log.Ldate|log.Ltime)
}

// LogError logs an error message to the appropriate file based on the logger name
func (el *ErrorLogger) LogError(loggerName, message string) {
	el.mu.Lock()
	defer el.mu.Unlock()

	logger, exists := el.loggers[loggerName]
	if !exists {
		log.Printf("[ERRLOG-ERROR] Logger not found for name: %s", loggerName)
		return
	}

	file := el.files[loggerName]
	if file == nil {
		log.Printf("[ERRLOG-ERROR] File not found for logger: %s", loggerName)
		return
	}

	mt, _ := file.Stat()
	if time.Now().Day() != mt.ModTime().Day() {
		// If it's a new day, rotate the log file
		file.Close()
		el.openLogFile(loggerName)
		logger = el.loggers[loggerName]
	}

	logger.Printf("[ERRLOG] %s", message)
}

func (el *ErrorLogger) LogToFile(logger Logger) {
	el.mu.Lock()
	defer el.mu.Unlock()

	loggerInstance, exists := el.loggers[logger.Name]
	if !exists {
		log.Printf("[ERRLOG-ERROR] Logger not found for name: %s", logger.Name)
		return
	}

	file := el.files[logger.Name]
	if file == nil {
		log.Printf("[ERRLOG-ERROR] File not found for logger: %s", logger.Name)
		return
	}

	mt, _ := file.Stat()
	if time.Now().Day() != mt.ModTime().Day() {
		// If it's a new day, rotate the log file
		file.Close()
		el.openLogFile(logger.Name)
		loggerInstance = el.loggers[logger.Name]
	}

	loggerInstance.Printf("[ERRLOG] %s", logger.Text)
}

// Close releases resources and closes all open files
func (el *ErrorLogger) Close() {
	el.mu.Lock()
	defer el.mu.Unlock()

	for name, file := range el.files {
		if file != nil {
			file.Close()
			delete(el.files, name)
			delete(el.loggers, name)
		}
	}
}
