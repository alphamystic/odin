// main.go
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"
	"strings"

	"github.com/alphamystic/odin/marketplace/reconjs"
)

func usageAndExit() {
	fmt.Println("Usage:")
	fmt.Println("  reconjs -mode <js|endpoints|secrets|check|fuzz|idor|all> -domain <domain> [-list <file>] [-outdir <dir>]")
	fmt.Println("Examples:")
	fmt.Println("  reconjs -mode js -domain example.com")
	fmt.Println("  reconjs -mode all -domain example.com")
	os.Exit(1)
}

func main() {
	mode := flag.String("mode", "", "one of: js | endpoints | secrets | check | fuzz | idor | all")
	domain := flag.String("domain", "", "target domain (e.g. example.com)")
	list := flag.String("list", "", "file with domains (one per line)")
	outdir := flag.String("outdir", ".", "directory to write JSON files")
	flag.Parse()

	if *mode == "" {
		usageAndExit()
	}
	var targets []string
	if *domain != "" {
		targets = append(targets, *domain)
	}
	if *list != "" {
		// read file
		data, err := os.ReadFile(*list)
		if err != nil {
			log.Fatalf("failed to read list file: %v", err)
		}
		for _, l := range strings.Split(string(data), "\n") {
			l = strings.TrimSpace(l)
			if l == "" {
				continue
			}
			targets = append(targets, l)
		}
	}
	if len(targets) == 0 {
		log.Println("no targets provided")
		usageAndExit()
	}

	// iterate targets
	for _, t := range targets {
		log.Printf("=== processing target: %s (mode=%s) ===", t, *mode)

		// create ReconData
		rData := reconjs.ReconData{
			Domain: t,
			Meta:   map[string]any{"started_at": fmt.Sprintf("%v", nowUTC())},
		}

		switch *mode {
		case "js":
			js, err := reconjs.DiscoverJSForDomain(t)
			if err != nil {
				log.Printf("discover js error: %v", err)
			}
			rData.JSFiles = js
			_ = reconjs.WriteJSONSafe(*outdir+"/js_"+safeFilename(t)+".json", rData)
			log.Printf("wrote js results for %s to %s", t, *outdir+"/js_"+safeFilename(t)+".json")
		case "endpoints":
			// First discover js (we can rely on user providing js file list, but for simplicity discover)
			js, _ := reconjs.DiscoverJSForDomain(t)
			rData.JSFiles = js
			endpoints, _ := reconjs.ExtractEndpointsFromJS(js, 6)
			rData.Endpoints = endpoints
			_ = reconjs.WriteJSONSafe(*outdir+"/endpoints_"+safeFilename(t)+".json", rData)
			log.Printf("wrote endpoints results for %s", t)
		case "secrets":
			js, _ := reconjs.DiscoverJSForDomain(t)
			rData.JSFiles = js
			secrets, _ := reconjs.ExtractSecretsFromJS(js, 6)
			rData.Secrets = secrets
			_ = reconjs.WriteJSONSafe(*outdir+"/secrets_"+safeFilename(t)+".json", rData)
			log.Printf("wrote secrets results for %s", t)
		case "check":
			js, _ := reconjs.DiscoverJSForDomain(t)
			rData.JSFiles = js
			endpoints, _ := reconjs.ExtractEndpointsFromJS(js, 6)
			rData.Endpoints = endpoints
			httpChecks, _ := reconjs.CheckEndpointsWithHttpx(endpoints)
			rData.HTTPChecks = httpChecks
			_ = reconjs.WriteJSONSafe(*outdir+"/httpchecks_"+safeFilename(t)+".json", rData)
		case "fuzz":
			// example fuzz run, you can adjust wordlist path and extra args
			res, _ := reconjs.RunFFUF("https://"+t, "/usr/share/wordlists/dirb/others/best1050.txt", "-mc 200 -e .php,.bak,.zip,.js,.cgi,.png,.jpeg,.bat,.asp,.aspx,.dll,.htm,.html,.jsp,.txt,.xml,.sql,.phtml,.inc,.jsa")
			rData.FuzzResults = []reconjs.FFUFResult{res}
			_ = reconjs.WriteJSONSafe(*outdir+"/fuzz_"+safeFilename(t)+".json", rData)
		case "idor":
			js, _ := reconjs.DiscoverJSForDomain(t)
			rData.JSFiles = js
			endpoints, _ := reconjs.ExtractEndpointsFromJS(js, 6)
			rData.Endpoints = endpoints
			idors, _ := reconjs.SimpleIDORTest(endpoints, "1234", 10)
			rData.IDORChecks = idors
			_ = reconjs.WriteJSONSafe(*outdir+"/idor_"+safeFilename(t)+".json", rData)
		case "all":
			// run pipeline: discover -> endpoints -> secrets -> httpx -> idor (fuzz optional)
			js, _ := reconjs.DiscoverJSForDomain(t)
			rData.JSFiles = js

			endpoints, _ := reconjs.ExtractEndpointsFromJS(js, 8)
			rData.Endpoints = endpoints

			secrets, _ := reconjs.ExtractSecretsFromJS(js, 8)
			rData.Secrets = secrets

			httpChecks, _ := reconjs.CheckEndpointsWithHttpx(endpoints)
			rData.HTTPChecks = httpChecks

			idors, _ := reconjs.SimpleIDORTest(endpoints, "1234", 10)
			rData.IDORChecks = idors

			_ = reconjs.WriteJSONSafe(*outdir+"/all_"+safeFilename(t)+".json", rData)
			log.Printf("wrote all results for %s", t)
		default:
			log.Printf("unknown mode %s", *mode)
			usageAndExit()
		}
	}
}

//
// minor helpers for main.go
//
func safeFilename(s string) string {
	// replace chars unsuitable for filename
	out := strings.ReplaceAll(s, ":", "_")
	out = strings.ReplaceAll(out, "/", "_")
	out = strings.ReplaceAll(out, " ", "_")
	return out
}

func nowUTC() string {
	return fmt.Sprintf("%v", timeNowUTC())
}

// small wrapper around time to make testing easy (and keep imports local)
func timeNowUTC() any {
	return strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(fmt.Sprint(time.Now().UTC()), " ", "_"), ":", "-"))
}
