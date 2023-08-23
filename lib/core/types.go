package core

/*
  A bunch of exported types and interfaces for whatever reason deemed fit.
  Do not import internal stuff in this package space to prevent cross importation
  Anything defined here is exported to everything odin
*/
import (
  "net"
)

type Server interface{
  Serve()
}

type ServerClient interface  {
  Serve() (*net.Listener,error)
  WClient() (*net.Conn,error)
}

type Work struct{
  UserId string
  OperatorId string
  Screenshot []string
  FD FileData
  CmdIn string
  CmdOut string
  Individual bool
  ForMS bool
}

type FileData struct {
  Nme string
  Data []byte
  Run bool
}
