package services


import (
  "fmt"
  "context"
  "net/url"
	//"encoding/json"
  dfn"github.com/alphamystic/odin/lib/definers"
  "github.com/alphamystic/odin/lib/utils"
)

// Reports are blogs which can be pentesting reports, Vulnerability/Pentesting and Audit Reports etc
type (
  Report interface {
    CreateBlog(ctx context.Context, blog dfn.Blog) error // this can replace the importing
    ViewBlog(ctx context.Context, token, uuid string) (dfn.Blog, error)
    UpdateBlog(ctx context.Context, blog dfn.Blog) error
    GetRecentBlogs(ctx context.Context) ([]dfn.Blog, error)
    ArchiveBlog(ctx context.Context, uuid string, archived bool) error
    ListBlogsByType(ctx context.Context, blogType string, limit,offset int) ([]dfn.Blog, error)
    ListBlogsByAuthor(ctx context.Context, author string, limit,offset int) ([]dfn.Blog, error)
    ListBlogsByMainTag(ctx context.Context, token, tag string, limit,offset int) ([]dfn.Blog, error)
    // Category & Enum Management
    CreateCategory(ctx context.Context, cat dfn.Category) error
    ListCategories(ctx context.Context, token string, limit,offset int) ([]dfn.Category, error)
    ListEnumTypes(ctx context.Context, token string, limit,offset int) ([]dfn.EnumType, error)
    // Comments
    CreateComment(ctx context.Context, comm dfn.Comment) error
    ListComments(ctx context.Context, token, blogUUID string, limit int) ([]dfn.Comment, error)
    GenerateReport(ctx context.Context) //Exporting a report
    ImportReport(ctx context.Context)
  }
  ReportService struct {SAC *ServerAPIConnector}
)

func NewReportService(sac *ServerAPIConnector) *ReportService{
  return &ReportService{SAC: sac}
}

func (s *ReportService) CreateBlog(ctx context.Context, blog dfn.Blog) error {
	_, err := s.SAC.Post("/api/report/createblog", blog)
	return err
}

func (s *ReportService) UpdateBlog(ctx context.Context, blog dfn.Blog) error {
	_, err := s.SAC.Post("/api/report/updateblog/", blog)
	return err
}

func (s *ReportService) ViewBlog(ctx context.Context, token, uuid string) (dfn.Blog, error) {
    var blog dfn.Blog
    endpoint := fmt.Sprintf("/api/report/viewblog?uuid=%s", uuid)
    // 2. Make the authorized request
    resp, err := s.SAC.AuthGet(endpoint, token)
    if err != nil {
        return blog, err
    }
    // 3. Extract the "data" field from the response map
    // Your API returns: {"status": "success", "data": [...], "message": "..."}
    dataField, ok := resp["data"]
    if !ok {
        return blog, fmt.Errorf("API response missing 'data' field")
    }

    // 4. Decode ONLY the data slice into your struct
    if err := s.SAC.Decode(dataField, &blog); err != nil {
        return blog, fmt.Errorf("failed to decode blog: %w", err)
    }
    utils.Warning(fmt.Sprintf("%v",blog))
    return blog, nil
}

func (s *ReportService) ListComments(ctx context.Context, token, blogUUID string, limit int) ([]dfn.Comment, error) {
	var comments []dfn.Comment
    endpoint := fmt.Sprintf("/api/report/listcomments/?uuid=%s&limit=20%d", blogUUID, limit)
    // 2. Make the authorized request
    resp, err := s.SAC.AuthGet(endpoint, token)
    if err != nil {
        return comments, err
    }
    // 3. Extract the "data" field from the response map
    // Your API returns: {"status": "success", "data": [...], "message": "..."}
    dataField, ok := resp["data"]
    if !ok {
        return comments, fmt.Errorf("API response missing 'data' field")
    }

    // 4. Decode ONLY the data slice into your struct
    if err := s.SAC.Decode(dataField, &comments); err != nil {
        return comments, fmt.Errorf("failed to decode blog: %w", err)
    }
    utils.Warning(fmt.Sprintf("%v",comments))
    return comments, nil
}

func (s *ReportService) ListBlogsByType(ctx context.Context, blogType string, limit,offset int) ([]dfn.Blog, error) {
	data, err := s.SAC.Get(fmt.Sprintf("/api/report/listblogsbytype/?type=%s&limit=%d&offset=%d", blogType, limit,offset))
	var blogs []dfn.Blog
	err = s.SAC.Decode(data, &blogs)
	return blogs, err
}

func (s *ReportService) ListBlogsByAuthor(ctx context.Context, author string, limit,offset int) ([]dfn.Blog, error) {
	data, err := s.SAC.Get(fmt.Sprintf("/api/report/llistblogsbyauthor?author=&limit=%d&offset=%d", author, limit, offset))
	var blogs []dfn.Blog
	err = s.SAC.Decode(data, &blogs)
	return blogs, err
}

func (s *ReportService) ListBlogsByMainTag(ctx context.Context, token, tag string, limit, offset int) ([]dfn.Blog, error) {
    // 1. Properly escape the tag (e.g., "Offensive Security" -> "Offensive%20Security")
    encodedTag := url.QueryEscape(tag)
    endpoint := fmt.Sprintf("/api/report/listblogsbymaintag?tag=%s&limit=%d&offset=%d", encodedTag, limit, offset)
    // 2. Make the authorized request
    resp, err := s.SAC.AuthGet(endpoint, token)
    if err != nil {
        return nil, err
    }
    // 3. Extract the "data" field from the response map
    // Your API returns: {"status": "success", "data": [...], "message": "..."}
    dataField, ok := resp["data"]
    if !ok {
        return nil, fmt.Errorf("API response missing 'data' field")
    }
    var blogs []dfn.Blog
    // 4. Decode ONLY the data slice into your struct
    if err := s.SAC.Decode(dataField, &blogs); err != nil {
        return nil, fmt.Errorf("failed to decode blogs: %w", err)
    }

    return blogs, nil
}

func (s *ReportService) ListCategories(ctx context.Context, token string, limit, offset int) ([]dfn.Category, error) {
    // 1. Properly escape the tag (e.g., "Offensive Security" -> "Offensive%20Security")
    endpoint := fmt.Sprintf("/api/report/category/list?limit=%d&offset=%d", limit, offset)
    // 2. Make the authorized request
    resp, err := s.SAC.AuthGet(endpoint, token)
    if err != nil {
        return nil, err
    }
    utils.Warning(fmt.Sprintf("%v",resp))
    // 3. Extract the "data" field from the response map
    // Your API returns: {"status": "success", "data": [...], "message": "..."}
    dataField, ok := resp["data"]
    if !ok {
        return nil, fmt.Errorf("API response missing 'data' field")
    }
    utils.Warning(fmt.Sprintf("%v",dataField))
    var categories []dfn.Category
    // 4. Decode ONLY the data slice into your struct
    if err := s.SAC.Decode(dataField, &categories); err != nil {
        return nil, fmt.Errorf("failed to decode categories: %w", err)
    }
    utils.Warning(fmt.Sprintf("%v",dataField))
    return categories, nil
}

func (s *ReportService) ListEnumTypes(ctx context.Context, token string, limit, offset int) ([]dfn.EnumType, error) {
    endpoint := fmt.Sprintf("/api/report/enum/list?limit=%d&offset=%d", limit, offset)
    // 2. Make the authorized request
    resp, err := s.SAC.AuthGet(endpoint, token)
    if err != nil {
        return nil, err
    }
    utils.Warning(fmt.Sprintf("%v",resp))
    // 3. Extract the "data" field from the response map
    // Your API returns: {"status": "success", "data": [...], "message": "..."}
    dataField, ok := resp["data"]
    if !ok {
        return nil, fmt.Errorf("API response missing 'data' field")
    }
    utils.Warning(fmt.Sprintf("%v",resp))
    var parent_categories []dfn.EnumType
    // 4. Decode ONLY the data slice into your struct
    if err := s.SAC.Decode(dataField, &parent_categories); err != nil {
        return nil, fmt.Errorf("failed to decode enum types: %w", err)
    }
    utils.Warning(fmt.Sprintf("%v",parent_categories))
    return parent_categories, nil
}







