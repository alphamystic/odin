package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func main() {
	if len(os.Args) < 6 {
		fmt.Println("Usage: go run add_server.go <type> <url> <name> <max_connections> <lb_url>")
		fmt.Println("Example: go run add_server.go API http://localhost:4002 API-3 100 https://localhost:8080")
		os.Exit(1)
	}

	payload := map[string]interface{}{
		"type":            os.Args[1],
		"url":             os.Args[2],
		"name":            os.Args[3],
		"max_connections": atoi(os.Args[4]),
	}

	data, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Failed to marshal JSON:", err)
		os.Exit(1)
	}

	lbURL := os.Args[5] + "/register-backend"
	req, err := http.NewRequest("POST", lbURL, bytes.NewBuffer(data))
	if err != nil {
		fmt.Println("Failed to create request:", err)
		os.Exit(1)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Failed to send request:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	fmt.Println("Response status:", resp.Status)
}

func atoi(s string) int {
	n, err := fmt.Sscanf(s, "%d", new(int))
	if err != nil || n != 1 {
		fmt.Println("Invalid number for max_connections:", s)
		os.Exit(1)
	}
	var val int
	fmt.Sscanf(s, "%d", &val)
	return val
}
