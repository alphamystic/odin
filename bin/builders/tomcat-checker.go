package main

import (
  "fmt"
  "net/http"
  "crypto/tls"
  "github.com/alphamystic/odin/plugins"
  "github.com/alphamystic/odin/lib/utils"
)


var Users = []string{"admin","manager","tomcat"}
var Passwords = []string{"admin","manager","tomcat","passwords"}

type TomcatChecker struct{}

func (tc *TomcatChecker) Check(host string,ssl bool,port uint64) *plugins.Result{
  var (
    resp *http.Response
    err error
    url string
    res *plugins.Result
    client *http.Client
    req *http.Request
  )
  res = new(plugins.Result)
  client = new(http.Client)
  if ssl {
		client = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		}
		url = fmt.Sprintf("https://%s:%d/manager/html", host)
	} else {
		client = &http.Client{}
		url = fmt.Sprintf("http://%s:%d/manager/html",host,port)
	}
  if resp,err = http.Head(url);err != nil {
    utils.Logerror(fmt.Errorf("Headerrequest failed: %v",err))
    return res
  }
  utils.PrintTextInASpecificColor("cyan","Host responded to /manager/html request")
  if resp.StatusCode != http.StatusUnauthorized || resp.Header.Get("www-Authenticate") == ""{
    utils.PrintTextInASpecificColor("yellow","Target doe not appear to require basic auth.")
    return res
  }
  if req,err = http.NewRequest("GET",url,nil); err != nil{
    utils.Logerror(fmt.Errorf("Unable to build get request. %v",err))
    return res
  }
  for _, user := range Users{
    for _,pass := range Passwords {
      req.SetBasicAuth(user,pass)
      if resp,err = client.Do(req); err != nil{
        utils.Logerror(fmt.Errorf("Unable to get request. %v",err))
        return res
      }
      if resp.StatusCode == http.StatusOK {
        res.Vulnerable = true
        res.Details = fmt.Sprintf("Valid credentials found. Username: %s. Password: %s",user,pass)
        return res
      }
    }
  }
  return res
}

func New()plugins.Checker{
  return new(TomcatChecker)
}
