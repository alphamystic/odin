package apihandlers


import (
  "sync"
  "net/http"
  //"encoding/json"
  "database/sql"
  dom"github.com/alphamystic/odin/lib/domain"
  //"github.com/dgrijalva/jwt-go"

  "github.com/alphamystic/odin/lib/utils"
  //dfn"github.com/alphamystic/odin/lib/definers"
)


// This package holds ALL structs for api calls.
type Auth struct {
  Email                  string `json:"email"`
  Password               string `json:"password,omitempty"`
  NewPassword            string `json:"new_password,omitempty"`
  ConfirmPassword        string `json:"confirmPassword,omitempty"`
  UpdatePasswordTokenStr string `json:"updatePasswordToken,omitempty"`
}


type AuthResponse struct {
  Status string `json:"status"`
  Message string `json:"message"`
  UserAuthorization string `json:"coockie"`
  UpdatePasswordTokenStr string `json:"updatePasswordToken,omitempty"`
}


type APIHandler struct {
  UE *utils.Crypter // This is the universal file encrypter and decryptor.
  Dom *dom.Domain
  Store *http.Cookie
  Dbs *sql.DB
  RL *utils.RequestLogger
  CanWriteLogs bool
  MaintenanceMode bool
  mu              sync.Mutex
  ShutdownChan,DoneChan chan bool // channels to write into
  InvalidTokens []string // You would want to have this cached ina caching sytem
}

type ListUsersFilter struct{
  Filter string `json:"filter"`
  FilterType string `json:"filter_type"`
}

type PostDataResponse struct {
	Status  string `json:"status"`
  RedirectUrl string `json: "redirecturl,omitempty"`
	Message string `json:"message"`
  Data interface{} `json:"data,omitempty"`
}

// healthResponse defines the JSON structure for health check responses.
type healthResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// toggleResponse defines the JSON structure for toggle responses.
type toggleResponse struct {
	MaintenanceMode bool   `json:"maintenance_mode"`
	Status          string `json:"status"`
}

type UserData struct {
  UserId string
  UserName string
  Admin bool
  OwnerID string  `json: "ownerid"`
  Anonymous bool  `json: "anonymous"`
}
