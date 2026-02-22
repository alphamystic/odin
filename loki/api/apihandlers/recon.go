package apihandlers

import (
  "fmt"
  "errors"
  "context"
  "net/http"
  "encoding/json"
  //"database/sql"

  "github.com/alphamystic/odin/lib/utils"
  //dom"github.com/alphamystic/odin/lib/domain"
  dfn"github.com/alphamystic/odin/lib/definers"
  png_hnd"github.com/alphamystic/odin/lib/handlers"
)

// create a scan
type ScanRequest struct {
    Name     string `json:"name"`
    ScanType string `json:"scantype"`
}

func (api_hnd *APIHandler) Createscan(res http.ResponseWriter, req *http.Request){
  if req.Method != http.MethodPost {
    api_hnd.MethodNotAllowed(res, "POST")
    return
  }
  ud := req.Context().Value("userData").(*UserData)
  if ud == nil {
    api_hnd.Unauthorized(res, "User data not found")
    return
  }
  // Parse JSON body
  var requestBody ScanRequest
  err := json.NewDecoder(req.Body).Decode(&requestBody)
  if err != nil {
    api_hnd.BadRequest(res, "Invalid JSON request body")
    return
  }
  // Validate fields
  if !utils.CheckifStringIsEmpty(requestBody.Name) {
    api_hnd.BadRequest(res, "Name for scan cannot be empty.")
    return
  }
  if !utils.CheckifStringIsEmpty(requestBody.ScanType) {
    api_hnd.BadRequest(res, "Scan Type cannot be empty.")
    return
  }
  if requestBody.ScanType != "Bug Bounty" && requestBody.ScanType != "Pentest" && requestBody.ScanType != "Black Ops" {
    api_hnd.BadRequest(res, "Scan Type can only be 'Bug Bounty', 'Pentest', or 'Black Ops'")
    return
  }
  // Process scan creation
  ctx := context.Background()
  scan_id, err := api_hnd.Dom.CreateScan(ctx, requestBody.Name, requestBody.ScanType, ud.UserId)
  if err != nil {
    utils.Logerror(err)
    api_hnd.InternalServerError(res, "Internal server error, please try again later.")
    return
  }
  api_hnd.Posdataresponse(res, PostDataResponse{
    Status:      "success",
    RedirectUrl: scan_id,
    Message:     fmt.Sprintf("Successfully created scan: %s", scan_id),
  })
  return
}



// list scans
func (api_hnd *APIHandler) Listscans(res http.ResponseWriter, req *http.Request){
  if req.Method != "GET"{
    api_hnd.MethodNotAllowed(res, "GET")
    return
  }
  ud := req.Context().Value("userData").(*UserData)
  if ud == nil {
      api_hnd.Unauthorized(res, "User data not found")
      return
  }
  ctx := context.Background()
  scan_type := req.URL.Query().Get("scantype")
  if scan_type != "Bug Bounty" && scan_type != "Pentest" && scan_type != "Black Ops" && scan_type !=  "ALL"{
    api_hnd.BadRequest(res, "Scan Type can only be 'Bug Bounty', 'Pentest', or 'Black Ops'")
    return
  }
  scans,err := api_hnd.Dom.ListScan(ctx,scan_type,ud.UserId)
  if err != nil {
    api_hnd.InternalServerError(res,"Internal server error, please try again later.")
    utils.Logerror(err)
    return
  }
  api_hnd.DynamicResponse(res,scans)
  return
}


// View a scan is listing targets on it and data on it.
func (api_hnd *APIHandler) Viewscan(res http.ResponseWriter, req *http.Request){
  if req.Method != "GET"{
    api_hnd.MethodNotAllowed(res, "GET")
    return
  }
  ud := req.Context().Value("userData").(*UserData)
  if ud == nil {
      api_hnd.Unauthorized(res, "User data not found")
      return
  }
  ctx := context.Background()
  scan_id := req.URL.Query().Get("scanid")
  if !utils.CheckifStringIsEmpty(scan_id) {
    api_hnd.BadRequest(res,"Scan ID must be a string.")
    return
  }
  targets,err := api_hnd.Dom.ListTargets(ctx,scan_id)
  if err != nil {
    api_hnd.InternalServerError(res,"Internal server error, please try again later.")
    utils.Logerror(err)
    return
  }
  api_hnd.DynamicResponse(res, map[string]interface{}{
    "status":  "success",
    "message": "Targets fetched successfully",
    "data":    targets,
  })
  return
}


type ReceivedTarget struct {
  TargetID string `json: "targetid"`
  ScanID string `json: "scanid"`
  Host string  `json: "host"`//can be null if not specified as a subdomain
  HostIp string`json: "hostip"`
  TargetIp string `json: "targetip"`
  FireWallName string `json: "firewallName"`
  Decoys []string`json: "decoys"`
}

// create a target, something to be enumerated on
func (api_hnd *APIHandler) Createtarget(res http.ResponseWriter, req *http.Request){
  if req.Method != "POST"{
    api_hnd.MethodNotAllowed(res, "POST")
    return
  }
  ud := req.Context().Value("userData").(*UserData)
  if ud == nil {
      api_hnd.Unauthorized(res, "User data not found")
      return
  }
  ctx := context.Background()
  var target ReceivedTarget
  //create a new target and parse the ip then decode it
  err := json.NewDecoder(req.Body).Decode(&target)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
  target.TargetID = utils.Md5Hash(utils.GenerateUUID())
  trgt := &png_hnd.Target{
    TargetID: utils.Md5Hash(utils.GenerateUUID()),
    ScanID: target.ScanID,
    Host: target.Host,
    HostIp : utils.StringToIP(target.HostIp),
    TargetIp: utils.StringToIP(target.TargetIp),
    FireWallName: target.FireWallName,
    Decoys: utils.StringsToIPArray(target.Decoys),
  }
  trgt.Touch()
  err = api_hnd.Dom.WriteTargetToDB(trgt,ctx)
  if err != nil {
    api_hnd.InternalServerError(res,"Error creating target to DB. Try again Later.")
    return
  }
  api_hnd.DynamicResponse(res,PostDataResponse{
    Status: "success",
    RedirectUrl: trgt.TargetID,
  	Message: fmt.Sprintf("Successfully created target: %s",target.TargetID),
  })
  return
}


// cretae recon data, create services and webdata
// create each of this Separetly, i.e services seperate from all this
type WebDataBody struct {
  TargetID string `json: "targetid"`
  Directories []string `jaon: "directories"`
  Parameters []string `json: "parameters"`
  Filepaths []string `json: "filepaths"`
}

func (api_hnd *APIHandler) Createrecondata(res http.ResponseWriter, req *http.Request){
  if req.Method != "POST"{
    api_hnd.MethodNotAllowed(res, "POST")
    return
  }
  ud := req.Context().Value("userData").(*UserData)
  if ud == nil {
      api_hnd.Unauthorized(res, "User data not found")
      return
  }
  ctx := context.Background()
  var wdBody WebDataBody
  err := json.NewDecoder(req.Body).Decode(&wdBody)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
  //validate that each web data value is not empty, yes it can but create a count to 3 or
  // something to ensure not all are empty
  // encode the webdata
  var (
    directory_path, parameter_path, file_path string
  )
  if len(wdBody.Directories) > 0 {
    directory_path,err = utils.ArrayToToken(wdBody.Directories)
    if err != nil {
      utils.NoticeError(fmt.Sprintf("%s",err))
      api_hnd.BadRequest(res,"Directories can not be empty.")
      return
    }
  }
  if len(wdBody.Parameters) > 0 {
    parameter_path,err = utils.ArrayToToken(wdBody.Parameters)
    if err != nil {
      utils.NoticeError(fmt.Sprintf("%s",err))
      api_hnd.BadRequest(res,"Parameters can not be empty.")
      return
    }
  }
  if len(wdBody.Filepaths) > 0 {
    file_path,err = utils.ArrayToToken(wdBody.Filepaths)
    if err != nil {
      utils.NoticeError(fmt.Sprintf("%s",err))
      api_hnd.BadRequest(res,"Filepaths can not be empty.")
      return
    }
  }
  if !utils.CheckifStringIsEmpty(wdBody.TargetID) {
    api_hnd.BadRequest(res,"Target ID can not be empty.")
    return
  }
  if err := api_hnd.Dom.CreateWebdata(ctx, wdBody.TargetID, directory_path, parameter_path, file_path); err != nil {
    utils.NoticeError(fmt.Sprintf("%s",err))
    api_hnd.InternalServerError(res,"Error creating web data. Please try again later :)")
    return
  }
  // return a success
  api_hnd.Success(res,"Successfully created web data.")
  return
}


// list recon data is just viewing a target and the recon data on it.
func (api_hnd *APIHandler) Viewtarget(res http.ResponseWriter, req *http.Request){
  if req.Method != "GET" {
    api_hnd.MethodNotAllowed(res, "GET")
    return
  }
  ud := req.Context().Value("userData").(*UserData)
  if ud == nil {
      api_hnd.Unauthorized(res, "User data not found")
      return
  }
  ctx := context.Background()
  target_id := req.URL.Query().Get("target-id")
  scan_id := req.URL.Query().Get("scan-id")
  // do a check to see if the  user owns the target via the scanid
  if !utils.CheckifStringIsEmpty(target_id) {
    api_hnd.BadRequest(res,"Target ID can not be empty.")
    return
  }
  if !utils.CheckifStringIsEmpty(scan_id) {
    api_hnd.BadRequest(res,"Scan ID can not be empty.")
    return
  }
  // get the target
  trg,err := api_hnd.Dom.ViewTarget(scan_id,target_id,ud.UserId,ctx)
  if err != nil {
    utils.Notice(fmt.Sprintf("%s",err))
    if errors.Is(err,dfn.ScanDoesNotExists){
      api_hnd.BadRequest(res,"Requested target id for scan does not exist.")
      return
    }
    if errors.Is(err,dfn.TargetDoesNotExist ){
      api_hnd.Success(res,"Target does not exist.")
      return
    }
    api_hnd.InternalServerError(res,"Error getting the target.")
    return
  }
  // get the webdata
  wd, count,err := api_hnd.Dom.GetWebData(ctx,target_id)
  if err != nil || count > 2 {
    utils.Notice(fmt.Sprintf("%s",err))
    if errors.Is(err,dfn.WebDataForTargetDoesNotExist) {
      api_hnd.Success(res,"Target doess not have web data available.")
      return
    }
    api_hnd.InternalServerError(res,"Internal Server error on getting webdata, try again later :).")
    return
  }
  services,err := api_hnd.Dom.GetServices(ctx, target_id)
  if err != nil {
    utils.Notice(fmt.Sprintf("%s",err))
    api_hnd.InternalServerError(res,"Internal Server error on getting services, try again later :).")
    return
  }
  // encode the response and send it over
  api_hnd.DynamicResponse(res, map[string]interface{}{
		"status":      "success",
		"redirecturl": "",
		"message":     "Target and it's WebData,serviives fetched successfully",
		"data": map[string]interface{}{
			"target":   trg, // Target details
			"web_data": wd,  // Web data
      "services": services, // services
		},
	})
  return
}

// write data to sb and decode it on read.
func (api_hnd *APIHandler)  Createservice(res http.ResponseWriter, req *http.Request) {
  ud := req.Context().Value("userData").(*UserData)
  if ud == nil {
    api_hnd.Unauthorized(res, "User data not found")
    return
  }
  var srvc,service png_hnd.Service
  err := json.NewDecoder(req.Body).Decode(&srvc)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
  ctx := context.Background()
  if len(srvc.Data) >= 2 {
    data := utils.Base64Decode(srvc.Data)
    if data == "" || len(data) <= 2{
      api_hnd.BadRequest(res,"The Data has to be base 64 encoded")
      return
    }
  }
  service = png_hnd.Service{
    TargetID: srvc.TargetID,
    ServiceName: srvc.ServiceName,
    Port: srvc.Port,
    Protocol: srvc.Protocol,
  	State: srvc.State,
  	Version: srvc.Version,
    Data: srvc.Data,
  }
  if err = api_hnd.Dom.CreateService(ctx, service); err != nil{
    api_hnd.InternalServerError(res,"Error creating service")
    return
  }
  api_hnd.Success(res,"Successfully created servcice.")
  return
}
