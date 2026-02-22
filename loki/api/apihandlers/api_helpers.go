package apihandlers

import (
  "fmt"
  "context"
  "net/http"
  "encoding/json"
  "github.com/dgrijalva/jwt-go"

  "github.com/alphamystic/odin/lib/utils"
  // dom"github.com/alphamystic/odin/lib/domain"
   dfn"github.com/alphamystic/odin/lib/definers"
)




// Helper methods for handling responses
func (api_hnd *APIHandler) MethodNotAllowed(res http.ResponseWriter, allowedMethod string) {
  res.WriteHeader(http.StatusMethodNotAllowed)
  res.Header().Set("Content-Type", "application/json")
  json.NewEncoder(res).Encode(map[string]string{"error": "Only " + allowedMethod + " requests are allowed."})
  return
}


func (api_hnd *APIHandler) Unauthorized(res http.ResponseWriter, message string) {
  res.WriteHeader(http.StatusUnauthorized)
  res.Header().Set("Content-Type", "application/json")
  json.NewEncoder(res).Encode(map[string]string{"error": message})
  return
}





func (api_hnd *APIHandler) WithUserData(next http.HandlerFunc) http.HandlerFunc {
  return func(res http.ResponseWriter, req *http.Request) {
    // --- CORS Headers ---
    res.Header().Set("Access-Control-Allow-Origin", "https://localhost:4040")
    res.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
    res.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
    res.Header().Set("Access-Control-Allow-Credentials", "true")

    // Handle preflight request (no auth required)
    if req.Method == http.MethodOptions {
      res.WriteHeader(http.StatusNoContent)
      return
    }

    // --- Auth Handling ---
    ud, err := api_hnd.GetUDFromRequest(req)
    if err != nil {
      api_hnd.Unauthorized(res, "Unauthorized request: "+err.Error())
      return
    }

    // Attach user data to context
    req = req.WithContext(context.WithValue(req.Context(), "userData", ud))

    next.ServeHTTP(res, req)
  }
}

func (api_hnd *APIHandler) BadRequest(res http.ResponseWriter, message string) {
  res.WriteHeader(http.StatusBadRequest)
  res.Header().Set("Content-Type", "application/json")
  json.NewEncoder(res).Encode(map[string]string{"error": message})
  return
}

func (api_hnd *APIHandler) InternalServerError(res http.ResponseWriter, message string) {
  res.WriteHeader(http.StatusInternalServerError)
  res.Header().Set("Content-Type", "application/json")
  json.NewEncoder(res).Encode(map[string]string{"error": message})
  return
}


func (api_hnd *APIHandler) Posdataresponse(res http.ResponseWriter, response PostDataResponse) {
  res.WriteHeader(http.StatusOK)
  res.Header().Set("Content-Type", "application/json")
  json.NewEncoder(res).Encode(response)
  return
}

func (api_hnd *APIHandler) DynamicResponse(res http.ResponseWriter, response interface{}) {
  res.Header().Set("Content-Type", "application/json")
  // Attempt encoding before writing status
  jsonData, err := json.Marshal(response)
  if err != nil {
    http.Error(res, fmt.Sprintf("error encoding response: %v", err), http.StatusInternalServerError)
    return
  }
  res.WriteHeader(http.StatusOK)
  res.Write(jsonData) // Writing actual JSON response
}



func (api_hnd *APIHandler) Success(res http.ResponseWriter, message string) {
  res.Header().Set("Content-Type", "application/json")
  json.NewEncoder(res).Encode(PostDataResponse{Status: "Success", Message: message})
  return
}



func (api_hnd *APIHandler) IsAdmin(req *http.Request) bool {
  ud, err := api_hnd.GetUDFromRequest(req)
  if  err != nil {
    utils.Danger(err)
    return false
  }
  return ud.Admin
  // if ud.RoleType != "ADMIN" {
  //   return false
  // }
  // ctx := context.Background()
  // user,err := api_hnd.Dom.ViewUser(ud.UUID, ctx)
  // if err != nil {
  //   utils.Danger(err)
  //   return false
  // }
  // if user.RoleType != "ADMIN"{
  //   return false
  // }
  // if user.Status != "active" {
  //   return false
  // }
  // return true
}

// Find a way to encrypt the tokenString
func (hnd *APIHandler) GetUDFromSvrToken(tokenString string) (*UserData,error) {
  // @TODO add functionality to check expiry for a jwt token and save it
  token,err := jwt.Parse(tokenString,func(tkn *jwt.Token)(interface{},error){
    if tkn.Method != jwt.SigningMethodHS256{
      return nil,fmt.Errorf("Unexepcted signing method: %v",tkn.Header["alg"])
    }
    return store,nil
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
