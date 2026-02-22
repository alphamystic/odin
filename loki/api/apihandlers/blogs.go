package apihandlers

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"

    dfn "github.com/alphamystic/odin/lib/definers"
    "github.com/alphamystic/odin/lib/utils"
)

// ====================================
// CREATE BLOG
// ====================================
func (api_hnd *APIHandler) CreateBlog(res http.ResponseWriter, req *http.Request) {
  if req.Method != http.MethodPost {
    api_hnd.MethodNotAllowed(res, "POST")
    return
  }

  var blog dfn.Blog
  if err := json.NewDecoder(req.Body).Decode(&blog); err != nil {
    api_hnd.BadRequest(res, "Invalid JSON body.")
    return
  }

  if utils.CheckifStringIsEmpty(blog.OwnerID) {
    api_hnd.BadRequest(res, "OwnerID is required.")
    return
  }
  if utils.CheckifStringIsEmpty(blog.Title) {
    api_hnd.BadRequest(res, "Title is required.")
    return
  }

  // Assign UUID + timestamps
  blog.UUID = utils.GenerateUUID()
  blog.Touch()

  ctx := context.Background()
  if err := api_hnd.Dom.CreateBlog(ctx, blog); err != nil {
    utils.Warning(fmt.Sprintf("Error creating blog: %v", err))
    api_hnd.InternalServerError(res, "Failed to create blog entry.")
    return
  }

  api_hnd.DynamicResponse(res, map[string]interface{}{
    "status":  "success",
    "message": "Blog created successfully.",
    "data":    blog.UUID,
  })
}

// ====================================
// VIEW BLOG
// ====================================
func (api_hnd *APIHandler) ViewBlog(res http.ResponseWriter, req *http.Request) {
  if req.Method != http.MethodGet {
    api_hnd.MethodNotAllowed(res, "GET")
    return
  }

  blogUUID := req.URL.Query().Get("uuid")
  if utils.CheckifStringIsEmpty(blogUUID) {
    api_hnd.BadRequest(res, "UUID is required.")
    return
  }

  ctx := context.Background()
  blog, err := api_hnd.Dom.ViewBlog(ctx, blogUUID)
  if err != nil {
    utils.Warning(fmt.Sprintf("Error viewing blog: %v", err))
    api_hnd.InternalServerError(res, "Failed to fetch blog.")
    return
  }

  api_hnd.DynamicResponse(res, map[string]interface{}{
    "status":  "success",
    "message": "Blog retrieved.",
    "data":    blog,
  })
}

// ====================================
// GET RECENT BLOGS
// ====================================
func (api_hnd *APIHandler) GetRecentBlogs(res http.ResponseWriter, req *http.Request) {
  if req.Method != http.MethodGet {
    api_hnd.MethodNotAllowed(res, "GET")
    return
  }
  ud := req.Context().Value("userData").(*UserData)
  if ud == nil {
    api_hnd.Unauthorized(res, "User data not found")
    return
  }

  ctx := context.Background()
  blogs, err := api_hnd.Dom.GetRecentBlogs(ctx, ud.OwnerID)
  if err != nil {
    utils.Warning(fmt.Sprintf("Error retrieving recent blogs: %v", err))
    api_hnd.InternalServerError(res, "Failed to fetch blogs.")
    return
  }

  res.Header().Set("Content-Type", "application/json")
  json.NewEncoder(res).Encode(blogs)
  return
}

// ====================================
// LIST BLOGS BY AUTHOR
// ====================================
func (api_hnd *APIHandler) ListBlogsByAuthor(res http.ResponseWriter, req *http.Request) {
  if req.Method != http.MethodGet {
    api_hnd.MethodNotAllowed(res, "GET")
    return
  }

  author := req.URL.Query().Get("author")
  if utils.CheckifStringIsEmpty(author) {
    api_hnd.BadRequest(res, "Author is required.")
    return
  }

  limit := utils.StringToInt(req.URL.Query().Get("limit"))
  offset := utils.StringToInt(req.URL.Query().Get("offset"))
  if limit <= 0 {
    limit = 20
  }

  ctx := context.Background()
  blogs, err := api_hnd.Dom.ListBlogByAuthor(ctx, author, limit, offset)
  if err != nil {
    utils.Warning(fmt.Sprintf("Error listing blogs by author: %v", err))
    api_hnd.InternalServerError(res, "Failed to list blogs.")
    return
  }

  res.Header().Set("Content-Type", "application/json")
  json.NewEncoder(res).Encode(blogs)
}

// ====================================
// LIST BLOGS BY MAIN TAG
// ====================================
func (api_hnd *APIHandler) ListBlogsByMainTag(res http.ResponseWriter, req *http.Request) {
  if req.Method != http.MethodGet {
    api_hnd.MethodNotAllowed(res, "GET")
    return
  }

  tag := req.URL.Query().Get("tag")
  if utils.CheckifStringIsEmpty(tag) {
    api_hnd.BadRequest(res, "MainTag is required.")
    return
  }

  limit := utils.StringToInt(req.URL.Query().Get("limit"))
  offset := utils.StringToInt(req.URL.Query().Get("offset"))
  if limit <= 0 {
    limit = 20
  }

  ctx := context.Background()
  blogs, err := api_hnd.Dom.ListBlogByMainTag(ctx, tag, limit, offset)
  if err != nil {
    utils.Warning(fmt.Sprintf("Error listing blogs by maintag: %v", err))
    api_hnd.InternalServerError(res, "Failed to fetch blogs.")
    return
  }

  res.Header().Set("Content-Type", "application/json")
  json.NewEncoder(res).Encode(blogs)
}

// ====================================
// LIST BLOGS BY TYPE
// ====================================
func (api_hnd *APIHandler) ListBlogsByType(res http.ResponseWriter, req *http.Request) {
  if req.Method != http.MethodGet {
    api_hnd.MethodNotAllowed(res, "GET")
    return
  }

  blogType := req.URL.Query().Get("type")
  if utils.CheckifStringIsEmpty(blogType) {
    api_hnd.BadRequest(res, "Type is required.")
    return
  }

  limit := utils.StringToInt(req.URL.Query().Get("limit"))
  offset := utils.StringToInt(req.URL.Query().Get("offset"))
  if limit <= 0 {
    limit = 20
  }

  ctx := context.Background()
  blogs, err := api_hnd.Dom.ListByType(ctx, blogType, limit, offset)
  if err != nil {
    utils.Warning(fmt.Sprintf("Error listing blogs by type: %v", err))
    api_hnd.InternalServerError(res, "Failed to fetch blogs.")
    return
  }

  res.Header().Set("Content-Type", "application/json")
  json.NewEncoder(res).Encode(blogs)
}

// ====================================
// LIST BLOGS BY CATEGORY
// ====================================
func (api_hnd *APIHandler) ListBlogsByCategory(res http.ResponseWriter, req *http.Request) {
  if req.Method != http.MethodGet {
    api_hnd.MethodNotAllowed(res, "GET")
    return
  }

  catStr := req.URL.Query().Get("category")
  category := utils.StringToInt(catStr)
  if category <= 0 {
    api_hnd.BadRequest(res, "Valid category ID is required.")
    return
  }

  limit := utils.StringToInt(req.URL.Query().Get("limit"))
  offset := utils.StringToInt(req.URL.Query().Get("offset"))
  if limit <= 0 {
    limit = 20
  }

  ctx := context.Background()
  blogs, err := api_hnd.Dom.ListByCategory(ctx, category, limit, offset)
  if err != nil {
    utils.Warning(fmt.Sprintf("Error listing blogs by category: %v", err))
    api_hnd.InternalServerError(res, "Failed to fetch blogs.")
    return
  }

  res.Header().Set("Content-Type", "application/json")
  json.NewEncoder(res).Encode(blogs)
}

// ====================================
// UPDATE BLOG (ADMIN ONLY)
// ====================================
func (api_hnd *APIHandler) UpdateBlog(res http.ResponseWriter, req *http.Request) {
  if req.Method != http.MethodPost {
    api_hnd.MethodNotAllowed(res, "POST")
    return
  }

  if !api_hnd.IsAdmin(req) {
    api_hnd.Unauthorized(res, "Only admins can update blogs.")
    return
  }

  var update dfn.Blog
  if err := json.NewDecoder(req.Body).Decode(&update); err != nil {
    api_hnd.BadRequest(res, "Invalid JSON body.")
    return
  }

  if utils.CheckifStringIsEmpty(update.UUID) {
    api_hnd.BadRequest(res, "UUID is required.")
    return
  }

  ctx := context.Background()
  if err := api_hnd.Dom.UpdateBlog(ctx, update); err != nil {
    utils.Warning(fmt.Sprintf("Error updating blog: %v", err))
    api_hnd.InternalServerError(res, "Failed to update blog.")
    return
  }

  api_hnd.DynamicResponse(res, map[string]interface{}{
    "status":  "success",
    "message": "Blog updated successfully.",
  })
}

// ====================================
// ARCHIVE / UNARCHIVE BLOG (ADMIN ONLY)
// ====================================
func (api_hnd *APIHandler) ArchiveBlog(res http.ResponseWriter, req *http.Request) {
  if req.Method != http.MethodPost {
    api_hnd.MethodNotAllowed(res, "POST")
    return
  }

  if !api_hnd.IsAdmin(req) {
    api_hnd.Unauthorized(res, "Only admins can archive or unarchive blogs.")
    return
  }

  uuid := req.URL.Query().Get("uuid")
  status := req.URL.Query().Get("archived")

  if utils.CheckifStringIsEmpty(uuid) {
    api_hnd.BadRequest(res, "UUID is required.")
    return
  }

  archive := status == "true"

  ctx := context.Background()
  if err := api_hnd.Dom.ArchiveOrRemoveFromArchive(ctx, uuid, archive); err != nil {
    api_hnd.InternalServerError(res, "Failed to update archive status.")
    return
  }

  api_hnd.DynamicResponse(res, map[string]interface{}{
    "status":  "success",
    "message": "Archive status updated.",
  })
}

// ====================================
// CREATE COMMENT
// ====================================
func (api_hnd *APIHandler) CreateComment(res http.ResponseWriter, req *http.Request) {
  if req.Method != http.MethodPost {
    api_hnd.MethodNotAllowed(res, "POST")
    return
  }

  var comment dfn.Comment
  if err := json.NewDecoder(req.Body).Decode(&comment); err != nil {
    api_hnd.BadRequest(res, "Invalid JSON body.")
    return
  }

  if utils.CheckifStringIsEmpty(comment.BlogUUID) {
    api_hnd.BadRequest(res, "BlogUUID is required.")
    return
  }
  if utils.CheckifStringIsEmpty(comment.Commentor) {
    api_hnd.BadRequest(res, "Commentor is required.")
    return
  }

  comment.CommentUUID = utils.GenerateUUID()
  comment.Touch()

  ctx := context.Background()
  if err := api_hnd.Dom.CreateComment(ctx, comment); err != nil {
    api_hnd.InternalServerError(res, "Failed to create comment.")
    return
  }

  api_hnd.DynamicResponse(res, map[string]interface{}{
    "status":  "success",
    "message": "Comment added.",
  })
}

// ====================================
// LIST COMMENTS FOR BLOG
// ====================================
func (api_hnd *APIHandler) ListComments(res http.ResponseWriter, req *http.Request) {
  if req.Method != http.MethodGet {
    api_hnd.MethodNotAllowed(res, "GET")
    return
  }

  blogUUID := req.URL.Query().Get("uuid")
  if utils.CheckifStringIsEmpty(blogUUID) {
    api_hnd.BadRequest(res, "UUID is required.")
    return
  }

  limit := utils.StringToInt(req.URL.Query().Get("limit"))
  offset := utils.StringToInt(req.URL.Query().Get("offset"))
  if limit <= 0 {
    limit = 20
  }

  ctx := context.Background()
  comments, err := api_hnd.Dom.ListComments(ctx, blogUUID, limit, offset)
  if err != nil {
    api_hnd.InternalServerError(res, "Failed to fetch comments.")
    return
  }

  res.Header().Set("Content-Type", "application/json")
  json.NewEncoder(res).Encode(comments)
}
