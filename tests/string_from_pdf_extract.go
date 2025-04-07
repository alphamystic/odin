package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func extractLinesWithString(inputPath, searchString string) ([]string, error) {
	// Convert the PDF to text using pdftotext command-line tool
	outputPath := inputPath + ".txt"
	cmd := exec.Command("pdftotext", inputPath, outputPath)
	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("failed to convert PDF to text: %v", err)
	}

	// Open the converted text file
	file, err := os.Open(outputPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open text file: %v", err)
	}
	defer file.Close()

	linesWithString := []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, searchString) {
			linesWithString = append(linesWithString, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading lines: %v", err)
	}

	// Clean up the temporary text file
	os.Remove(outputPath)

	return linesWithString, nil
}

func writeLinesToFile(lines []string, outputPath string) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %v", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, line := range lines {
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			return fmt.Errorf("failed to write to output file: %v", err)
		}
	}
	return writer.Flush()
}

func main() {
	if len(os.Args) < 3 {
		fmt.Printf("Usage: go run main.go input.pdf output.txt\n")
		os.Exit(1)
	}

	inputPath := os.Args[1]
	outputPath := os.Args[2]
	searchString := "BSSEY3"

	lines, err := extractLinesWithString(inputPath, searchString)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	err = writeLinesToFile(lines, outputPath)
	if err != nil {
		fmt.Printf("Error writing lines to file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Extracted lines have been written to %s\n", outputPath)
}
