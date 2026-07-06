package handlers

import (
	"fmt"
	"bytes"
	"net/http"

	"github.com/alphamystic/odin/lib/utils"
    png_hnd"github.com/alphamystic/odin/lib/handlers"
)

// RenderScans displays active scans based on categories
func (hnd *Handler) RenderScans(res http.ResponseWriter, req *http.Request) {
	ud, authenticated := hnd.AuthenticateUser(res, req)
	if !authenticated {
		return
	}
	token, _ := hnd.GetToken(req)
	scanType := req.URL.Query().Get("scantype")

	// FIXED: Removed duplicate 'scan_type :=' short declaration error
	if scanType != "Bug Bounty" && scanType != "Pentest" && scanType != "Black Ops" && scanType != "ALL" {
		scanType = "ALL"
	}

	ctx := req.Context()
	// FIXED: Routed call cleanly to the modern ReconService namespace passing standard token configurations
	scans, err := hnd.SRVCS.ReconSrvs.ListScans(ctx, token, scanType)
	if err != nil {
		utils.Warning(fmt.Sprintf("Error Getting Scans: %v", err))
		hnd.RenderErrorPage(res, req, ErrorPage{
			ErrorCode: 502,
			Message:   "Unable to retrieve scans matrix records from backend pipeline.",
			Back:      "Back to Dashboard",
			Direction: "/",
		})
		return
	}

	tpl, err := hnd.Pages.GetATemplate("scans_view", "scans-view.tmpl")
	if err != nil {
		utils.Warning(fmt.Sprintf("Template Load Error: %s", err))
		hnd.Internalserverror(res, req)
		return
	}

	err = tpl.ExecuteTemplate(res, "scans_view", map[string]interface{}{
		"UserData":  ud,
		"Scans":     scans,
		"ScanType":  scanType,
	})
	if err != nil {
		utils.Warning(fmt.Sprintf("Template Execution Error: %v", err))
	}
    return
}

// ListTargets displays all targeted hosts linked to a distinct scan context
func (hnd *Handler) ListTargets(res http.ResponseWriter, req *http.Request) {
	ud, authenticated := hnd.AuthenticateUser(res, req)
	if !authenticated {
		return
	}
	token, _ := hnd.GetToken(req)
	scanID := req.URL.Query().Get("scanid")

	if scanID == "" {
		hnd.RenderErrorPage(res, req, ErrorPage{
			ErrorCode: 400,
			Message:   "Missing Scan Tracking ID context parameter.",
			Back:      "Go to Scans",
			Direction: "/scans",
		})
		return
	}

	ctx := req.Context()
	targets, err := hnd.SRVCS.ReconSrvs.ListTargets(ctx, token, scanID)
	if err != nil {
		utils.Warning(fmt.Sprintf("Error Getting Targets for Scan %s: %v", scanID, err))
		hnd.RenderErrorPage(res, req, ErrorPage{
			ErrorCode: 502,
			Message:   "Failed loading target index profile elements.",
			Back:      "Go to Scans",
			Direction: "/scans",
		})
		return
	}

	tpl, err := hnd.Pages.GetATemplate("targets_list", "targets-list.tmpl")
	if err != nil {
		utils.Warning(fmt.Sprintf("Template Load Error: %s", err))
		hnd.Internalserverror(res, req)
		return
	}

	_ = tpl.ExecuteTemplate(res, "targets_list", map[string]interface{}{
		"UserData": ud,
		"Targets":  targets,
		"ScanID":   scanID,
	})
    return
}
// ViewTarget evaluates a target's completed attack data surfaces (Web Data, Port Services, Vulnerabilities)
func (hnd *Handler) ViewTarget(res http.ResponseWriter, req *http.Request) {
	ud, authenticated := hnd.AuthenticateUser(res, req)
	if !authenticated {
		return
	}
	token, _ := hnd.GetToken(req)
	scanID := req.URL.Query().Get("scan-id")
	targetID := req.URL.Query().Get("target-id")

	if scanID == "" || targetID == "" {
		hnd.RenderErrorPage(res, req, ErrorPage{
			ErrorCode: 400,
			Message:   "Required scanning tracking path markers are missing.",
			Back:      "Go to Scans",
			Direction: "/scans",
		})
		return
	}

	ctx := req.Context()
	// Fetch core target telemetry metrics (ReconData wrapper maps)
	targetData, err := hnd.SRVCS.ReconSrvs.ViewTarget(ctx, token, scanID, targetID)
	if err != nil {
		utils.Warning(fmt.Sprintf("Error Fetching Target Profiles %s: %v", targetID, err))
		hnd.RenderErrorPage(res, req, ErrorPage{
			ErrorCode: 502,
			Message:   "Unable to recover target surface intelligence parameters.",
			Back:      "Go back to Targets",
			Direction: "/scans/view-target?scanid=" + scanID,
		})
		return
	}

	// ========================================================================
	// 1. DATA EMPTINESS & HYDRATION SAFEGUARDS
	// ========================================================================
	// Ensure that if the API keys do not exist or are explicitly null, they
	// are safely initialized into empty structures. This completely stops
	// template evaluation lookups from throwing missing key panics.
	if targetData == nil {
		targetData = make(map[string]interface{})
	}

	if targetData["target"] == nil {
		targetData["target"] = make(map[string]interface{})
	}

	if targetData["web_data"] == nil {
		targetData["web_data"] = map[string]interface{}{
			"directory":  []string{},
			"parameters": []string{},
			"files":      []string{},
		}
	} else {
		// Even if web_data exists, verify that sub-arrays are not null references
		wd, ok := targetData["web_data"].(map[string]interface{})
		if ok {
			if wd["directory"] == nil {
				wd["directory"] = []string{}
			}
			if wd["parameters"] == nil {
				wd["parameters"] = []string{}
			}
			if wd["files"] == nil {
				wd["files"] = []string{}
			}
		}
	}

	if targetData["services"] == nil {
		targetData["services"] = []map[string]interface{}{}
	}

	// Fetch confirmed vulnerabilities tied to this asset to render threat metrics inline
	vulnerabilities, err := hnd.SRVCS.ReconSrvs.ListVulnerabilities(ctx, token, targetID)
	if err != nil {
		utils.Warning(fmt.Sprintf("Non-fatal logging warning: failed loading vulnerabilities list: %v", err))
		vulnerabilities = []png_hnd.Vulnerabilities{}
	}

	tpl, err := hnd.Pages.GetATemplate("target_view", "target-view.tmpl")
	if err != nil {
		utils.Warning(fmt.Sprintf("Template Load Error: %s", err))
		hnd.Internalserverror(res, req)
		return
	}

	// ========================================================================
	// 2. BUFFERED TEMPLATE EXECUTION & RECOVERY PIPELINE
	// ========================================================================
	// By compiling to an in-memory byte buffer first, we guarantee that if
	// a formatting flaw exists, the app fails gracefully into an error screen
	// instead of delivering a partial rendering or a blank white screen.
	var buf bytes.Buffer

	err = tpl.ExecuteTemplate(&buf, "target_view", map[string]interface{}{
		"UserData":        ud,
		"Target":          targetData["target"],
		"WebData":         targetData["web_data"],
		"Services":        targetData["services"],
		"Vulnerabilities": vulnerabilities,
		"ScanID":          scanID,
	})

	if err != nil {
		// Log the explicit compilation or runtime pipeline error for immediate diagnostic verification
		utils.Warning(fmt.Sprintf("Template compilation failed securely in memory buffer: %v", err))

		// Render a rich error screen safely since no bytes have hit the active ResponseWriter socket
		hnd.RenderErrorPage(res, req, ErrorPage{
			ErrorCode: 500,
			Message:   "An internal processing layout constraint broke during UI compilation. Please check your data fields alignment.",
			Data:      err.Error(),
			Back:      "Go back to Targets",
			Direction: "/scans/view-target?scanid=" + scanID,
		})
		return
	}

	// If compilation succeeded without errors, safely flush the HTML output to the network line
	res.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = buf.WriteTo(res)
}