package cmd

import (
  "fmt"
  "time"
  "sync"
  "github.com/alphamystic/odin/lib/utils"
  "github.com/spf13/cobra"
)
/*
  * Runs a configurator for a default operator profile or a mothership in the background
*/

type Mothership struct{
  Name string
  ID string
  startTime string
  PID int
}

type MSRunner struct{
  MS map[string]*Mothership
  mu sync.RWMutex
}

var MSR = new(MSRunner)

func NewMSR()*MSRunner{
  return &MSRunner{
    MS: make(map[string]*Mothership),
    mu: sync.RWMutex{},
  }
}
func (rm *MSRunner) runMS(id,name,pathToExec string)error{
  rm.mu.Lock()
	defer rm.mu.Unlock()
  pid,err := utils.RunExecutable(pathToExec)
  if err != nil{
    return fmt.Errorf("Error starting Mothership: %s with id %s.\nERROR: %q",name,id,err)
  }
  ms := &Mothership{
    Name:name,
    ID: id,
    startTime: StartTime.Format(time.UnixDate),
    PID: pid,
  }
  rm.MS[ms.ID] = ms
  utils.PrintTextInASpecificColor("cyan",fmt.Sprintf("Started mothership %s with id %s and Pid value of %d",name,id,pid))
  return nil
}

func (rm *MSRunner) killMs(id string) error{
  rm.mu.Lock()
	defer rm.mu.Unlock()
  for _,ms := range rm.MS{
    if ms.ID == id{
      err := utils.KillExec(ms.PID)
      if err != nil {
        return fmt.Errorf("Error killing mothership wuthid of: %s\nERROR: %q",id,err)
      }
      delete(rm.MS,ms.ID)
      utils.PrintTextInASpecificColor("blue","Successfully terminated mothership.")
      return nil
    }
  }
  return fmt.Errorf("Error no thothership with id of %s ",id)
}


var cmdKillMS = &cobra.Command{
  Use: "kill-ms",
  Short: "Kill a particular running mothership",
  Long: "",
  Run: func(cmd *cobra.Command, args []string){
    id, err := cmd.Flags().GetString("id")
		if err != nil {
			utils.PrintTextInASpecificColor("red","ID can not be empty")
			return
		}
    if err := MSR.killMs(id); err != nil{
      utils.Logerror(err);return
    }
  },
}

var cmdStartMS = &cobra.Command{
  Use: "start-ms",
  Short: "Run a mothership in the background.",
  Long: "",
  Run: func(cmd *cobra.Command, args []string){
    id, err := cmd.Flags().GetString("id")
		if err != nil {
			utils.PrintTextInASpecificColor("red","ID can not be empty")
			return
		}
    name, err := cmd.Flags().GetString("name")
		if err != nil {
			utils.PrintTextInASpecificColor("red","name can not be empty")
			return
		}
    path, err := cmd.Flags().GetString("path")
		if err != nil {
			utils.PrintTextInASpecificColor("red","path can not be empty")
			return
		}
    if err := MSR.runMS(id,name,path); err != nil{
      utils.Logerror(err);return
    }
  },
}

func init(){
  MSR = NewMSR()
  cmdStartMS.Flags().String("id","[id]","ID for the mothership(c2)")
  cmdStartMS.Flags().String("name","name","Name of your c2(Preferably use the one you used while creating it)")
  cmdStartMS.Flags().String("path","../bin/temp/mothership","Path to the mS Executable.")
  cmdKillMS.Flags().String("id","[id]","ID for the mothership(c2)")
  RootCmd.AddCommand(cmdStartMS)
  RootCmd.AddCommand(cmdKillMS)
}
