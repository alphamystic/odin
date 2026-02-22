package apihandlers

import (
  "fmt"
  "log"
  "time"
  //"sync"
  "net/http"
  //"database/sql"
  "encoding/json"
  "github.com/dgrijalva/jwt-go"

  "github.com/alphamystic/odin/lib/utils"
  dom"github.com/alphamystic/odin/lib/domain"
  dfn"github.com/alphamystic/odin/lib/definers"
)

const (
  SighnInKey = "1234567890!#$%^&*()QWERYUIPASDFGHLMNBVCXZqwertplhgfdsazxcvbnm"
)


var (
  store = []byte("Odin-Loki Api")
)

// Health handles health check requests and returns a JSON-encoded status.
func (api_hnd *APIHandler) Health(w http.ResponseWriter, r *http.Request) {
	api_hnd.mu.Lock()
	defer api_hnd.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")

	var resp healthResponse
	if api_hnd.MaintenanceMode {
		w.WriteHeader(http.StatusServiceUnavailable)
		resp = healthResponse{
			Status:  "error",
			Message: "Backend is under maintenance.",
		}
		log.Printf("Health check request received. Responding with 503 Service Unavailable (Maintenance Mode).")
	} else {
		w.WriteHeader(http.StatusOK)
		resp = healthResponse{
			Status:  "ok",
			Message: "Backend is healthy!",
		}
		log.Printf("Health check request received. Responding with 200 OK.")
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Error encoding health response: %v", err)
	}
  return
}

// ToggleMaintenanceModeHandler toggles the maintenance mode and returns the new state as JSON.
func (api_hnd *APIHandler) ToggleMaintenanceModeHandler(w http.ResponseWriter, r *http.Request) {
	api_hnd.mu.Lock()
	defer api_hnd.mu.Unlock()

	api_hnd.MaintenanceMode = !api_hnd.MaintenanceMode
	status := "disabled"
	if api_hnd.MaintenanceMode {
		status = "enabled"
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	resp := toggleResponse{
		MaintenanceMode: api_hnd.MaintenanceMode,
		Status:          status,
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Error encoding toggle response: %v", err)
	}

	log.Printf("Maintenance mode toggled to %s.", status)
  return
}
func NewAPIHandler(shutdownCh chan bool, doneCh chan bool,rl *utils.RequestLogger,mode string) (*APIHandler,error) {
  utils.PrintTextInASpecificColorInBold("white",fmt.Sprintf(" Starting API server at: %s",utils.GetCurrentTime()))
  // create db configurations
  dbConfig := dfn.InitializeConnector(mode)
  dbConn,err := dfn.NewMySQLConnector(dbConfig)
  if err != nil {
    utils.Warning(fmt.Sprintf("Error connecting to the DB. \n[-]   ERROR: %s",err))
    return nil,err
  }
  dom := dom.NewDomain(dbConn,10,5)
  return &APIHandler {
    Dom: dom,
    Dbs: dbConn,
    CanWriteLogs: true,
    ShutdownChan: shutdownCh,
    MaintenanceMode: false,
    DoneChan: doneCh,
    RL: rl,
  },nil
}

func (api_hnd *APIHandler) GenerateJWT(ud *UserData) (string,error) {
  expTime := time.Now().Add(time.Hour * 72)
  token := jwt.NewWithClaims(jwt.SigningMethodHS256,jwt.MapClaims{
    "ud": map[string]interface{}{
      "UserId":   ud.UserId,
      "Admin":    ud.Admin,
      "Username": ud.UserName,
    },
    "exp": expTime.Unix(),
  })
  sighnedToken,err := token.SignedString([]byte(SighnInKey))
  if err != nil {
    return "",fmt.Errorf("Error signing token: %q",err)
  }
  return sighnedToken,nil
}

// Find a way to encrypt the tokenString
func (api_hnd *APIHandler) GetUDFromRequest(req *http.Request) (*UserData,error) {
  var tokenString string
  if api_hnd.AuthenticateServer(req) {
    //Extract token from server header.
    tokenString = req.Header.Get("User-Token")
    if utils.CheckifStringIsEmpty(tokenString){
      return &UserData{},fmt.Errorf("Empty Token string provided by the server.")
    }
  } else{
    cookie,err := req.Cookie("Authorization")
    if err != nil{
      return nil,err
    }
    tokenString = cookie.Value
  }
  // @TODO add functionality to check expiry for a jwt token and save it
  token,err := jwt.Parse(tokenString,func(tkn *jwt.Token)(interface{},error){
    if _, ok := tkn.Method.(*jwt.SigningMethodHMAC); !ok {
      return nil, fmt.Errorf("Unexpected signing method: %v", tkn.Header["alg"])
    }
    return []byte(SighnInKey),nil
  })
  if err != nil {
    return nil,fmt.Errorf("Signing error. %q",err)
  }
  if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
    if exp, ok := claims["exp"].(float64); ok {
      if time.Now().Unix() > int64(exp) {
        return nil, fmt.Errorf("Token expired")
      }
    }

    if runtimeMap, ok := claims["ud"].(map[string]interface{}); ok {
      return &UserData{
        UserId:   runtimeMap["UserId"].(string),
        Admin:    runtimeMap["Admin"].(bool),
        UserName: runtimeMap["Username"].(string),
      }, nil
    }
  }
  return nil,dfn.NoClaimsError
}

// Find a way to encrypt the tokenString: Moved to now GetUDFromRequest
func (api_hnd *APIHandler) GetUDFromToken(req *http.Request) (*UserData,error) {
  cookie,err := req.Cookie("Authorization")
  if err != nil{
    return nil,err
  }
  tokenString := cookie.Value
  // @TODO add functionality to check expiry for a jwt token and save it
  token,err := jwt.Parse(tokenString,func(tkn *jwt.Token)(interface{},error){
    if _, ok := tkn.Method.(*jwt.SigningMethodHMAC); !ok {
      return nil, fmt.Errorf("Unexpected signing method: %v", tkn.Header["alg"])
  }
    return []byte(SighnInKey), nil
  })
  if err != nil {
    return nil,fmt.Errorf("Signing error. %q",err)
  }
  if claims,ok := token.Claims.(jwt.MapClaims); ok &&  token.Valid {
    if runtimeMap,ok := claims["ud"].(map[string]interface{}); ok {
      return &UserData{
        UserId: runtimeMap["UserId"].(string),
        Admin: runtimeMap["Admin"].(bool),
        UserName: runtimeMap["Username"].(string),
      },nil
    }
  }
  return nil,dfn.NoClaimsError
}
