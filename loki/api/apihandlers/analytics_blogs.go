package apihandlers

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"

    "github.com/alphamystic/odin/lib/utils"
)

// ====================================
// CREATE ENUM TYPE (ADMIN ONLY)
// ====================================
func (api_hnd *APIHandler) CreateEnumType(res http.ResponseWriter, req *http.Request) {
  if req.Method != http.MethodPost {
    api_hnd.MethodNotAllowed(res, "POST")
    return
  }

  if !api_hnd.IsAdmin(req) {
    api_hnd.Unauthorized(res, "Only admins can create enum types.")
    return
  }

  var body struct {
    Name        string `json:"name"`
    Description string `json:"description"`
  }

  if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
    api_hnd.BadRequest(res, "Invalid JSON body.")
    return
  }

  if utils.CheckifStringIsEmpty(body.Name) {
    api_hnd.BadRequest(res, "Enum type name is required.")
    return
  }

  ctx := context.Background()
  if err := api_hnd.Dom.CreateEnumType(ctx, body.Name, body.Description); err != nil {
    utils.Warning(fmt.Sprintf("Error creating enum type: %v", err))
    api_hnd.InternalServerError(res, "Failed to create enum type.")
    return
  }

  api_hnd.DynamicResponse(res, map[string]interface{}{
    "status":  "success",
    "message": "Enum type created.",
  })
  return
}

// ====================================
// LIST ENUM TYPES
// ====================================
func (api_hnd *APIHandler) ListEnumTypes(res http.ResponseWriter, req *http.Request) {
  if req.Method != http.MethodGet {
    api_hnd.MethodNotAllowed(res, "GET")
    return
  }

  limit := utils.StringToInt(req.URL.Query().Get("limit"))
  offset := utils.StringToInt(req.URL.Query().Get("offset"))
  if limit <= 0 {
    limit = 50
  }

  ctx := context.Background()
  enums, err := api_hnd.Dom.ListEnumTypes(ctx, limit, offset)
  if err != nil {
    utils.Warning(fmt.Sprintf("Error listing enum types: %v", err))
    api_hnd.InternalServerError(res, "Failed to retrieve enum types.")
    return
  }

  res.Header().Set("Content-Type", "application/json")
  json.NewEncoder(res).Encode(enums)
  return
}

// ====================================
// CREATE CATEGORY (ADMIN ONLY)
// ====================================
func (api_hnd *APIHandler) CreateCategory(res http.ResponseWriter, req *http.Request) {
  if req.Method != http.MethodPost {
    api_hnd.MethodNotAllowed(res, "POST")
    return
  }

  if !api_hnd.IsAdmin(req) {
    api_hnd.Unauthorized(res, "Only admins can create categories.")
    return
  }

  var body struct {
    Name string `json:"name"`
  }

  if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
    api_hnd.BadRequest(res, "Invalid JSON body.")
    return
  }

  if utils.CheckifStringIsEmpty(body.Name) {
    api_hnd.BadRequest(res, "Category name is required.")
    return
  }

  ctx := context.Background()
  if err := api_hnd.Dom.CreateCategory(ctx, body.Name); err != nil {
    utils.Warning(fmt.Sprintf("Error creating category: %v", err))
    api_hnd.InternalServerError(res, "Failed to create category.")
    return
  }

  api_hnd.DynamicResponse(res, map[string]interface{}{
    "status":  "success",
    "message": "Category created.",
  })
  return
}

// ====================================
// LIST CATEGORIES
// ====================================
func (api_hnd *APIHandler) ListCategories(res http.ResponseWriter, req *http.Request) {
  if req.Method != http.MethodGet {
    api_hnd.MethodNotAllowed(res, "GET")
    return
  }

  limit := utils.StringToInt(req.URL.Query().Get("limit"))
  offset := utils.StringToInt(req.URL.Query().Get("offset"))
  if limit <= 0 {
    limit = 50
  }

  ctx := context.Background()
  cats, err := api_hnd.Dom.ListCategories(ctx, limit, offset)
  if err != nil {
    utils.Warning(fmt.Sprintf("Error listing categories: %v", err))
    api_hnd.InternalServerError(res, "Failed to retrieve categories.")
    return
  }

  res.Header().Set("Content-Type", "application/json")
  json.NewEncoder(res).Encode(cats)
  return
}

// ====================================
// SEARCH BLOGS BY TAG
// ====================================
func (api_hnd *APIHandler) SearchBlogsByTag(res http.ResponseWriter, req *http.Request) {
  if req.Method != http.MethodGet {
    api_hnd.MethodNotAllowed(res, "GET")
    return
  }

  tag := req.URL.Query().Get("tag")
  if utils.CheckifStringIsEmpty(tag) {
    api_hnd.BadRequest(res, "Tag is required.")
    return
  }

  limit := utils.StringToInt(req.URL.Query().Get("limit"))
  offset := utils.StringToInt(req.URL.Query().Get("offset"))
  if limit <= 0 {
    limit = 20
  }

  ctx := context.Background()
  blogs, err := api_hnd.Dom.SearchBlogsByTag(ctx, tag, limit, offset)
  if err != nil {
    utils.Warning(fmt.Sprintf("Error searching blogs by tag: %v", err))
    api_hnd.InternalServerError(res, "Failed to search blogs.")
    return
  }

  res.Header().Set("Content-Type", "application/jsoVBn")
  json.NewEncoder(res).Encode(blogs)
  return
}
