package main

import (
  "fmt"
  "flag"
  "github.com/alphamystic/odin/loki/api"
)


func main() {
	// Define CLI flags
	mode := flag.String("mode", "DEV", "Run mode: DEV or PROD")
	ports := flag.Int("ports", 5001, "Port to listen on for HTTPS")
  port := flag.Int("port", 5000, "Port to listen on for HTTP")
	tls := flag.Bool("tls", true, "Enable TLS (true/false)")

	flag.Parse()

	fmt.Printf("Starting API server in %s mode on port %d (TLS=%v)\n", *mode, *port, *tls)

	api_server := &api.APIServers{
		Address: "0.0.0.0",
		PortS:   *ports,
    Port:     *port,
		Tls:     *tls,
		TlsCert: "",
		TlsKey:  "",
		ApiKey:  "",
	}

	httpSvr, httpsSvr := api_server.CreateServer()
	rtr := api.NewRouter(httpSvr, httpsSvr, *mode)
	rtr.Tls = *tls
	rtr.RunAPI(true)
}
