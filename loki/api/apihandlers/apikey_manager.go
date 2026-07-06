package apihandlers

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"

    dfn "github.com/alphamystic/odin/lib/definers"
    "github.com/alphamystic/odin/lib/utils"
)

// CreateApiKey — Register new encrypted tool credentials (Admin Only)
func (api_hnd *APIHandler) CreateApiKey(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		api_hnd.MethodNotAllowed(res, "POST")
		return
	}

	// Security check: Only admins should handle raw API keys/secrets
	if !api_hnd.IsAdmin(req) {
		api_hnd.Unauthorized(res, "Only admins can register API keys.")
		return
	}

	var api dfn.Api
	if err := json.NewDecoder(req.Body).Decode(&api); err != nil {
		api_hnd.BadRequest(res, "Invalid JSON body.")
		return
	}
    ud := req.Context().Value("userData").(*UserData)
    if ud == nil {
      api_hnd.Unauthorized(res, "User data not found")
      return
    }
    api.OwnerID = ud.UserId

	// Validation: Ensure mandatory fields are present
	if !utils.CheckifStringIsEmpty(api.OwnerID) {
		api_hnd.BadRequest(res, "OwnerID is required.")
		return
	}
	if !utils.CheckifStringIsEmpty(api.ToolName) {
		api_hnd.BadRequest(res, "ToolName is required.")
		return
	}

	// Metadata Assignment
	api.KeyID = utils.GenerateUUID()
	api.Touch() // Assigns timestamps

	ctx := context.Background()
	if err := api_hnd.Dom.CreateApiKey(ctx, api); err != nil {
		utils.Warning(fmt.Sprintf("Error creating ApiKey: %v", err))
		api_hnd.InternalServerError(res, "Failed to store credentials.")
		return
	}

	api_hnd.DynamicResponse(res, map[string]interface{}{
		"status":  "success",
		"message": "Credential stored and encrypted successfully.",
		"data":    api.KeyID,
	})
    return
}

// ListApiKeys — List all tool credentials for a specific owner
func (api_hnd *APIHandler) ListApiKeys(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		api_hnd.MethodNotAllowed(res, "GET")
		return
	}


    ud := req.Context().Value("userData").(*UserData)
    if ud == nil {
      api_hnd.Unauthorized(res, "User data not found")
      return
    }
    ownerID := ud.UserId

	active := req.URL.Query().Get("active") != "false"
	limit := utils.StringToInt(req.URL.Query().Get("limit"))
	if limit <= 0 {
		limit = 20
	}
	offset := utils.StringToInt(req.URL.Query().Get("offset"))

	ctx := req.Context()
	keys, err := api_hnd.Dom.ListApiKeys(ctx, ownerID, active, limit, offset)
	if err != nil {
		utils.Warning(fmt.Sprintf("Error listing ApiKeys for %s: %v", ownerID, err))
		api_hnd.InternalServerError(res, "Failed to retrieve credentials.")
		return
	}

    api_hnd.DynamicResponse(res, map[string]interface{}{
		"status":  "success",
		"message": "User APIKeys rerturned successfully.",
		"data":    keys,
	})
    return
// 	res.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(res).Encode(keys)
}

// ViewApiKey — Retrieve details for a specific key (Decrypted for App Use)
func (api_hnd *APIHandler) ViewApiKey(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		api_hnd.MethodNotAllowed(res, "GET")
		return
	}

	keyID := req.URL.Query().Get("keyid")
	if !utils.CheckifStringIsEmpty(keyID) {
		api_hnd.BadRequest(res, "keyid cannot be empty.")
		return
	}

	ctx := context.Background()
	key, err := api_hnd.Dom.ViewApiKey(ctx, keyID)
	if err != nil {
		utils.Warning(fmt.Sprintf("Error viewing ApiKey %s: %v", keyID, err))
		api_hnd.InternalServerError(res, "Failed to fetch credential details.")
		return
	}

	api_hnd.DynamicResponse(res, map[string]interface{}{
		"status":  "success",
		"message": "Credential retrieved and decrypted.",
		"data":    key,
	})
    return
}

// UpdateApiKey — Modify tool credential fields (Admin Only)
func (api_hnd *APIHandler) UpdateApiKey(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		api_hnd.MethodNotAllowed(res, "POST")
		return
	}

	if !api_hnd.IsAdmin(req) {
		api_hnd.Unauthorized(res, "Only admins can update credentials.")
		return
	}

	var update dfn.Api
	if err := json.NewDecoder(req.Body).Decode(&update); err != nil {
		api_hnd.BadRequest(res, "Invalid JSON body.")
		return
	}

	if !utils.CheckifStringIsEmpty(update.KeyID) {
		api_hnd.BadRequest(res, "KeyID is required for updates.")
		return
	}

	update.Touch() // Refresh the updated_at timestamp

	ctx := context.Background()
	if err := api_hnd.Dom.UpdateKey(ctx, update); err != nil {
		utils.Warning(fmt.Sprintf("Error updating ApiKey: %v", err))
		api_hnd.InternalServerError(res, "Failed to update credentials.")
		return
	}

	api_hnd.DynamicResponse(res, map[string]interface{}{
		"status":  "success",
		"message": "Credentials updated successfully.",
	})
    return
}

// DeactivateApiKey — Set a tool's key to inactive
func (api_hnd *APIHandler) DeactivateApiKey(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		api_hnd.MethodNotAllowed(res, "POST")
		return
	}

	if !api_hnd.IsAdmin(req) {
		api_hnd.Unauthorized(res, "Only admins can deactivate keys.")
		return
	}

	keyID := req.URL.Query().Get("keyid")
	if !utils.CheckifStringIsEmpty(keyID) {
		api_hnd.BadRequest(res, "keyid is required.")
		return
	}

	ctx := context.Background()
	key, err := api_hnd.Dom.ViewApiKey(ctx, keyID)
	if err != nil {
		api_hnd.InternalServerError(res, "Credential not found.")
		return
	}

	key.Active = false
	key.Touch()

	if err := api_hnd.Dom.UpdateKey(ctx, *key); err != nil {
		utils.Warning(fmt.Sprintf("Error deactivating ApiKey: %v", err))
		api_hnd.InternalServerError(res, "Failed to deactivate key.")
		return
	}

	api_hnd.DynamicResponse(res, map[string]interface{}{
		"status":  "success",
		"message": "ApiKey deactivated successfully.",
	})
    return
}

// ValidateApiKey — Server-to-server validation (Check if key matches)
func (api_hnd *APIHandler) ValidateApiKey(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		api_hnd.MethodNotAllowed(res, "POST")
		return
	}

	var body struct {
		ApiKey  string `json:"api_key"`
		OwnerID string `json:"ownerid"`
	}

	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		api_hnd.BadRequest(res, "Invalid body.")
		return
	}

	ctx := context.Background()
	isValid := api_hnd.Dom.CheckIfApiKey(ctx, body.ApiKey, body.OwnerID)

	api_hnd.DynamicResponse(res, map[string]interface{}{
		"status": "success",
		"valid":  isValid,
	})
    return
}