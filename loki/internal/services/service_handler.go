package services

import (
	"io"
	"fmt"
	"bytes"
	"net/url"
	"net/http"
	"encoding/json"
)

// ServerAPIConnector manages API requests with API key authentication.
type ServerAPIConnector struct {
	BaseURL string
	apiKey  string
}

// NewServerAPIConnector creates a new instance of ServerAPIConnector.
func NewServerAPIConnector(baseURL, apiKey string) *ServerAPIConnector {
  ///@ TODO
  // create a request to validate this KEY
	return &ServerAPIConnector{
		BaseURL: baseURL,
		apiKey:  apiKey,
	}
}

// AuthenticateServer verifies the request has the correct API key.
func (api *ServerAPIConnector) AuthenticateServer(req *http.Request) bool {
	apiKey := req.Header.Get("WEB-SERVER-API-Key")
	return apiKey == api.apiKey
}

// makeRequest is a helper function for handling API requests.
func (api *ServerAPIConnector) makeRequest(method, endpoint string, body interface{}) (map[string]interface{}, error) {
	fullURL, err := url.JoinPath(api.BaseURL, endpoint)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	var reqBody io.Reader
	if body != nil {
		jsonData, _ := json.Marshal(body)
		reqBody = bytes.NewBuffer(jsonData)
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, fullURL, reqBody)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("WEB-SERVER-API-Key", api.apiKey)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	// Decode JSON response into a generic map
	var apiResp map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return apiResp, nil
}

// Get sends a GET request.
func (api *ServerAPIConnector) Get(endpoint string) (map[string]interface{}, error) {
	return api.makeRequest("GET", endpoint, nil)
}

// Post sends a POST request.
func (api *ServerAPIConnector) Post(endpoint string, body interface{}) (map[string]interface{}, error) {
	return api.makeRequest("POST", endpoint, body)
}

// Put sends a PUT request.
func (api *ServerAPIConnector) Put(endpoint string, body interface{}) (map[string]interface{}, error) {
	return api.makeRequest("PUT", endpoint, body)
}

// Delete sends a DELETE request.
func (api *ServerAPIConnector) Delete(endpoint string, body interface{}) (map[string]interface{}, error) {
	return api.makeRequest("DELETE", endpoint, body)
}
