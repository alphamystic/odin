package apihandlers

import (
	"fmt"
	"net/http"
	"strconv"
	"context"
	"encoding/json"

	"github.com/alphamystic/odin/lib/utils"
	dfn "github.com/alphamystic/odin/lib/definers"
)

// CreateMothership - admin creates a new mothership
func (api_hnd *APIHandler) CreateMothership(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		api_hnd.MethodNotAllowed(res, "POST")
		return
	}

	if !api_hnd.IsAdmin(req) {
		api_hnd.Unauthorized(res, "Only admins can create motherships.")
		return
	}

	var ms dfn.Mothership
	if err := json.NewDecoder(req.Body).Decode(&ms); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	// Basic setup
	ms.MSId = utils.GenerateUUID()
	ms.Touch() // sets CreatedAt / UpdatedAt timestamps

	ctx := context.Background()
	if err := api_hnd.Dom.CreateMothership(ms, ctx); err != nil {
		utils.Warning(fmt.Sprintf("Error creating mothership: %v", err))
		api_hnd.InternalServerError(res, "Failed to create mothership.")
		return
	}

	api_hnd.DynamicResponse(res, map[string]interface{}{
		"status":  "success",
		"message": "Mothership created successfully.",
		"data":    ms.MSId,
	})
	return
}

// ListMotherships - list motherships by owner (paginated)
func (api_hnd *APIHandler) ListMotherships(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		api_hnd.MethodNotAllowed(res, "GET")
		return
	}

	ownerID := req.URL.Query().Get("ownerid")
	if utils.CheckifStringIsEmpty(ownerID) {
		api_hnd.BadRequest(res, "OwnerID is required.")
		return
	}

	activeStr := req.URL.Query().Get("active")
	limitStr := req.URL.Query().Get("limit")
	offsetStr := req.URL.Query().Get("offset")

	active := true
	if activeStr == "false" {
		active = false
	}

	limit := utils.StringToInt(limitStr)
	offset := utils.StringToInt(offsetStr)
	if limit <= 0 {
		limit = 20
	}

	ctx := context.Background()
	motherships, err := api_hnd.Dom.ListMotherships(ownerID, active, limit, offset, ctx)
	if err != nil {
		utils.Warning(fmt.Sprintf("Error listing motherships: %v", err))
		api_hnd.InternalServerError(res, "Failed to list motherships.")
		return
	}

	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(motherships)
	return
}

// ListFilteredMotherships - list by owner with filters (active, online)
func (api_hnd *APIHandler) ListFilteredMotherships(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		api_hnd.MethodNotAllowed(res, "GET")
		return
	}

	ownerID := req.URL.Query().Get("ownerid")
	if utils.CheckifStringIsEmpty(ownerID) {
		api_hnd.BadRequest(res, "OwnerID is required.")
		return
	}

	activeStr := req.URL.Query().Get("active")
	onlineStr := req.URL.Query().Get("online")
	limitStr := req.URL.Query().Get("limit")
	offsetStr := req.URL.Query().Get("offset")

	active := activeStr != "false"
	online := onlineStr == "true"

	limit := utils.StringToInt(limitStr)
	offset := utils.StringToInt(offsetStr)
	if limit <= 0 {
		limit = 20
	}

	ctx := context.Background()
	motherships, err := api_hnd.Dom.ListFilteredMotherships(ownerID, active, online, limit, offset, ctx)
	if err != nil {
		utils.Warning(fmt.Sprintf("Error filtering motherships: %v", err))
		api_hnd.InternalServerError(res, "Failed to filter motherships.")
		return
	}

	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(motherships)
	return
}

// ViewMothership - view single mothership by msid
func (api_hnd *APIHandler) ViewMothership(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		api_hnd.MethodNotAllowed(res, "GET")
		return
	}

	msid := req.URL.Query().Get("msid")
	if utils.CheckifStringIsEmpty(msid) {
		api_hnd.BadRequest(res, "MSID cannot be empty.")
		return
	}

	ctx := context.Background()
	ms, err := api_hnd.Dom.ViewMothership(msid, ctx)
	if err != nil {
		utils.Warning(fmt.Sprintf("Error viewing mothership: %v", err))
		api_hnd.InternalServerError(res, "Failed to fetch mothership details.")
		return
	}

	api_hnd.DynamicResponse(res, map[string]interface{}{
		"status":  "success",
		"message": "Mothership fetched successfully.",
		"data":    ms,
	})
	return
}

// UpdateMothership - updates existing mothership details
func (api_hnd *APIHandler) UpdateMothership(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		api_hnd.MethodNotAllowed(res, "POST")
		return
	}

	if !api_hnd.IsAdmin(req) {
		api_hnd.Unauthorized(res, "Only admins can update motherships.")
		return
	}

	var update dfn.Mothership
	if err := json.NewDecoder(req.Body).Decode(&update); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	if utils.CheckifStringIsEmpty(update.MSId) || utils.CheckifStringIsEmpty(update.OwnerID) {
		api_hnd.BadRequest(res, "MSID and OwnerID are required.")
		return
	}

	ctx := context.Background()
	if err := api_hnd.Dom.UpdateMothership(update.OwnerID, update.MSId, update, ctx); err != nil {
		utils.Warning(fmt.Sprintf("Error updating mothership: %v", err))
		api_hnd.InternalServerError(res, "Failed to update mothership.")
		return
	}

	api_hnd.DynamicResponse(res, map[string]interface{}{
		"status":  "success",
		"message": "Mothership updated successfully.",
	})
	return
}

// DeactivateMothership - marks mothership as inactive
func (api_hnd *APIHandler) DeactivateMothership(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		api_hnd.MethodNotAllowed(res, "POST")
		return
	}

	if !api_hnd.IsAdmin(req) {
		api_hnd.Unauthorized(res, "Only admins can deactivate motherships.")
		return
	}

	msid := req.URL.Query().Get("msid")
	ownerID := req.URL.Query().Get("ownerid")

	if utils.CheckifStringIsEmpty(msid) || utils.CheckifStringIsEmpty(ownerID) {
		api_hnd.BadRequest(res, "MSID and OwnerID are required.")
		return
	}

	ctx := context.Background()
	if err := api_hnd.Dom.DeactivateMothership(msid, ownerID, ctx); err != nil {
		utils.Warning(fmt.Sprintf("Error deactivating mothership: %v", err))
		api_hnd.InternalServerError(res, "Failed to deactivate mothership.")
		return
	}

	api_hnd.DynamicResponse(res, map[string]interface{}{
		"status":  "success",
		"message": "Mothership deactivated successfully.",
	})
	return
}
