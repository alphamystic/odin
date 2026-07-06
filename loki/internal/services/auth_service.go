package services

import (
	"bytes"
	"context"
	"encoding/json"
	"net/url"
	"crypto/tls"
	"fmt"
	"io"
	"strings"
	"net/http"

	dfn "github.com/alphamystic/odin/lib/definers"
)

type (
	Authorize interface {
		Login() (*dfn.User, error)
		Logout(string) error
	}
	AuthorizeService struct {
		SAC *ServerAPIConnector
	}

	// LoginResponse matches the specific JSON keys returned by your API
	LoginResponse struct {
		Status  string `json:"status"`  // Matches "status"
		Message string `json:"message"` // Matches "message"
		Token   string `json:"coockie"` // Matches the typo "coockie" from your API
	}
)

func NewAuthorizeService(sac *ServerAPIConnector) *AuthorizeService {
	return &AuthorizeService{SAC: sac}
}

func (as *AuthorizeService) Login(ctx context.Context, pass, email string) (string, error) {
    // 1. Safe URL Construction
    base, err := url.Parse(as.SAC.BaseURL)
    if err != nil {
        return "", fmt.Errorf("invalid base URL: %w", err)
    }
    rel, err := url.Parse("/api/login")
    if err != nil {
        return "", fmt.Errorf("invalid endpoint: %w", err)
    }
    fullURL := base.ResolveReference(rel).String()

    // 2. Prepare Payload
    payload := map[string]string{"email": email, "password": pass}
    jsonData, err := json.Marshal(payload)
    if err != nil {
       return "", err
    }

    // 3. Configure HTTPS Client to skip TLS verification
    client := &http.Client{
        Transport: &http.Transport{
            TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
        },
    }

    req, err := http.NewRequestWithContext(ctx, "POST", fullURL, bytes.NewBuffer(jsonData))
    if err != nil {
       return "", fmt.Errorf("error creating request: %w", err)
    }

    // 4. Set Headers
    req.Header.Set("WEB-SERVER-API-Key", as.SAC.apiKey)
    req.Header.Set("Content-Type", "application/json")

    // 5. Execute Request
    resp, err := client.Do(req)
    if err != nil {
       return "", fmt.Errorf("error making login request: %w", err)
    }
    defer resp.Body.Close()

    // 6. Handle Specific Status Codes
    if resp.StatusCode != http.StatusOK {
       if resp.StatusCode == http.StatusUnauthorized {
          return "", dfn.WrongPassword
       }
       return "", fmt.Errorf("login failed with status: %s", resp.Status)
    }

    body, err := io.ReadAll(resp.Body)
    if err != nil {
       return "", fmt.Errorf("error reading response body: %w", err)
    }

    // 7. Unmarshal directly into LoginResponse
    var result LoginResponse
    if err := json.Unmarshal(body, &result); err != nil {
       return "", fmt.Errorf("error decoding response: %w. Raw: %s", err, string(body))
    }

    // 8. Validate Business Logic Status
    status := strings.ToLower(result.Status)
    if status != "success" && status != "ok" {
       return "", fmt.Errorf("API returned error status [%s]: %s", result.Status, result.Message)
    }

    // 9. Ensure the token exists
    if result.Token == "" {
       return "", fmt.Errorf("API indicated success but token is empty. Raw: %s", string(body))
    }

    return result.Token, nil
}
func (as *AuthorizeService) ChangePassword() error { return nil }
func (as *AuthorizeService) Logout(token string) error { return nil }