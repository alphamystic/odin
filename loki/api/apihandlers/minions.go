package apihandlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/alphamystic/odin/lib/utils"
	dfn "github.com/alphamystic/odin/lib/definers"
)

// ==============================
// CREATE MINION
// ==============================
func (api_hnd *APIHandler) CreateMinion(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		api_hnd.MethodNotAllowed(res, "POST")
		return
	}

	var minion dfn.Minion
	if err := json.NewDecoder(req.Body).Decode(&minion); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	if utils.CheckifStringIsEmpty(minion.OwnerID) {
		api_hnd.BadRequest(res, "OwnerID is required.")
		return
	}
	if utils.CheckifStringIsEmpty(minion.MothershipID) {
		api_hnd.BadRequest(res, "MothershipID is required.")
		return
	}

	// Assign UUID and timestamps
	minion.MinionID = utils.GenerateUUID()
	minion.Touch()

	ctx := context.Background()
	if err := api_hnd.Dom.CreateMinion(ctx, minion); err != nil {
		utils.Warning(fmt.Sprintf("Error creating minion: %v", err))
		api_hnd.InternalServerError(res, "Failed to create minion.")
		return
	}

	api_hnd.DynamicResponse(res, map[string]interface{}{
		"status":  "success",
		"message": "Minion registered successfully.",
		"data":    minion.MinionID,
	})
	return
}

// ==============================
// LIST MINIONS
// ==============================
func (api_hnd *APIHandler) ListMinions(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		api_hnd.MethodNotAllowed(res, "GET")
		return
	}

	ownerID := req.URL.Query().Get("ownerid")
	if utils.CheckifStringIsEmpty(ownerID) {
		api_hnd.BadRequest(res, "OwnerID is required.")
		return
	}

	// mothershipID := req.URL.Query().Get("mothershipid")
	// activeStr := req.URL.Query().Get("active")
	// hardwareStr := req.URL.Query().Get("hardware")
	limitStr := req.URL.Query().Get("limit")
	offsetStr := req.URL.Query().Get("offset")

	active := true
	if activeStr == "false" {
		active = false
	}

	hardware := false
	if hardwareStr == "true" {
		hardware = true
	}

	limit := utils.StringToInt(limitStr)
	offset := utils.StringToInt(offsetStr)
	if limit <= 0 {
		limit = 20
	}

	ctx := context.Background()
	minions, err := api_hnd.Dom.ListAllMinions(ctx, limit, offset)
	if err != nil {
		utils.Warning(fmt.Sprintf("Error listing minions: %v", err))
		api_hnd.InternalServerError(res, "Failed to list minions.")
		return
	}

	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(minions)
	return
}

// ==============================
// VIEW MINION
// ==============================
func (api_hnd *APIHandler) ViewMinion(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		api_hnd.MethodNotAllowed(res, "GET")
		return
	}

	minionID := req.URL.Query().Get("minionid")
	if utils.CheckifStringIsEmpty(minionID) {
		api_hnd.BadRequest(res, "MinionID is required.")
		return
	}

	ctx := context.Background()
	minion, err := api_hnd.Dom.ViewMinion(ctx, minionID)
	if err != nil {
		utils.Warning(fmt.Sprintf("Error viewing minion: %v", err))
		api_hnd.InternalServerError(res, "Failed to fetch minion details.")
		return
	}

	api_hnd.DynamicResponse(res, map[string]interface{}{
		"status":  "success",
		"message": "Minion fetched successfully.",
		"data":    minion,
	})
	return
}

// ==============================
// UPDATE MINION
// ==============================
func (api_hnd *APIHandler) UpdateMinion(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		api_hnd.MethodNotAllowed(res, "POST")
		return
	}

	if !api_hnd.IsAdmin(req) {
		api_hnd.Unauthorized(res, "Only admins can update minions.")
		return
	}

	var update dfn.Minion
	if err := json.NewDecoder(req.Body).Decode(&update); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	if utils.CheckifStringIsEmpty(update.MinionID) {
		api_hnd.BadRequest(res, "MinionID is required.")
		return
	}

	ctx := context.Background()
	if err := api_hnd.Dom.UpdateMinion(ctx, update); err != nil {
		utils.Warning(fmt.Sprintf("Error updating minion: %v", err))
		api_hnd.InternalServerError(res, "Failed to update minion.")
		return
	}

	api_hnd.DynamicResponse(res, map[string]interface{}{
		"status":  "success",
		"message": "Minion updated successfully.",
	})
	return
}

// ==============================
// DEACTIVATE MINION
// ==============================
func (api_hnd *APIHandler) DeactivateMinion(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		api_hnd.MethodNotAllowed(res, "POST")
		return
	}

	// if !api_hnd.IsAdmin(req) {
	// 	api_hnd.Unauthorized(res, "Only admins can deactivate minions.")
	// 	return
	// }
	//
	// minionID := req.URL.Query().Get("minionid")
	// if utils.CheckifStringIsEmpty(minionID) {
	// 	api_hnd.BadRequest(res, "MinionID is required.")
	// 	return
	// }
	//
	// ctx := context.Background()
	// minion, err := api_hnd.Dom.ViewMinion(ctx, minionID)
	// if err != nil {
	// 	api_hnd.InternalServerError(res, "Minion not found or cannot be retrieved.")
	// 	return
	// }
	//
	// minion.Active = false
	// if err := api_hnd.Dom.UpdateMinion(ctx, *minion); err != nil {
	// 	utils.Warning(fmt.Sprintf("Error deactivating minion: %v", err))
	// 	api_hnd.InternalServerError(res, "Failed to deactivate minion.")
	// 	return
	// }

	api_hnd.DynamicResponse(res, map[string]interface{}{
		"status":  "success",
		"message": "Minion deactivated successfully.",
	})
	return
}
