package apihandlers

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"

    dfn "github.com/alphamystic/odin/lib/definers"
    "github.com/alphamystic/odin/lib/utils"
)

//
// ================================
// CREATE API KEY  (ADMIN ONLY)
// ================================
func (api_hnd *APIHandler) CreateApiKey(res http.ResponseWriter, req *http.Request) {
  if req.Method != http.MethodPost {
    api_hnd.MethodNotAllowed(res, "POST")
    return
  }

  if !api_hnd.IsAdmin(req) {
    api_hnd.Unauthorized(res, "Only admins can create API keys.")
    return
  }

  var body struct {
    ApiKey  string `json:"apikey"`
    OwnerID string `json:"ownerid"`
    Active  bool   `json:"active"`
  }

  if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
    api_hnd.BadRequest(res, "Invalid JSON body.")
    return
  }

  if utils.CheckifStringIsEmpty(body.ApiKey) || utils.CheckifStringIsEmpty(body.OwnerID) {
    api_hnd.BadRequest(res, "apikey and ownerid are required.")
    return
  }

  apiObj := dfn.Api{
    ApiKey:    body.ApiKey,
    OwnerID:   body.OwnerID,
    Active:    body.Active,
  }
  apiObj.Touch()

  ctx := context.Background()
  if err := api_hnd.Dom.CreateApiKey(apiObj, ctx); err != nil {
    utils.Warning(fmt.Sprintf("Error creating apikey: %v", err))
    api_hnd.InternalServerError(res, "Failed to create apikey.")
    return
  }

  api_hnd.DynamicResponse(res, map[string]interface{}{
    "status":  "success",
    "message": "API key created successfully.",
  })
  return
}

//
// ================================
// LIST API KEYS (ADMIN ONLY)
// ================================
func (api_hnd *APIHandler) ListApiKeys(res http.ResponseWriter, req *http.Request) {
  if req.Method != http.MethodGet {
    api_hnd.MethodNotAllowed(res, "GET")
    return
  }

  if !api_hnd.IsAdmin(req) {
    api_hnd.Unauthorized(res, "Only admins can list API keys.")
    return
  }

  active := utils.StringToBool(req.URL.Query().Get("active"))

  ctx := context.Background()
  keys, err := api_hnd.Dom.ListApiKeys(active, ctx)
  if err != nil {
    utils.Warning(fmt.Sprintf("Error listing API keys: %v", err))
    api_hnd.InternalServerError(res, "Failed to list API keys.")
    return
  }

  res.Header().Set("Content-Type", "application/json")
  json.NewEncoder(res).Encode(keys)
  return
}

//
// ================================
// VIEW API KEY (ADMIN ONLY)
// ================================
func (api_hnd *APIHandler) ViewApiKey(res http.ResponseWriter, req *http.Request) {
  if req.Method != http.MethodGet {
    api_hnd.MethodNotAllowed(res, "GET")
    return
  }

  if !api_hnd.IsAdmin(req) {
    api_hnd.Unauthorized(res, "Only admins can view API key details.")
    return
  }

  keyId := req.URL.Query().Get("apikey")
  if utils.CheckifStringIsEmpty(keyId) {
    api_hnd.BadRequest(res, "apikey parameter is required.")
    return
  }

  ctx := context.Background()
  key, err := api_hnd.Dom.ViewApiKey(keyId, ctx)
  if err != nil {
    api_hnd.BadRequest(res, err.Error())
    return
  }

  res.Header().Set("Content-Type", "application/json")
  json.NewEncoder(res).Encode(key)
  return
}

//
// ================================
// UPDATE API KEY (ADMIN ONLY)
// ================================
func (api_hnd *APIHandler) UpdateApiKey(res http.ResponseWriter, req *http.Request) {
  if req.Method != http.MethodPut && req.Method != http.MethodPatch {
    api_hnd.MethodNotAllowed(res, "PUT/PATCH")
    return
  }

  if !api_hnd.IsAdmin(req) {
    api_hnd.Unauthorized(res, "Only admins can update API keys.")
    return
  }

  var body struct {
    OwnerID string `json:"ownerid"`
    ApiKey  string `json:"apikey"`
  }

  if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
    api_hnd.BadRequest(res, "Invalid JSON body.")
    return
  }

  if utils.CheckifStringIsEmpty(body.OwnerID) || utils.CheckifStringIsEmpty(body.ApiKey) {
    api_hnd.BadRequest(res, "ownerid and apikey are required.")
    return
  }

  ctx := context.Background()
  if err := api_hnd.Dom.UpdateKey(body.OwnerID, body.ApiKey, ctx); err != nil {
    utils.Warning(fmt.Sprintf("Error updating API key: %v", err))
    api_hnd.InternalServerError(res, "Failed to update API key.")
    return
  }

  api_hnd.DynamicResponse(res, map[string]interface{}{
    "status":  "success",
    "message": "API key updated.",
  })
  return
}

//
// ================================
// DEACTIVATE API KEY (ADMIN ONLY)
// ================================
func (api_hnd *APIHandler) DeactivateApiKey(res http.ResponseWriter, req *http.Request) {
  if req.Method != http.MethodPost && req.Method != http.MethodPatch {
    api_hnd.MethodNotAllowed(res, "POST/PATCH")
    return
  }

  if !api_hnd.IsAdmin(req) {
    api_hnd.Unauthorized(res, "Only admins can deactivate API keys.")
    return
  }

  var body struct {
    OwnerID string `json:"ownerid"`
    ApiKey  string `json:"apikey"`
  }

  if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
    api_hnd.BadRequest(res, "Invalid JSON body.")
    return
  }

  if utils.CheckifStringIsEmpty(body.OwnerID) || utils.CheckifStringIsEmpty(body.ApiKey) {
    api_hnd.BadRequest(res, "ownerid and apikey are required.")
    return
  }

  ctx := context.Background()
  if err := api_hnd.Dom.DeactivateKey(body.OwnerID, body.ApiKey, ctx); err != nil {
    utils.Warning(fmt.Sprintf("Error deactivating API key: %v", err))
    api_hnd.InternalServerError(res, "Failed to deactivate API key.")
    return
  }

  api_hnd.DynamicResponse(res, map[string]interface{}{
    "status":  "success",
    "message": "API key deactivated.",
  })
  return
}

//
// ================================
// VALIDATE API KEY (Server-to-server)
// ================================
// NOTE: This is NOT admin-only.
//       Used for internal integrations.
// ================================
func (api_hnd *APIHandler) ValidateApiKey(res http.ResponseWriter, req *http.Request) {
  if req.Method != http.MethodPost {
    api_hnd.MethodNotAllowed(res, "POST")
    return
  }

  var body struct {
    ApiKey  string `json:"apikey"`
    OwnerID string `json:"ownerid"`
  }

  if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
    api_hnd.BadRequest(res, "Invalid JSON body.")
    return
  }

  if utils.CheckifStringIsEmpty(body.ApiKey) || utils.CheckifStringIsEmpty(body.OwnerID) {
    api_hnd.BadRequest(res, "apikey and ownerid are required.")
    return
  }

  ctx := context.Background()
  valid := api_hnd.Dom.CheckIfApiKey(body.ApiKey, body.OwnerID, ctx)

  api_hnd.DynamicResponse(res, map[string]interface{}{
    "status":  "success",
    "valid":   valid,
    "message": "API key validation complete.",
  })
  return
}
