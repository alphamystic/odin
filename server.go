package main

import (

  "github.com/alphamystic/odin/loki"
  "github.com/alphamystic/odin/loki/ui/router"
)
func main(){
  Loki := &loki.Loki {
    Address: "0.0.0.0",
    PortS: 4001,
    //Port: 4000,
    TlsCert: "",
    TlsKey: "",
    Tls: true,
    ApiKey: "", // servers api keey to chat service at main
  }
  httpSvr, httpsSvr := Loki.CreateServer()
  rtr := router.NewRouter(httpSvr, httpsSvr)
  rtr.Tls = true
  rtr.Run(true)
}
