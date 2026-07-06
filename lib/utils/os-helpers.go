package utils

import (
"os"
	"fmt"
	"os/exec"
	"path/filepath"
	"runtime"
)

// GetCurrentOS returns the string representation of the current operating system (e.g., "windows", "linux", "darwin")
func GetCurrentOS() string {
	return runtime.GOOS
}

// RunExecutable executes a file path or binary command asynchronously across multiple operating systems.
// It formats the path dynamically and hands off execution to the native shell interpreter, returning the Process ID (PID).
func UniversalRunExecutable(targetPath string) (int, error) {
	var cmd *exec.Cmd

	// Clean and normalize the path configuration for the current host OS filesystem rules
	cleanPath := filepath.Clean(targetPath)

	switch runtime.GOOS {
	case "windows":
		// Windows: Execute via the command interpreter environment
		cmd = exec.Command("cmd", "/C", cleanPath)

	case "linux", "darwin", "freebsd":
		// Unix-like systems: Check if we are running a direct local path string
		// Ensure execution permission hooks or explicit shell bindings are handled gracefully
		cmd = exec.Command("/bin/sh", "-c", cleanPath)

	default:
		// Fallback mechanism for rare runtime environments
		return 0, fmt.Errorf("unsupported platform environment target: %s", runtime.GOOS)
	}

	// Start the process asynchronously (does not block waiting for execution completion)
	if err := cmd.Start(); err != nil {
		return 0, fmt.Errorf("failed to launch asynchronous process: %w", err)
	}

	// Return the successfully created system Process identifier
	return cmd.Process.Pid, nil
}


// KillExec inspects the platform type and terminates the specified Process ID (PID)
func KillExec(pid int) error {
	proc, err := os.FindProcess(pid)
	if err != nil {
		return fmt.Errorf("unable to locate process with PID %d: %w", pid, err)
	}

	// Windows doesn't handle os.Kill natively in some situations; use standard platform procedures
	if runtime.GOOS == "windows" {
		return proc.Kill()
	}

	// Unix targets (Linux, macOS) respond natively to standard SIGKILL signals
	return proc.Signal(os.Kill)
}