package handlers

import (
	"fmt"
	"strings"
	"net/http"
	"encoding/json"
	"html/template"
	"path/filepath"
	"github.com/alphamystic/odin/lib/utils"
	dfn "github.com/alphamystic/odin/lib/definers"

	//"github.com/luandm/go-docx" // Example for docx to html
    "code.sajari.com/docconv"   // Multi-format converter
)

func (hnd *Handler) ParseDocument(res http.ResponseWriter, req *http.Request) {
    if req.Method != http.MethodPost {
        http.Error(res, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    // 1. Receive file
    req.ParseMultipartForm(10 << 20)
    file, header, err := req.FormFile("file")
    if err != nil {
        utils.Warning(fmt.Sprintf("Upload Error: %v", err))
        http.Error(res, "File not found", http.StatusBadRequest)
        return
    }
    defer file.Close()

    // 2. Translate to HTML
    mimeType := header.Header.Get("Content-Type")

    // docconv.Convert returns (*docconv.Response, error)
    resBody, err := docconv.Convert(file, mimeType, true)
    if err != nil {
        utils.Warning(fmt.Sprintf("Conversion Error: %v", err))
        http.Error(res, "Conversion failed", http.StatusInternalServerError)
        return
    }

    // 3. Respond with HTML string for the editor
    // Fix: resBody is a struct, use resBody.Body to get the actual string
    res.Header().Set("Content-Type", "application/json")
    json.NewEncoder(res).Encode(map[string]string{
        "html":  resBody.Body, // Accessing the Body string field
        "title": strings.TrimSuffix(header.Filename, filepath.Ext(header.Filename)),
    })
    return
}

// ManageBlog handles the creation UI (GET) and the AJAX save (POST)
func (hnd *Handler) ManageBlog(res http.ResponseWriter, req *http.Request) {
	ud, authenticated := hnd.AuthenticateUser(res, req)
	if !authenticated {
		return
	}
    token, _ := hnd.GetToken(req)

	if req.Method == http.MethodGet {
	    // Get the Enum Types
	    ctx := req.Context()
	    parent_categories, err := hnd.SRVCS.ReportSrvs.ListEnumTypes(ctx, token, 100, 0)
	    if err != nil {
	        utils.Warning(fmt.Sprintf("Error Getting Parent Categories: %v", err))
	        parent_categories = []dfn.EnumType{}
	    }
        categories, err := hnd.SRVCS.ReportSrvs.ListCategories(ctx, token, 100, 0)
        if err != nil {
            utils.Warning(fmt.Sprintf("Error Getting  Categories: %v", err))
            categories = []dfn.Category{}
        }
	    // Get the categories
		tpl, err := hnd.Pages.GetATemplate("create_blog", "create-blog.tmpl")
		if err != nil {
		    utils.Warning(fmt.Sprintf("Template Load Error: %v", err))
			hnd.RenderErrorPage(res, req, ErrorPage{
				ErrorCode: 500,
				Message:   "Failed to load the editor template.",
				Back:      "Go to Dashboard",
				Direction: "/dashboard",
			})
			return
		}
		err = tpl.ExecuteTemplate(res, "create_blog", map[string]interface{}{
			"UserData": ud,
			"parent_categories":  parent_categories,
			"categories": categories,
		})
        if err != nil {
            utils.Warning(fmt.Sprintf("Template(create_blog) Execution Error: %v", err))
        }
		return
	}

	if req.Method == http.MethodPost {
		var blog dfn.Blog
		if err := json.NewDecoder(req.Body).Decode(&blog); err != nil {
		utils.Warning(fmt.Sprintf("Json Decode Error: %v", err))
		hnd.RenderErrorPage(res, req, ErrorPage{
                ErrorCode: 400,
                Message:   "Invalid Request Format. :)",
                Back:      "Go to Dashboard",
                Direction: "/dashboard",
            })
			return
		}

		blog.Author = ud.UserName
		blog.OwnerID = ud.UserId

		ctx := req.Context()
		if err := hnd.SRVCS.ReportSrvs.CreateBlog(ctx, blog); err != nil {
            utils.Warning(fmt.Sprintf("%v",err))
            hnd.RenderErrorPage(res, req, ErrorPage{
                            ErrorCode: 500,
                            Message:   "Failed to save blog to API server.",
                            Back:      "Go to Dashboard",
                            Direction: "/dashboard",
                        })
                return
        }
		http.Redirect(res, req, "/reports/list", http.StatusSeeOther)
		return
	}
}

// ViewReport displays a specific report or renders an error if not found
func (hnd *Handler) ViewReport(res http.ResponseWriter, req *http.Request) {
    ud, authenticated := hnd.AuthenticateUser(res, req)
    if !authenticated {
        return
    }

    uuid := req.URL.Query().Get("uuid")
    limit := utils.StringToInt(req.URL.Query().Get("limit"))
    if limit <= 0 {
        limit = 20
    }

    if uuid == "" {
        hnd.RenderErrorPage(res, req, ErrorPage{
            ErrorCode: 400,
            Message:   "Missing Report UUID. Please select a valid report.",
            Back:      "Back to Reports",
            Direction: "/listreports",
        })
        return
    }

    ctx := req.Context()
    token, _ := hnd.GetToken(req)

    // 1. Fetch Blog
    blog, err := hnd.SRVCS.ReportSrvs.ViewBlog(ctx, token, uuid)
    if err != nil {
        utils.Warning(fmt.Sprintf("Error from service: %v", err))
        hnd.RenderErrorPage(res, req, ErrorPage{
            ErrorCode: 404,
            Message:   fmt.Sprintf("Could not find report: %s", uuid),
            Data:      err.Error(),
            Back:      "Back to Reports",
            Direction: "/listreports",
        })
        return
    }

    // Wrap for HTML safety
    type SafeBlog struct {
        dfn.Blog
        SafeContent template.HTML
    }
    displayBlog := SafeBlog{
        Blog:        blog,
        SafeContent: template.HTML(blog.Content),
    }

    // 2. Fetch Comments (Don't return on error!)
    comments, err := hnd.SRVCS.ReportSrvs.ListComments(ctx, token, uuid, limit)
    if err != nil {
        utils.Warning(fmt.Sprintf("Error getting comments: %v", err))
        // Ensure comments is an empty slice so the template range doesn't break
        comments = []dfn.Comment{}
    }

    // 3. Render Template
    tpl, err := hnd.Pages.GetATemplate("report_view", "report-view.tmpl")
    if err != nil {
        utils.Warning(fmt.Sprintf("Template Load Error: %v", err))
        hnd.Internalserverror(res, req)
        return
    }

    // Ensure we pass the data even if comments are empty
    err = tpl.ExecuteTemplate(res, "report_view", map[string]interface{}{
        "Blog":     displayBlog,
        "Comments": comments,
        "UserData": ud,
    })
    if err != nil {
        utils.Warning(fmt.Sprintf("Template Execution Error: %v", err))
    }
}

// BugBountyReports populates the bug bounty listing
func (hnd *Handler) BugBountyReports(res http.ResponseWriter, req *http.Request) {
	ud, authenticated := hnd.AuthenticateUser(res, req)
	if !authenticated {
		return
	}
    token, _ := hnd.GetToken(req)
    limit := utils.StringToInt(req.URL.Query().Get("limit"))
    offset := utils.StringToInt(req.URL.Query().Get("offset"))
    if limit <= 0 {
      limit = 20
    }

	reports, err := hnd.SRVCS.ReportSrvs.ListBlogsByMainTag(req.Context(), token, "Bug Bounty Report",limit,offset)
	if err != nil {
	    utils.Warning(fmt.Sprintf("Error from service: %v", err))
		hnd.RenderErrorPage(res, req, ErrorPage{
			ErrorCode: 502,
			Message:   "Unable to fetch Bug Bounty reports from the API.",
			Data:      err.Error(),
			Back:      "Reload Page",
			Direction: "/reports/bugbounty",
		})
		return
	}

	tpl, err := hnd.Pages.GetATemplate("reports", "sample_blogs_reports.tmpl")
	if err != nil {
	    utils.Warning(fmt.Sprintf("Template Load Error: %v", err))
		hnd.Internalserverror(res, req)
		return
	}

	tpl.ExecuteTemplate(res, "reports", map[string]interface{}{
		"Reports":  reports,
		"UserData": ud,
		"Title":    "Bug Bounty Reports",
	})
    return
}

// PentestsReports populates the pentesting listing
func (hnd *Handler) PentestsReports(res http.ResponseWriter, req *http.Request) {
	ud, authenticated := hnd.AuthenticateUser(res, req)
	if !authenticated {
		return
	}
    token, _ := hnd.GetToken(req)

    limit := utils.StringToInt(req.URL.Query().Get("limit"))
    offset := utils.StringToInt(req.URL.Query().Get("offset"))
    if limit <= 0 {
      limit = 20
    }

	reports, err := hnd.SRVCS.ReportSrvs.ListBlogsByMainTag(req.Context(), token, "Pentesting Report", limit, offset)
	if err != nil {
	    utils.Warning(fmt.Sprintf("Error from service: %v", err))
		hnd.RenderErrorPage(res, req, ErrorPage{
			ErrorCode: 502,
			Message:   "Gateway Error: Failed to retrieve Pentest reports.",
			Data:      err.Error(),
			Back:      "Back to Dashboard",
			Direction: "/dashboard",
		})
		return
	}

	tpl, err := hnd.Pages.GetATemplate("pentest_reports", "pentest_reports.tmpl")
	if err != nil {
	    utils.Warning(fmt.Sprintf("Template Load Error: %v", err))
		hnd.Internalserverror(res, req)
		return
	}

	tpl.ExecuteTemplate(res, "pentest_reports", map[string]interface{}{
		"Reports":  reports,
		"UserData": ud,
	})
    return
}

// Zerodays populates specialized zero-day research reports
func (hnd *Handler) Zerodays(res http.ResponseWriter, req *http.Request) {
	ud, authenticated := hnd.AuthenticateUser(res, req)
	if !authenticated {
		return
	}
    token, _ := hnd.GetToken(req)

    limit := utils.StringToInt(req.URL.Query().Get("limit"))
    offset := utils.StringToInt(req.URL.Query().Get("offset"))
    if limit <= 0 {
      limit = 20
    }

	reports, err := hnd.SRVCS.ReportSrvs.ListBlogsByMainTag(req.Context(), token, "Zero-Day", limit, offset)
	if err != nil {
	    utils.Warning(fmt.Sprintf("Error from service: %v", err))
		hnd.RenderErrorPage(res, req, ErrorPage{
			ErrorCode: 500,
			Message:   "Critical error accessing Zero-Day database.",
			Back:      "Logout",
			Direction: "/logout",
		})
		return
	}

	tpl, err := hnd.Pages.GetATemplate("blank", "blank.tmpl")
	if err != nil {
	    utils.Warning(fmt.Sprintf("Template Load Error: %v", err))
		hnd.Internalserverror(res, req)
		return
	}

	tpl.ExecuteTemplate(res, "blank", map[string]interface{}{
		"Reports":  reports,
		"UserData": ud,
		"Title":    "Zero Day Research",
	})
    return
}

// ReportViewModel helps pass both the struct and its JSON representation to the UI
type ReportViewModel struct {
    Data     dfn.Blog // Use concrete type instead of interface{}
    FullJSON string
}

func (hnd *Handler) RenderBlogList(res http.ResponseWriter, req *http.Request) {
	ud, authenticated := hnd.AuthenticateUser(res, req)
	if !authenticated {
		return
	}
    token, _ := hnd.GetToken(req)

	limit := utils.StringToInt(req.URL.Query().Get("limit"))
	offset := utils.StringToInt(req.URL.Query().Get("offset"))
	if limit <= 0 {
		limit = 20
	}

	// Fetch reports from Service B
	rawReports, err := hnd.SRVCS.ReportSrvs.ListBlogsByMainTag(req.Context(), token, "Offensive Security", limit, offset)
	if err != nil {
		utils.Warning(fmt.Sprintf("Service Error: %v", err))
		hnd.RenderErrorPage(res, req, ErrorPage{
			ErrorCode: 502,
			Message:   "Gateway Error: Failed to retrieve Pentest reports.",
			Data:      err.Error(),
			Back:      "Back to Dashboard",
			Direction: "/dashboard",
		})
		return
	}

	// Pre-process reports into ViewModels to avoid Template FuncMap issues
	var decoratedReports []ReportViewModel
    for _, r := range rawReports {
        jsonData, err := json.Marshal(r)
        if err != nil {
            utils.Warning(fmt.Sprintf("Marshal Error: %v", err))
            continue
        }
        decoratedReports = append(decoratedReports, ReportViewModel{
            Data:     r, // Now correctly typed
            FullJSON: string(jsonData),
        })
    }

	// Load template
	tpl, err := hnd.Pages.GetATemplate("pentest_reports", "bloglist.tmpl")
	if err != nil {
		utils.Warning(fmt.Sprintf("Template Load Error: %v", err))
		hnd.Internalserverror(res, req)
		return
	}

	// Execute with decorated reports
	utils.Warning(fmt.Sprintf("%s",decoratedReports))
	if err := tpl.ExecuteTemplate(res, "pentest_reports", map[string]interface{}{
		"Reports":  decoratedReports,
		"UserData": ud,
	}); err != nil {
		utils.Warning(fmt.Sprintf("Template Execution Error: %v", err))
	}
    return
}


// Filters to use, use the MainTag as Enum Types
// Endpoint Security
// Exploitation Proofs
// Compliance Audits
// Internal Controls
// SOP's Standard Operating Procedures
// Offensive Security  //using this as pentests
// Security Advisory
// Bug Bounty