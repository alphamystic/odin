package apihandlers

import (
  "fmt"
  "time"
  "errors"
  "net/http"
  "encoding/json"

  "github.com/alphamystic/odin/lib/utils"
  dfn"github.com/alphamystic/odin/lib/definers"
)


// Package to authenticate a user
// Returns a userdata on success and a failure on error
func (api_hnd *APIHandler) Authenticateuser(res http.ResponseWriter, req *http.Request){
  if req.Method != "POST"{
    api_hnd.MethodNotAllowed(res, "POST")
    return
  }
  ctx := req.Context()
  var reqAuth Auth
  err := json.NewDecoder(req.Body).Decode(&reqAuth)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
  if !utils.IsValidEmail(reqAuth.Email) {
    api_hnd.BadRequest(res,"Invalid email address provided.")
    return
  }
  if !utils.CheckifStringIsEmpty(reqAuth.Password) {
    api_hnd.BadRequest(res,"Password cannot be empty.")
    return
  }

  user, err := api_hnd.Dom.Authenticate(ctx, reqAuth.Password, reqAuth.Email)
  if err != nil {
    utils.Logerror(err)
    if errors.Is(err, dfn.WrongPassword) {
      api_hnd.Unauthorized(res,"Wrong email or password provided.")
      return
    }
    api_hnd.InternalServerError(res,"Internal server error. Please try again later.")
    return
  }

  ud := &UserData{
      UserId:   user.UserID,
      UserName: user.UserName,
      Admin:    user.Admin,
      OwnerID: user.OwnerID,
      Anonymous: user.Anonymous,
  }
  token, err := api_hnd.GenerateJWT(ud)
  if err != nil {
    utils.Danger(err)
    api_hnd.InternalServerError(res,"Internal server error. Please try again later.")
    return
  }

  if api_hnd.AuthenticateServer(req) {
    res.WriteHeader(http.StatusOK)
    json.NewEncoder(res).Encode(AuthResponse{
        Status:            "Success",
        Message:           "Signed in successfully",
        UserAuthorization: token,
    })
    return
  } else {
    http.SetCookie(res, &http.Cookie{
        Name:     "Authorization",
        Value:    token,
        Path:     "/",
        MaxAge:   72000,
        HttpOnly: true,
        Secure:   true,
        SameSite: http.SameSiteLaxMode,
    })
  }
  res.WriteHeader(http.StatusOK)
  return
}


// Authenticate the Server Separetly
func (api_hnd *APIHandler) AuthenticateServer(req *http.Request) bool {
    apiKey := req.Header.Get("WEB-SERVER-API-Key")
    return apiKey == "12345"
}



// Logout user
func (api_hnd *APIHandler) Logout(res http.ResponseWriter, req *http.Request){
  // Check if the Authorization cookie exists
  _, err := req.Cookie("Authorization")
  if err == http.ErrNoCookie {
      api_hnd.Unauthorized(res,"Please login first.")
      return
  } else if err != nil {
      fmt.Println("[+] Some internal error. \nERROR: ", err)
      api_hnd.InternalServerError(res,"Internal error, try again later.")
      return
  }
  // Invalidate the cookie by setting it with a past expiration date
  http.SetCookie(res, &http.Cookie{
      Name:     "Authorization",
      Value:    "",
      Path:     "/",
      Expires:  time.Now().Add(-1 * time.Hour),
      HttpOnly: true,
      Secure:   true,
      SameSite: http.SameSiteLaxMode,
  })
  res.WriteHeader(http.StatusOK)
  json.NewEncoder(res).Encode(map[string]string{"message": "Logged out successfully. ADIOS!!!"})
}

// Change Password/UpdatePassword
func (api_hnd *APIHandler) Updatepassword(res http.ResponseWriter, req *http.Request){
  if req.Method != "POST"{
    api_hnd.MethodNotAllowed(res, "POST")
    return
  }
  return
}
