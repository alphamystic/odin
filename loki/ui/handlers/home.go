package handlers

import(
  "fmt"
  "log"
  "net/http"
  "encoding/json"
  "github.com/alphamystic/odin/lib/utils"
)



// healthResponse defines the JSON structure for health check responses.
type healthResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// toggleResponse defines the JSON structure for toggle responses.
type toggleResponse struct {
	MaintenanceMode bool   `json:"maintenance_mode"`
	Status          string `json:"status"`
}

// Health handles health check requests and returns a JSON-encoded status.
func (api_hnd *Handler) Health(w http.ResponseWriter, r *http.Request) {
  fmt.Println("Running health request..")
	api_hnd.mu.Lock()
	defer api_hnd.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")

	var resp healthResponse
	if api_hnd.MaintenanceMode {
		w.WriteHeader(http.StatusServiceUnavailable)
		resp = healthResponse{
			Status:  "error",
			Message: "Backend is under maintenance.",
		}
		log.Printf("Health check request received. Responding with 503 Service Unavailable (Maintenance Mode).")
	} else {
		w.WriteHeader(http.StatusOK)
		resp = healthResponse{
			Status:  "ok",
			Message: "Backend is healthy!",
		}
		log.Printf("Health check request received. Responding with 200 OK.")
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Error encoding health response: %v", err)
	}
}

// ToggleMaintenanceModeHandler toggles the maintenance mode and returns the new state as JSON.
func (api_hnd *Handler) ToggleMaintenanceModeHandler(w http.ResponseWriter, r *http.Request) {
	api_hnd.mu.Lock()
	defer api_hnd.mu.Unlock()

	api_hnd.MaintenanceMode = !api_hnd.MaintenanceMode
	status := "disabled"
	if api_hnd.MaintenanceMode {
		status = "enabled"
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	resp := toggleResponse{
		MaintenanceMode: api_hnd.MaintenanceMode,
		Status:          status,
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Error encoding toggle response: %v", err)
	}

	log.Printf("Maintenance mode toggled to %s.", status)
}

func (hnd *Handler) Blank(res http.ResponseWriter, req *http.Request){
  tpl,err := hnd.Pages.GetATemplate("blank","blank.tmpl")
  if err != nil {
    utils.Warning(fmt.Sprintf("%s", err))
    hnd.Internalserverror(res, req)
		return
  }
  tpl.ExecuteTemplate(res,"blank",nil)
  return
}



func (hnd *Handler) Internalserverror(res http.ResponseWriter, req *http.Request) {
  tpl,err := hnd.Pages.GetATemplate("error","error.tmpl")
  if err != nil{
    utils.Warning(fmt.Sprintf("%s",err))
    http.Error(res, "An error occurred", http.StatusInternalServerError)
  }
  tpl.ExecuteTemplate(res,"error",nil)
  return
}

func (hnd *Handler) Home(res http.ResponseWriter, req *http.Request){
  _, authenticated := hnd.AuthenticateUser(res, req)
  if !authenticated {
    return // User is redirected in the helper
  }
  tpl,err := hnd.Pages.GetATemplate("home","home.tmpl")
  if err != nil{
    utils.Warning(fmt.Sprintf("%s",err))
    http.Error(res, "An error occurred", http.StatusInternalServerError)
  }
  tpl.ExecuteTemplate(res,"home",nil)
  return
}
