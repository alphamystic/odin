package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Config constants
const (
	HostIP   = "https://gcp-wazuh-01.ipsl.co.ke" // Replace with your HOST_IP
	Username = "sodhiambo@ipsl.co.ke"      // Replace with your Username
	Password = "S@muel001!"  // Replace with your Password
	BaseURL  = HostIP//"https://" + HostIP + ":55000"
)

// AuthResponse matches the JSON structure provided
type AuthResponse struct {
	Data struct {
		Token string `json:"token"`
	} `json:"data"`
	Error int `json:"error"`
}

func main() {
	// 1. Setup insecure transport (equivalent to curl -k)
	// Warning: Only use InsecureSkipVerify for development/internal testing.
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Transport: tr,
		Timeout:   time.Second * 10,
	}

	// 2. Authenticate and get JWT Token
	token, err := getJWT(client)
	if err != nil {
		fmt.Printf("Authentication failed: %v\n", err)
		return
	}
	fmt.Println("Successfully authenticated. Token acquired.")

	// 3. Perform the "Ping" (Example endpoint: /manager/status)
	// Replace "/manager/status" with your specific ping endpoint
	pingEndpoint := "/manager/status"
	err = performSecureRequest(client, "GET", pingEndpoint, token)
	if err != nil {
		fmt.Printf("Ping request failed: %v\n", err)
	}
}

func getJWT(client *http.Client) (string, error) {
	authURL := BaseURL + "/security/user/authenticate"

	req, err := http.NewRequest("POST", authURL, nil)
	if err != nil {
		return "", err
	}

	// Set Basic Auth header
	req.SetBasicAuth(Username, Password)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("auth failed with status %d: %s", resp.StatusCode, string(body))
	}

	var authRes AuthResponse
	if err := json.Unmarshal(body, &authRes); err != nil {
		return "", err
	}

	return authRes.Data.Token, nil
}

func performSecureRequest(client *http.Client, method, endpoint, token string) error {
	url := BaseURL + endpoint

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return err
	}

	// Add the JWT to the Authorization header
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("Response from %s:\n%s\n", endpoint, string(body))

	return nil
}