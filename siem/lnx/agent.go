// +build linux

package lnx

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

// LogEvent represents a log event with an alert level.
type LogEvent struct {
	Timestamp  time.Time `json:"timestamp"`
	AlertLevel string    `json:"alert_level"`
	EventText  string    `json:"event_text"`
}

// LogTracker keeps track of the last log entry read.
type LogTracker struct {
	LastLogTime time.Time
}

func RunAgent() {
	logTracker := &LogTracker{}
	// get the last logged time
	for {
		// Monitor auditd logs
		auditLogs, err := executeCommand("ausearch", "--input-logs", "--start", logTracker.LastLogTime.Format("2006-01-02T15:04:05"), "--interpret", "--raw")
		if err != nil {
			handleError("Error reading auditd logs:", err)
		}
		processLogs("Audit Logs", auditLogs, "info")

		// Monitor kernel logs (dmesg)
		kernelLogs, err := executeCommand("dmesg")
		if err != nil {
			handleError("Error reading kernel logs:", err)
		}
		processLogs("Kernel Logs", kernelLogs, "info")

		// Update the last log time to the current time
		logTracker.LastLogTime = time.Now()

		// Sleep for a specified duration before checking logs again
		time.Sleep(30 * time.Second) // Adjust the duration as needed
	}
}

func executeCommand(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

func processLogs(logType string, logText string, alertLevel string) {
	// Split logText into individual lines
	logLines := strings.Split(logText, "\n")

	for _, line := range logLines {
		if strings.TrimSpace(line) != "" {
			// Create a LogEvent and send it as JSON
			logEvent := LogEvent{
				Timestamp:  time.Now(),
				AlertLevel: alertLevel,
				EventText:  line,
			}

			logJSON, err := json.Marshal(logEvent)
			if err != nil {
				handleError("Error marshalling log event:", err)
			}

			// Print or send the log event (you can modify this part)
			fmt.Println(string(logJSON))
		}
	}
}

func handleError(message string, err error) {
	// Handle errors (you can customize this part)
	fmt.Println(message, err)
	os.Exit(1)
}
