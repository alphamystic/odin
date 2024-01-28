package loki

/*
  * Exporter for loki server
  Exports the UI and API Server(will be implemented later)
*/

import (
  rtr"github.com/alphamystic/loki/ui/router"
)

type Loki struct {
  Address string
  PortS int
  Port int
  TlsCert string
  TlsKey string
  Tls bool
  ApiKey string // servers api keey to chat service at main
}

// we'll orbably need a struct to hold if HTTP/HTTPS
func (l *Loki) Server(address string,port int){
  // start the main server or the api server.
  //"0.0.0.0:9000"
  if l.Tls {
    // write tls work here
  } else {

  }

  addr := fmt.Sprintf("%s:%d",address,port)
  rtr := router.NewRouter(addr)
  rtr.Run(true)
}



// https://pkg.go.dev/golang.org/x/net/http2/h2demo#section-readme
func (l *Loki) CreateServer() (*http.Server,*http.Server) {
  httpServer := &http.Server {
    Addr: fmt.Sprintf(":%d", l.Port),
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
    Addr: fmt.Sprintf(":%d",l.PortS),
    TLSConfig: config,
    TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
  }
  return httpServer,httpsServer
}

// openssl ecparam -genkey -name secp384r1 -out server.key
// openssl req -new -x509 -sha256 -key server.key -out server.crt -days 3650
