package main

import (
  "github.com/alphamystic/odin/plugins"
)

type CNC struct{}

func (cnc *CNC) Drop() (func(string, ...string)error){
  return Implementor()
}

func NewDropper() plugins.Dropper{
  return new(CNC)
}

func Implementor()func(string, ...string)error{
  return func(command string,args ...string)error{
    switch command {
      case "ia":
        fmt.Println("Interacting with implant")
      case "ci":
        fmt.Println("Creating implant")
      default:
        return errors.New("Unknown arguement")
    }
    return nil
  }
}
