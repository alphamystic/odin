package router

import (
  "os"
  "log"
  "fmt"
  "time"
  "syscall"
  "context"
  "os/signal"
  "net/http"
  "github.com/alphamystic/odin/lib/utils"
  dfn"github.com/alphamystic/odin/lib/definers"
  "github.com/alphamystic/odin/loki/ui/handlers"
)

// A low Level router exposing the default http

type Router struct {
  Mux *http.ServeMux
  HTTPSvr *http.Server
  HTTPSSvr *http.Server
  Tls bool
}

// should probably receive a server
func NewRouter(httpsSvr,httpSvr *http.Server) *Router {
  return &Router{
    Mux: http.NewServeMux(),
    HTTPSvr: httpSvr,
    HTTPSSvr: httpsSvr,
  }
}

func (rtr *Router) Run(reg bool){
  handlers.Registration = reg
  rtr.HTTPSvr.Handler = rtr.Mux
  rtr.HTTPSSvr.Handler = rtr.Mux
  // start channels to write logs
  ShutdownCh := make(chan bool)
  DoneCh := make(chan bool)
  var err error
  // create a file server for the static files
  fs := http.FileServer(http.Dir("./loki/ui/static"))
  //rtr.Mux.Handle("/static/",http.StripPrefix("/static",fs))
  // Cache static files for 1 hour (adjust as needed)
  rtr.Mux.Handle("/static/", http.StripPrefix("/static", http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
    res.Header().Set("Cache-Control", "max-age=3600")
    fs.ServeHTTP(res,req)
  })))

  // connect to DB (Add your onw connection or load from the environment)
  dbConfig := dfn.IntitializeConnector("root","","localhost","odin")
  dbConn,err := dfn.NewMySQLConnector(dbConfig)
  if err != nil {
    utils.Warning(fmt.Sprintf("Error connecting to the DB. \n[-]   ERROR: %s",err))
    return
  }

  // create a request logger
  rl := utils.NewRequestLogger("./.data/logs/requests/",066)

  // initiate new handler
  hnd,err := handlers.NewHandler(dbConn, ShutdownCh, DoneCh, rl)
  if err != nil {
    utils.Danger(err);return
  }

  // Handlers
  //panel shortcuts (We handle the rest as they come. bit by bit baba)
  rtr.Mux.HandleFunc("/profile",hnd.Profile)
  rtr.Mux.HandleFunc("/updateprofile",hnd.Updateprofile)
  rtr.Mux.HandleFunc("/securityprofile",hnd.Securityprofile)
  rtr.Mux.HandleFunc("/notificationprofile",hnd.Notificationsprofile)

  //panel side
  rtr.Mux.HandleFunc("/test",hnd.Blank)
  rtr.Mux.HandleFunc("/",hnd.Home)
  rtr.Mux.HandleFunc("/mkubwa",hnd.Signin)
  rtr.Mux.HandleFunc("/signout",hnd.Logout)
  rtr.Mux.HandleFunc("/register",hnd.Register)

  rtr.Mux.HandleFunc("/apt",hnd.Apt)
  rtr.Mux.HandleFunc("/edr",hnd.Edr)
  rtr.Mux.HandleFunc("/odin-net",hnd.Odinnet)

  rtr.Mux.HandleFunc("/bb",hnd.Bugbounty)
  rtr.Mux.HandleFunc("/pentests",hnd.Pentests)
  rtr.Mux.HandleFunc("/bo",hnd.Blackops)

  rtr.Mux.HandleFunc("/ms",hnd.Motherships)

  rtr.Mux.HandleFunc("/agents",hnd.Minions)
  rtr.Mux.HandleFunc("/iot",hnd.IotRouters)
  rtr.Mux.HandleFunc("/phone",hnd.AndroidIOs)

  rtr.Mux.HandleFunc("/bbreports",hnd.BugBountyReports)
  rtr.Mux.HandleFunc("/pentestsreport",hnd.PentestsReports)

  rtr.Mux.HandleFunc("/pendingscans",hnd.Pendingscans)
  rtr.Mux.HandleFunc("/phishinglinks",hnd.Phishinglinks)
  rtr.Mux.HandleFunc("/zerodays",hnd.Zerodays)

  rtr.Mux.HandleFunc("/events",hnd.Events)
  rtr.Mux.HandleFunc("/regulars",hnd.RegularUsers)
  rtr.Mux.HandleFunc("/admins",hnd.Admins)
  rtr.Mux.HandleFunc("/activeprojects",hnd.Activeprojects)
  rtr.Mux.HandleFunc("/archivedprojects",hnd.Archivedprojects)

  rtr.Mux.HandleFunc("/bds",hnd.Backdoors)
  rtr.Mux.HandleFunc("/bd-generator",hnd.Backdoorgenerator)

  rtr.Mux.HandleFunc("/listcontacts",hnd.Listcontacts)
  rtr.Mux.HandleFunc("/createcontact",hnd.Createcontact)

  rtr.Mux.HandleFunc("/listapikeys",hnd.Listapikeys)
  rtr.Mux.HandleFunc("/createapikeys",hnd.Createapikeys)

  rtr.Mux.HandleFunc("/issues",hnd.CurrentIssues)
  rtr.Mux.HandleFunc("/appointments",hnd.Viewappointments)

  rtr.Mux.HandleFunc("/docs",hnd.Documentation)
  // End of handlers

  // Start the server on the background
  go func(){
    if err := rtr.HTTPSvr.ListenAndServe(); err != http.ErrServerClosed {
      log.Fatalf("[-] Error starting server: %s\n",err.Error())
    }
  }()
  if rtr.Tls {
    go func(){
      // we need to find a better way of supplying this
      if err := rtr.HTTPSSvr.ListenAndServeTLS("../../../certs/server.crt", "../../../certs/server.key"); err != http.ErrServerClosed {
        log.Fatalf("[-] Error starting HTTPS server: %s\n",err.Error())
      }
    }()
  }
  fmt.Println("Servers are here running")
  interruptChan := make(chan os.Signal,1)
  signal.Notify(interruptChan,os.Interrupt, syscall.SIGTERM)
  //sedn a close channel to the handler
  hnd.ShutdownChan <- true
  // wait for the receiver to finish writing all logs
  <-hnd.DoneChan
  // read from the interrupt chan and shutdown
  <-interruptChan
  shutdownCtx,shutdownCancel := context.WithTimeout(context.Background(),5 * time.Second)
  defer shutdownCancel()
  err = rtr.HTTPSvr.Shutdown(shutdownCtx)
  if err != nil {
    log.Fatalf("[-] HTTP Server shutdown error: %s\n",err.Error())
  }
  if rtr.Tls {
    new_err := rtr.HTTPSSvr.Shutdown(shutdownCtx)
    if new_err != nil {
      log.Fatalf("[-] HTTPS Server shutdown error: %s\n",new_err.Error())
    }
  }
  fmt.Println("Server are off")
  log.Println("[+] Server gracefully stopped.")
}
