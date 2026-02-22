package services

import(
  "io"
  "fmt"
	"bytes"
  "context"
	"net/http"
	"encoding/json"
  //"github.com/alphamystic/odin/lib/utils"
  //dom"github.com/alphamystic/odin/lib/domain"
  dfn"github.com/alphamystic/odin/lib/definers"
)

type (
  Authorize interface {
    Login() (*dfn.User,error)
    Logout(string) error
  }
  AuthorizeService struct{
    //authDomain domain.
    SAC *ServerAPIConnector
  }
)


func NewAuthorizeService(sac *ServerAPIConnector) *AuthorizeService{
  return &AuthorizeService{SAC: sac}
}


func (as *AuthorizeService) Login(ctx context.Context, pass, email string) (string, error) {
	loginURL := fmt.Sprintf("%s/api/login", as.SAC.BaseURL)
	payload := map[string]string{"email": email, "password": pass}
	jsonData, err := json.Marshal(payload)
  if err != nil {
    return "",err
  }

	req, err := http.NewRequest("POST", loginURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}
  req.Header.Set("WEB-SERVER-API-Key", as.SAC.apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error making login request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusUnauthorized {
			return "", dfn.WrongPassword
		}
		return "", fmt.Errorf("login failed with status: %s", resp.Status)
	}

	// Read response body into []byte
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %w", err)
	}

	// Parse JSON response
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("error decoding response: %w", err)
	}

	// Extract fields dynamically
	if status, ok := result["Status"].(string); !ok {
		return "", fmt.Errorf("Status: %s", status)
	}

	if message, ok := result["Message"].(string); !ok {
		return "", fmt.Errorf("Message: %s", message)
	}

	token, ok := result["UserAuthorization"].(string)
	if !ok {
		return "", fmt.Errorf("missing authorization token in response")
	}

	return token, nil
}


func (as *AuthorizeService) ChangePassword() error {
  return nil
}

func (as *AuthorizeService) Logout(token string) error  {
  return nil
}
