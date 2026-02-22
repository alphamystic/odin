package api


import (
  "os"
  "log"
  "fmt"
  "time"
  "sync"
  "syscall"
  "context"
  "net/http"
  "os/signal"
  "crypto/tls"
  "github.com/alphamystic/odin/lib/utils"
  //dfn"github.com/alphamystic/odin/lib/definers"
  "github.com/alphamystic/odin/loki/api/apihandlers"
)

// A low Level router exposing the default http

type Router struct {
  Mux *http.ServeMux
  HTTPSvr *http.Server
  HTTPSSvr *http.Server
  Tls bool
	Mode string // DEV or PROD
}

// should probably receive a server
func NewRouter(httpsSvr, httpSvr *http.Server, mode string) *Router {
	return &Router{
		Mux:      http.NewServeMux(),
		HTTPSvr:  httpSvr,
		HTTPSSvr: httpsSvr,
		Mode:  mode,
	}
}


type APIServers struct {
  Address string
  PortS int
  Port int
  TlsCert string
  TlsKey string
  Tls bool
  ApiKey string // servers api keey to chat service at main
}

func (api_server *APIServers) CreateServer() (*http.Server,*http.Server) {
  httpServer := &http.Server {
    Addr: fmt.Sprintf(":%d", api_server.Port),
    ReadTimeout: 5 * time.Second,
    WriteTimeout: 10 * time.Second,
    IdleTimeout: 120 * time.Second,
	}
  config := &tls.Config {
    MinVersion: tls.VersionTLS12,
    CurvePreferences: []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
    PreferServerCipherSuites: true,
    CipherSuites: []uint16 {
      tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
      tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
      tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
      tls.TLS_RSA_WITH_AES_256_CBC_SHA,
    },
  }
  httpsServer := &http.Server {
    Addr: fmt.Sprintf(":%d",api_server.PortS),
    TLSConfig: config,
    TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
    ReadTimeout: 5 * time.Second,
    WriteTimeout: 10 * time.Second,
    IdleTimeout: 120 * time.Second,
  }
  return httpServer,httpsServer
}

func (rtr *Router) RunAPI(reg bool){
  rtr.HTTPSvr.Handler = rtr.Mux
  rtr.HTTPSSvr.Handler = rtr.Mux
  // start channels to write logs
  ShutdownCh := make(chan bool)
  DoneCh := make(chan bool)
  var err error
  rl := utils.NewRequestLogger("./.logs/requests/",0644)
  api_hnd, err := apihandlers.NewAPIHandler(ShutdownCh, DoneCh, rl,rtr.Mode)
	if err != nil {
		utils.Danger(err)
		return
	}

  //*********** ROUTES **********//
  // ********** AUTH LOGS *******//
  rtr.Mux.HandleFunc("/api/login",api_hnd.Authenticateuser)
  rtr.Mux.HandleFunc("/api/logout",api_hnd.Logout)
  // **** END OF AUTH ROUTES *********//

  // *****USER ROUTES ***********
  rtr.Mux.HandleFunc("/api/users/createuser",api_hnd.Createuser)
  rtr.Mux.HandleFunc("/api/users/listusers",api_hnd.Listusers)
  rtr.Mux.HandleFunc("/api/users/listadmins",api_hnd.Listadmins)
  // ******** END OF USER ROUTES ********//


  // ******** START OF AUTH ROUTES *********

  // *********** END OF AUTH ROUTES *******//


  // ************ START OF RECON ROUTES ***********
  rtr.Mux.HandleFunc("/api/recon/createscan/",api_hnd.WithUserData(api_hnd.Createscan))
  rtr.Mux.HandleFunc("/api/recon/listscans/",api_hnd.WithUserData(api_hnd.Listscans))
  rtr.Mux.HandleFunc("/api/recon/viewscans/",api_hnd.WithUserData(api_hnd.Viewscan))
  rtr.Mux.HandleFunc("/api/recon/createtarget",api_hnd.WithUserData(api_hnd.Createtarget))
  rtr.Mux.HandleFunc("/api/recon/viewtarget/",api_hnd.WithUserData(api_hnd.Viewtarget))
  rtr.Mux.HandleFunc("/api/recon/createrecondata",api_hnd.WithUserData(api_hnd.Createrecondata))
  rtr.Mux.HandleFunc("/api/recon/createservice",api_hnd.WithUserData(api_hnd.Createservice))
  rtr.Mux.HandleFunc("/api/health",api_hnd.Health)
  // ************* END OF RECON ROUTES ************//


  // ******** START OF AUTH ROUTES *********
  rtr.Mux.HandleFunc("/api/asssets/createasset",api_hnd.WithUserData(api_hnd.Createscan))
  rtr.Mux.HandleFunc("/api/asssets/viewasset/",api_hnd.WithUserData(api_hnd.Createscan))
  rtr.Mux.HandleFunc("/api/asssets/listassets/",api_hnd.WithUserData(api_hnd.Createscan))
  rtr.Mux.HandleFunc("/api/asssets/listassetsbyfilter/",api_hnd.WithUserData(api_hnd.Createscan))
  // *********** END OF AUTH ROUTES *******//

  // ******** START OF Asssets ROUTES *********
  // *********** END OF Assets ROUTES *******//

  // ******** START OF Blog ROUTES *********
  rtr.Mux.HandleFunc("/api/recon/createblog",api_hnd.WithUserData(api_hnd.CreateBlog))
  rtr.Mux.HandleFunc("/api/recon/viewblog/",api_hnd.WithUserData(api_hnd.ViewBlog))
  rtr.Mux.HandleFunc("/api/recon/createscan/",api_hnd.WithUserData(api_hnd.GetRecentBlogs))
  rtr.Mux.HandleFunc("/api/recon/createscan/",api_hnd.WithUserData(api_hnd.ListBlogsByAuthor))
  rtr.Mux.HandleFunc("/api/recon/createscan/",api_hnd.WithUserData(api_hnd.ListBlogsByMainTag))
  rtr.Mux.HandleFunc("/api/recon/createscan/",api_hnd.WithUserData(api_hnd.ListBlogsByType))
  rtr.Mux.HandleFunc("/api/recon/createscan/",api_hnd.WithUserData(api_hnd.ListBlogsByCategory))
  rtr.Mux.HandleFunc("/api/recon/createscan/",api_hnd.WithUserData(api_hnd.UpdateBlog))
  rtr.Mux.HandleFunc("/api/recon/createscan/",api_hnd.WithUserData(api_hnd.ArchiveBlog))
  rtr.Mux.HandleFunc("/api/recon/createscan/",api_hnd.WithUserData(api_hnd.CreateComment))
  rtr.Mux.HandleFunc("/api/recon/createscan/",api_hnd.WithUserData(api_hnd.ListComments))

  // *********** END OF Blog ROUTES *******//

  // ******** START OF Blog Analytics ROUTES *********
  // *********** END OF Blog Analytics ROUTES *******//

  // ******** START OF Mothership ROUTES *********
  // *********** END OF Mothership ROUTES *******//

  // ******** START OF Minion ROUTES *********
  // *********** END OF Minion ROUTES *******//

  // ******** START OF User Manager ROUTES *********
  // *********** END OF User Manager ROUTES *******//

  // ******** START OF APIKey Manager ROUTES *********
  // *********** END OF APIKey Manager ROUTES *******//





 // **** Middleware functionalities **** //

 // ********* END OF ROUTES *********** //

 // WaitGroup to wait for the server goroutines to finish
  var wg sync.WaitGroup
  wg.Add(2) // We will wait for 2 servers (HTTP and HTTPS)
  go func(){
    defer wg.Done()
    if err := rtr.HTTPSSvr.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("API HTTP server error: %v", err)
		}
  }()
  if rtr.Tls {
    go func(){
      defer wg.Done()
      // we need to find a better way of supplying this
      //if err := rtr.HTTPSSvr.ListenAndServeTLS("../../../certs/server.crt", "../../../certs/server.key"); err != http.ErrServerClosed {
      if err := rtr.HTTPSvr.ListenAndServeTLS("./certs/server.crt", "./certs/server.key"); err != nil && err != http.ErrServerClosed {
				log.Fatalf("API HTTPS server error: %v", err)
			}
    }()
  }
  utils.PrintInformation("API Servers are here running")
  interruptChan := make(chan os.Signal,1)
  signal.Notify(interruptChan,os.Interrupt, syscall.SIGTERM)
  //sedn a close channel to the handler
  api_hnd.ShutdownChan <- true
  // wait for the receiver to finish writing all logs
  <-api_hnd.DoneChan
  // read from the interrupt chan and shutdown
  <-interruptChan
  log.Println("Received shutdown signal, shutting down servers...")
  shutdownCtx,shutdownCancel := context.WithTimeout(context.Background(),5 * time.Second)
  defer shutdownCancel()
  err = rtr.HTTPSvr.Shutdown(shutdownCtx)
  if err != nil {
    log.Fatalf("[-] API HTTP Server shutdown error: %s\n",err.Error())
  }
  if rtr.Tls {
    new_err := rtr.HTTPSSvr.Shutdown(shutdownCtx)
    if new_err != nil {
      log.Fatalf("[-] API HTTPS Server shutdown error: %s\n",new_err.Error())
    }
  }
   wg.Wait()
  fmt.Println("API API Servers are off.............")
  log.Println("[+] API Server gracefully stopped.")
}
