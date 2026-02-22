package apihandlers

import (
  "fmt"
  "context"
  "net/http"
  "encoding/json"
  "github.com/alphamystic/odin/lib/utils"
  dfn"github.com/alphamystic/odin/lib/definers"
)


func (api_hnd *APIHandler) Createuser(res http.ResponseWriter, req *http.Request){
  if req.Method != "POST"{
    res.WriteHeader(http.StatusMethodNotAllowed)
    res.Header().Set("Content-Type", "application/json")
    json.NewEncoder(res).Encode("Only post requests allowed.")
    return
  }
  var user dfn.User
  ctx := context.Background()
	err := json.NewDecoder(req.Body).Decode(&user)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
  user.UserID = utils.GenerateUUID()
  user.Touch()
  if err := api_hnd.Dom.CreateUser(ctx, user); err != nil{
    response := PostDataResponse{
  		Status:  "Failed",
  		Message: "Server encountered an error creating user, try again later or contact admin. :)",
  	}
    utils.Warning(fmt.Sprintf("%s",err))
    res.WriteHeader(http.StatusInternalServerError)
    res.Header().Set("Content-Type", "application/json")
  	json.NewEncoder(res).Encode(response)
    return
  }
  response := PostDataResponse{
    Status:  "Success",
    Message: "User created successfully",
    RedirectUrl: user.UserID,
  }
  res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(response)
}


type ListMysUserRequest struct {
  OwnerID string `json: "ownerid"`
  Verified bool `json: "verified"`
  Active bool `json: "active"`
}


func (api_hnd *APIHandler) Listusers(res http.ResponseWriter, req *http.Request) {
  if req.Method != http.MethodGet {
    res.WriteHeader(http.StatusMethodNotAllowed)
    res.Header().Set("Content-Type", "application/json")
    json.NewEncoder(res).Encode(map[string]string{"error": "Only GET requests are allowed."})
    return
  }

  var lmur ListMysUserRequest
  err := json.NewDecoder(req.Body).Decode(&lmur)
  if err != nil {
    res.WriteHeader(http.StatusBadRequest)
    res.Header().Set("Content-Type", "application/json")
    json.NewEncoder(res).Encode(map[string]string{"error": "Invalid request body: " + err.Error()})
    return
  }

  var users []dfn.User
  ctx := context.Background()
  users, err = api_hnd.Dom.ListMyUsers(ctx, lmur.OwnerID, lmur.Active, lmur.Verified)
  if err != nil {
    response := map[string]string{
      "status":  "Failed",
      "message": "Server encountered an error listing users. Please try again later or contact the admin.",
    }
    utils.Warning(fmt.Sprintf("%s", err))
    res.WriteHeader(http.StatusInternalServerError)
    res.Header().Set("Content-Type", "application/json")
    json.NewEncoder(res).Encode(response)
    return
  }
  res.Header().Set("Content-Type", "application/json")
  json.NewEncoder(res).Encode(users)
  return
}


func (api_hnd *APIHandler) Listadmins(res http.ResponseWriter, req *http.Request) {
  if req.Method != http.MethodGet {
    res.WriteHeader(http.StatusMethodNotAllowed)
    res.Header().Set("Content-Type", "application/json")
    json.NewEncoder(res).Encode(map[string]string{"error": "Only GET requests are allowed."})
    return
  }

  var lmur ListMysUserRequest
  err := json.NewDecoder(req.Body).Decode(&lmur)
  if err != nil {
    res.WriteHeader(http.StatusBadRequest)
    res.Header().Set("Content-Type", "application/json")
    json.NewEncoder(res).Encode(map[string]string{"error": "Invalid request body: " + err.Error()})
    return
  }

  var users []dfn.User
  ctx := context.Background()
  users, err = api_hnd.Dom.AdminListUsers(lmur.Active, ctx)
  if err != nil {
    response := map[string]string{
      "status":  "Failed",
      "message": "Server encountered an error listing admin users. Please try again later or contact the admin.",
    }
    utils.Warning(fmt.Sprintf("%s", err))
    res.WriteHeader(http.StatusInternalServerError)
    res.Header().Set("Content-Type", "application/json")
    json.NewEncoder(res).Encode(response)
    return
  }
  res.Header().Set("Content-Type", "application/json")
  json.NewEncoder(res).Encode(users)
  return
}
