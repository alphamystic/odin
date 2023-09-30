package main

/*
  * Load functions from private
*/

import (
  "fmt"
  "net/http"
  "github.com/alphamystic/odin/lib/handlers"
)
//we should definietlty change this or give it something universal for proper private loaading
func NewPrivateScanner() hadlers.Scanner{
  return new(Private)
}

type Private struct{}

func (p *Private) Scan(rd handlers.ReconData) []handlers.Vulnerability{
  var vulns []handlers.Vulnerability
  // create a channel that takes in al the various vulnerabilities and writes them into a vulnerabilities channel
  //les't segregate recon data and know what to call where.
  return vulns
}

func HeaderEnumerator(ssl bool)error{
  return nil
}
// bruteforce all web logins (if cms then ....)
func Login(ssl bool,dir string)error{
  var (
    err error
    url string
    client *http.Client
  )
  client = new(http.Client)
  if ssl {
		client = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		}
		url = fmt.Sprintf("https://%s:%d/"+dir, trg.TargetIp.String(),)
	} else {
		client = &http.Client{}
		url = fmt.Sprintf("http://%s:%d/manager/html",host,port)
	}
  return nil
}

func BruteForceLogin()error{
  return nil
}

func CreateClient(s)
