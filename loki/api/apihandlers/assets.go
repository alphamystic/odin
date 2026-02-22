package apihandlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/alphamystic/odin/lib/utils"
	dfn "github.com/alphamystic/odin/lib/definers"
)

// CreateAsset — Admins can register new assets
func (api_hnd *APIHandler) CreateAsset(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		api_hnd.MethodNotAllowed(res, "POST")
		return
	}


	var asset dfn.Asset
	if err := json.NewDecoder(req.Body).Decode(&asset); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	if utils.CheckifStringIsEmpty(asset.OwnerID) {
		api_hnd.BadRequest(res, "OwnerID is required.")
		return
	}

	// Assign UUID and timestamps
	asset.AssetID = utils.GenerateUUID()
	asset.Touch()

	ctx := context.Background()
	if err := api_hnd.Dom.CreateAsset(ctx, asset); err != nil {
		utils.Warning(fmt.Sprintf("Error creating asset: %v", err))
		api_hnd.InternalServerError(res, "Failed to create asset.")
		return
	}

	api_hnd.DynamicResponse(res, map[string]interface{}{
		"status":  "success",
		"message": "Asset created successfully.",
		"data":    asset.AssetID,
	})
}

// ListAssets — lists all assets owned by a user with optional filters (active, hardware, pagination)
func (api_hnd *APIHandler) ListAssets(res http.ResponseWriter, req *http.Request) {
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
	hardwareStr := req.URL.Query().Get("hardware")
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
	assets, err := api_hnd.Dom.ListAssetsByFilter(ctx, ownerID, active, hardware, limit, offset)
	if err != nil {
		utils.Warning(fmt.Sprintf("Error listing assets: %v", err))
		api_hnd.InternalServerError(res, "Failed to list assets.")
		return
	}

	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(assets)
	return
}

// ViewAsset — retrieves details for one asset
func (api_hnd *APIHandler) ViewAsset(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		api_hnd.MethodNotAllowed(res, "GET")
		return
	}

	assetID := req.URL.Query().Get("assetid")
	if utils.CheckifStringIsEmpty(assetID) {
		api_hnd.BadRequest(res, "AssetID cannot be empty.")
		return
	}

	ctx := context.Background()
	asset, err := api_hnd.Dom.ViewAsset(ctx, assetID)
	if err != nil {
		utils.Warning(fmt.Sprintf("Error viewing asset: %v", err))
		api_hnd.InternalServerError(res, "Failed to fetch asset details.")
		return
	}

	api_hnd.DynamicResponse(res, map[string]interface{}{
		"status":  "success",
		"message": "Asset fetched successfully.",
		"data":    asset,
	})
	return
}

// UpdateAsset — allows admins to modify asset fields
func (api_hnd *APIHandler) UpdateAsset(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		api_hnd.MethodNotAllowed(res, "POST")
		return
	}

	if !api_hnd.IsAdmin(req) {
		api_hnd.Unauthorized(res, "Only admins can update assets.")
		return
	}

	var update dfn.Asset
	if err := json.NewDecoder(req.Body).Decode(&update); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	if utils.CheckifStringIsEmpty(update.AssetID) {
		api_hnd.BadRequest(res, "AssetID is required.")
		return
	}

	ctx := context.Background()
	if err := api_hnd.Dom.UpdateAsset(ctx, update); err != nil {
		utils.Warning(fmt.Sprintf("Error updating asset: %v", err))
		api_hnd.InternalServerError(res, "Failed to update asset.")
		return
	}

	api_hnd.DynamicResponse(res, map[string]interface{}{
		"status":  "success",
		"message": "Asset updated successfully.",
	})
	return
}

// DeactivateAsset — sets asset active=false
func (api_hnd *APIHandler) DeactivateAsset(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		api_hnd.MethodNotAllowed(res, "POST")
		return
	}

	if !api_hnd.IsAdmin(req) {
		api_hnd.Unauthorized(res, "Only admins can deactivate assets.")
		return
	}

	assetID := req.URL.Query().Get("assetid")
	if utils.CheckifStringIsEmpty(assetID) {
		api_hnd.BadRequest(res, "AssetID is required.")
		return
	}

	// Get the asset first
	ctx := context.Background()
	asset, err := api_hnd.Dom.ViewAsset(ctx, assetID)
	if err != nil {
		api_hnd.InternalServerError(res, "Asset not found or cannot be retrieved.")
		return
	}

	asset.Active = false
	if err := api_hnd.Dom.UpdateAsset(ctx, *asset); err != nil {
		utils.Warning(fmt.Sprintf("Error deactivating asset: %v", err))
		api_hnd.InternalServerError(res, "Failed to deactivate asset.")
		return
	}

	api_hnd.DynamicResponse(res, map[string]interface{}{
		"status":  "success",
		"message": "Asset deactivated successfully.",
	})
	return
}
