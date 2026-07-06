package services

import (
	"io"
	"fmt"
	"bytes"
	"net/url"
	"strings"
	"net/http"
	"crypto/tls"
	"encoding/json"

	"github.com/alphamystic/odin/lib/utils"
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

	// client := &http.Client{}
	client := &http.Client{
        Transport: &http.Transport{
            TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
        },
    }
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



// Add these methods to your existing ServerAPIConnector struct
func (api *ServerAPIConnector) makeAuthorizedRequest(method, endpoint, token string, body interface{}) (map[string]interface{}, error) {
      // 1. Handle URL Construction safely
      base, err := url.Parse(api.BaseURL)
      if err != nil {
          return nil, fmt.Errorf("invalid base URL: %w", err)
      }
      rel, err := url.Parse(endpoint)
      if err != nil {
          return nil, fmt.Errorf("invalid endpoint: %w", err)
      }
      fullURL := base.ResolveReference(rel).String()
      utils.Warning(fmt.Sprintf("%s",fullURL))

      var reqBody io.Reader
      if body != nil {
         jsonData, _ := json.Marshal(body)
         reqBody = bytes.NewBuffer(jsonData)
      }

      // 2. CONFIGURE CLIENT TO SKIP TLS VERIFICATION
      client := &http.Client{
          Transport: &http.Transport{
              TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
          },
      }

      req, err := http.NewRequest(method, fullURL, reqBody)
      if err != nil {
         return nil, fmt.Errorf("error creating request: %w", err)
      }

      req.Header.Set("Content-Type", "application/json")

      // Clean and set token
      token = strings.TrimPrefix(token, "Authorization=")
      if token != "" {
         req.Header.Set("Cookie", fmt.Sprintf("Authorization=%s", token))
      }

      resp, err := client.Do(req)
      if err != nil {
         return nil, fmt.Errorf("error making request: %w", err)
      }
      defer resp.Body.Close()

      // 3. Read Body and Unmarshal
      bodyBytes, err := io.ReadAll(resp.Body)
      if err != nil {
          return nil, fmt.Errorf("failed to read response body: %w", err)
      }

      if resp.StatusCode != http.StatusOK {
          return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(bodyBytes))
      }

      var apiResp map[string]interface{}
      if err := json.Unmarshal(bodyBytes, &apiResp); err != nil {
          return nil, fmt.Errorf("unmarshal error: %w. Raw body: %s", err, string(bodyBytes))
      }

      return apiResp, nil
 }

// Helper Wrappers for Authorized calls
func (api *ServerAPIConnector) AuthGet(endpoint, token string) (map[string]interface{}, error) {
	return api.makeAuthorizedRequest("GET", endpoint, token, nil)
}

func (api *ServerAPIConnector) AuthPost(endpoint, token string, body interface{}) (map[string]interface{}, error) {
	return api.makeAuthorizedRequest("POST", endpoint, token, body)
}
// One function to allow for marshaling and typecasting
func (api *ServerAPIConnector) Decode(input interface{}, target interface{}) error {
	// A common way to convert map[string]interface{} to a struct in Go
	data, err := json.Marshal(input)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, target)
}