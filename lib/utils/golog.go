package utils

/*
	* THis package logs a given request and the random data
	* Each day is logged out differently
*/
import (
	"os"
  "io"
  "log"
  "fmt"
	"time"
	"sync"
  "net/http"
	"path/filepath"
)

type RequestLogger struct {
	Dir  string
	Perm os.FileMode
	mu   sync.Mutex
	logger *log.Logger
	file   *os.File
}

// initiated at server startup, opens a fil to log into that way no run race for multiple requests
func NewRequestLogger(dir string, perm os.FileMode) *RequestLogger {
	rl := &RequestLogger{
		Dir:  dir,
		Perm: perm,
	}
	rl.openLogFile()
	return rl
}

// open a log file for writting
func (rl *RequestLogger) openLogFile() {
	currentDate := time.Now().Format("2006-01-02")
	fileName := fmt.Sprintf("%s.log", currentDate)
	filePath := filepath.Join(rl.Dir, fileName)
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, rl.Perm)
	if err != nil {
		log.Printf("[REQLOG-ERROR]  Error opening log file: %v", err)
		return
	}
	rl.file = file
	rl.logger = log.New(io.MultiWriter(rl.file, os.Stdout), "", log.Ldate|log.Ltime)
}

// log a given request
func (rl *RequestLogger) LogRequestDetails(req *http.Request, data string) {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	mt, _ := rl.file.Stat()
	if time.Now().Day() != mt.ModTime().Day() {
		// If it's a new day, open a new log file
		rl.file.Close()
		rl.openLogFile()
	}
	rl.logger.Printf("[REQLOG]  IP: %s | Method: %s | Path: %s  | Data: %s", req.RemoteAddr, req.Method, req.URL.Path, data)
}

// Close should be called to release the resources when done
func (rl *RequestLogger) Close() {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	rl.file.Close()
}
