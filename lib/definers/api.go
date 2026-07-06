package definers

import (
  "github.com/alphamystic/odin/lib/utils"
)

type Api struct {
    KeyID          string `json:"key_id"`
    ToolName       string `json:"tool_name"`
    OwnerID        string `json:"ownerid"`
    CredentialType string `json:"credential_type"` // apikey, oauth2, token
    Data           SecretData `json:"data"`        // Working data (decrypted)
    EncryptedBlob  string     `json:"-"`           // Encrypted string for DB
    Active         bool       `json:"active"`
    utils.TimeStamps
}

type SecretData struct {
    KeyID        string `json:"key_id"`
    ApiKey       string `json:"api_key,omitempty"`
    ClientID     string `json:"client_id,omitempty"`
    ClientSecret string `json:"client_secret,omitempty"`
    Token        string `json:"token,omitempty"`
}
