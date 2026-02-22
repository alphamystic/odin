package utils

/*
	* This package logs a given request and the random data
	* Each day is logged out differently
*/

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type RequestLogger struct {
	Dir     string
	Perm    os.FileMode
	mu      sync.Mutex
	logger  *log.Logger
	file    *os.File
	once    sync.Once
}

// NewRequestLogger initializes the logger and ensures log files are ready.
func NewRequestLogger(dir string, perm os.FileMode) *RequestLogger {
	rl := &RequestLogger{
		Dir:  dir,
		Perm: perm,
	}

	// Ensure log file is created on startup
	rl.once.Do(func() {
		rl.openLogFile()
	})

	return rl
}

// openLogFile ensures the log file is properly created or reopened.
func (rl *RequestLogger) openLogFile() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	currentDate := time.Now().Format("2006-01-02")
	fileName := fmt.Sprintf("%s.log", currentDate)
	filePath := filepath.Join(rl.Dir, fileName)

	// Attempt to open or create the log file
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, rl.Perm)
	if err != nil {
		log.Printf("[REQLOG-ERROR] Failed to open log file: %v", err)
		return
	}

	// Close previous log file if open
	if rl.file != nil {
		rl.file.Close()
	}

	// Assign new file and logger
	rl.file = file
	rl.logger = log.New(io.MultiWriter(rl.file, os.Stdout), "", log.Ldate|log.Ltime)
}

// LogRequestDetails logs HTTP request details safely.
func (rl *RequestLogger) LogRequestDetails(req *http.Request, data string) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	// Ensure log file is initialized before writing
	if rl.file == nil {
		fmt.Println("Warning: Log file is nil, attempting to reinitialize...")
		rl.openLogFile()
		if rl.file == nil {
			fmt.Println("Error: Failed to initialize log file, dropping log entry")
			return
		}
	}

	// Rotate log file if date changes
	mt, err := rl.file.Stat()
	if err != nil {
		fmt.Println("Warning: Failed to retrieve log file info:", err)
		return
	}

	if time.Now().Day() != mt.ModTime().Day() {
		rl.file.Close()
		rl.openLogFile()
	}

	// Log the request details
	rl.logger.Printf("[REQLOG] IP: %s | Method: %s | Path: %s | Data: %s",
		req.RemoteAddr, req.Method, req.URL.Path, data)
}

// Close releases file resources gracefully.
func (rl *RequestLogger) Close() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.file != nil {
		rl.file.Close()
	}
}
