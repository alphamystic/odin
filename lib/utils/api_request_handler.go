package utils

import (
	"io"
	"log"
	"fmt"
	"sync"
	"bytes"
	"errors"
	"net/url"
	"net/http"
	"encoding/json"
)

// OdinAPIClient manages authentication and API requests.
type OdinAPIClient struct {
	client   *http.Client
	cookie   string
	loggedIn bool
	baseURL  string
	mu       sync.Mutex
}

// APIResponse represents a structured JSON response.
type APIResponse struct {
	Status  string      `json:"status"`
  RedirectUrl string `json: "redirecturl"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Singleton instance
var instance *OdinAPIClient
var once sync.Once

// GetClient returns the singleton instance of OdinAPIClient.
func GetClient(baseURL string) *OdinAPIClient {
	once.Do(func() {
		instance = &OdinAPIClient{
			client:  &http.Client{},
			baseURL: baseURL,
		}
	})
	return instance
}

// Login authenticates and stores the session cookie.
func (o *OdinAPIClient) Login(email, password string) error {
	o.mu.Lock()
	defer o.mu.Unlock()

	loginURL := fmt.Sprintf("%s/api/login", o.baseURL)
	payload := map[string]string{"email": email, "password": password}
	jsonData, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", loginURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := o.client.Do(req)
	if err != nil {
		return fmt.Errorf("error making login request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("login failed with status: %s", resp.Status)
	}

	// Extract the authentication cookie
	for _, cookie := range resp.Cookies() {
		if cookie.Name == "Authorization" {
			o.cookie = cookie.Value
			o.loggedIn = true
			break
		}
	}

	if o.cookie == "" {
		return errors.New("authorization cookie not found")
	}

	fmt.Println("Login successful. Session stored.")
	return nil
}

// IsLoggedIn checks if the client is authenticated.
func (o *OdinAPIClient) IsLoggedIn() bool {
	o.mu.Lock()
	defer o.mu.Unlock()
	return o.loggedIn
}


func (o *OdinAPIClient) DoRequest(method,endpoint string,payload interface{}) (*APIResponse,error){
	o.mu.Lock()
	defer o.mu.Unlock()
	if !o.loggedIn {
		return nil, errors.New("client is not logged in")
	}
	fullURL, err := url.JoinPath(o.baseURL, endpoint)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}
	var reqBody io.Reader
	if payload != nil {
		reqBody = bytes.NewBufferString(payload.(string))
	}
	req, err := http.NewRequest(method, fullURL, reqBody)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cookie", "Authorization="+o.cookie)
	resp, err := o.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		Notice(fmt.Sprintf("Status code returned: %s",resp.StatusCode))
		return nil,fmt.Errorf("%v", fmt.Sprintf("%s",resp.Body))
	}
	// Decode JSON response
	var apiResp APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &apiResp, nil
}

// makeRequest is a helper function for handling API requests.
func (o *OdinAPIClient) makeRequest(method, endpoint string, body interface{}) (*APIResponse, error) {
	o.mu.Lock()
	defer o.mu.Unlock()

	if !o.loggedIn {
		return nil, errors.New("client is not logged in")
	}

	fullURL, err := url.JoinPath(o.baseURL, endpoint)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	var reqBody io.Reader
	if body != nil {
		jsonData, _ := json.Marshal(body)
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, fullURL, reqBody)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cookie", "Authorization="+o.cookie)

	resp, err := o.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()
	body, err = io.ReadAll(resp.Body)
	resp.Body.Close()
	if resp.StatusCode > 299 {
		log.Fatalf("Response failed with status code: %d and\nbody: %s\n", resp.StatusCode, body)
	}
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s", body)

	// Decode JSON response
	var apiResp APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &apiResp, nil
}

// Get sends a GET request.
func (o *OdinAPIClient) Get(endpoint string) (*APIResponse, error) {
	return o.makeRequest("GET", endpoint, nil)
}

// Post sends a POST request.
func (o *OdinAPIClient) Post(endpoint string, body interface{}) (*APIResponse, error) {
	fmt.Println("Doing post for scans")
	fmt.Println(body)
	return o.makeRequest("POST", endpoint, body)
}

// Put sends a PUT request.
func (o *OdinAPIClient) Put(endpoint string, body interface{}) (*APIResponse, error) {
	return o.makeRequest("PUT", endpoint, body)
}

// Delete sends a DELETE request.
func (o *OdinAPIClient) Delete(endpoint string, body interface{}) (*APIResponse, error) {
	return o.makeRequest("DELETE", endpoint, body)
}
