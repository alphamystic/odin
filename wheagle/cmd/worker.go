package cmd

/*
  Will add things rather operator configs as they come along
*/
import (
  "os"
  "fmt"
  "github.com/alphamystic/odin/lib/core"
  "github.com/alphamystic/odin/lib/utils"
)

type RuntimeWorker struct {
  OperatorId string
  Cave *core.Cave
  Dir string
}

func CreateRW()*RuntimeWorker{
  return &RuntimeWorker{
    OperatorId: "3l0r@cle",
    Cave: core.InitializeTunnelMan(),
    Dir: "/home/sam/Documents/3l0racle/github.com/alphamystic/odin/wheagle",
  }
}

func (rw *RuntimeWorker) CreateWheagleDir(dir string) error{
  if !utils.CheckifStringIsEmpty(dir){
    rw.Dir = dir;return nil
  }
  dir,err := os.UserHomeDir()
  if err != nil {
    return fmt.Errorf("Error getting home directory. \n%q",err)
  }
  rw.Dir = dir + "./odin"
  return nil
}
