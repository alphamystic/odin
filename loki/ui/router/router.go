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
  //dfn"github.com/alphamystic/odin/lib/definers"
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

  // create a request logger
  rl := utils.NewRequestLogger("./.data/logs/requests/",066)

  // initiate new handler
  hnd,err := handlers.NewHandler(ShutdownCh, DoneCh, rl,"https://localhost:5001","12345")
  if err != nil {
    utils.Danger(err);return
  }

  // Handlers
  // Portfolio Handler
  //rtr.Mux.HandleFunc("/",hnd.Portfolio)
  rtr.Mux.HandleFunc("/health",hnd.Health)
  rtr.Mux.HandleFunc("/portfolio",hnd.Portfolio)
  rtr.Mux.HandleFunc("/blogs",hnd.BlogSite)
  rtr.Mux.HandleFunc("/createblog",hnd.Createblog) /// this renders perfect, use this to reports
  rtr.Mux.HandleFunc("/writeblog",hnd.Writeblog)

  // Odin Net
  rtr.Mux.HandleFunc("/odin-net",hnd.OdinNet)

  //panel shortcuts (We handle the rest as they come. bit by bit baba)
  rtr.Mux.HandleFunc("/profile",hnd.Profile)
  rtr.Mux.HandleFunc("/updateprofile",hnd.Updateprofile)
  rtr.Mux.HandleFunc("/securityprofile",hnd.Securityprofile)
  rtr.Mux.HandleFunc("/notificationprofile",hnd.Notificationsprofile)

  //panel side
  rtr.Mux.HandleFunc("/test",hnd.Blank)
  rtr.Mux.HandleFunc("/",hnd.Home) //.Home
  rtr.Mux.HandleFunc("/mkubwa",hnd.Signin)
  rtr.Mux.HandleFunc("/signout",hnd.Logout)
  rtr.Mux.HandleFunc("/register",hnd.Register)

  rtr.Mux.HandleFunc("/apt",hnd.Apt)
  rtr.Mux.HandleFunc("/edr",hnd.Edr)

  rtr.Mux.HandleFunc("/bb",hnd.Bugbounty) //lists
  rtr.Mux.HandleFunc("/pentests",hnd.Pentests) //lists
  rtr.Mux.HandleFunc("/bo",hnd.Blackops)

  //Scan tests
  rtr.Mux.HandleFunc("/scans",hnd.RenderScans) //lists
  rtr.Mux.HandleFunc("/scans/list-targets",hnd.ListTargets) //lists
  rtr.Mux.HandleFunc("/scans/view-target",hnd.ViewTarget)

  rtr.Mux.HandleFunc("/bbreports",hnd.BugBountyReports)//rendering
  rtr.Mux.HandleFunc("/createreport", hnd.ManageBlog)
  rtr.Mux.HandleFunc("/reports/import", hnd.ParseDocument)
  rtr.Mux.HandleFunc("/listreports",hnd.RenderBlogList) // use this to list various blogs(templates)
  rtr.Mux.HandleFunc("/viewsreports",hnd.ViewReport)
  rtr.Mux.HandleFunc("/pentestsreport",hnd.PentestsReports) //rendering

  rtr.Mux.HandleFunc("/ms",hnd.ListMotherships)

  // Asset manager
  rtr.Mux.HandleFunc("/assets/list", hnd.ListAssets)
  rtr.Mux.HandleFunc("/agents",hnd.Minions)


  rtr.Mux.HandleFunc("/pendingscans",hnd.Pendingscans)
  rtr.Mux.HandleFunc("/phishinglinks",hnd.Phishinglinks)
  rtr.Mux.HandleFunc("/zerodays",hnd.Zerodays)

  // Adding a report manager to create an aggregate functionality
  rtr.Mux.HandleFunc("/reports/create",hnd.ListFiles)

  rtr.Mux.HandleFunc("/listfiles",hnd.ListFiles)

  rtr.Mux.HandleFunc("/events",hnd.Events)
  rtr.Mux.HandleFunc("/regulars",hnd.RegularUsers)
  rtr.Mux.HandleFunc("/admins",hnd.Admins)
  rtr.Mux.HandleFunc("/activeprojects",hnd.Activeprojects)
  rtr.Mux.HandleFunc("/archivedprojects",hnd.Archivedprojects)

  /* Yara Handlers */
  rtr.Mux.HandleFunc("/listyararules",hnd.ListYaraRule)
  rtr.Mux.HandleFunc("/createyararule",hnd.CreateYaraRule)

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
  	if err := rtr.HTTPSSvr.ListenAndServe(); err != http.ErrServerClosed {
  		log.Fatalf("HTTP server error: %v", err)
  	}
  }()
  if rtr.Tls {
    go func(){
      // we need to find a better way of supplying this
      if err := rtr.HTTPSvr.ListenAndServeTLS("./certs/server.crt", "./certs/server.key"); err != http.ErrServerClosed {
  			log.Fatalf("HTTPS server error: %v", err)
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
