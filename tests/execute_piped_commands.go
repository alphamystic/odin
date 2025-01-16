package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// ExecutePipedCommand takes a piped command (e.g., "ls -l | grep .go") and executes it.
func LinuxExecutePipedCommand(pipedCommand string) (string, error) {
	// Split the piped command into individual commands
	commands := strings.Split(pipedCommand, "|")
	if len(commands) < 1 {
		return "", fmt.Errorf("no commands provided")
	}

	// Initialize the input for the first command
	var inputBuffer *bytes.Buffer
	var outputBuffer bytes.Buffer
	var errBuffer bytes.Buffer

	for i, cmd := range commands {
		// Trim spaces and split the command into its executable and arguments
		parts := strings.Fields(strings.TrimSpace(cmd))
		if len(parts) == 0 {
			return "", fmt.Errorf("empty command detected")
		}

		// Create the command
		execCmd := exec.Command(parts[0], parts[1:]...)

		// Set up stdin for the current command
		if i == 0 {
			// First command takes input from the default stdin
			execCmd.Stdin = nil
		} else {
			// Subsequent commands take input from the output of the previous command
			execCmd.Stdin = inputBuffer
		}

		// Set up stdout and stderr
		if i == len(commands)-1 {
			// Last command captures the final output
			execCmd.Stdout = &outputBuffer
		} else {
			// Intermediate commands pipe their output forward
			inputBuffer = &bytes.Buffer{}
			execCmd.Stdout = inputBuffer
		}
		execCmd.Stderr = &errBuffer

		// Run the command
		if err := execCmd.Run(); err != nil {
			return "", fmt.Errorf("error executing command %q: %v\nstderr: %s", cmd, err, errBuffer.String())
		}
	}

	// Return the final output
	return outputBuffer.String(), nil
}

func main() {
	// Example usage
	command := "ps -ef | grep firefox"
	output, err := ExecutePipedCommand(command)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println(output)
	}
}
