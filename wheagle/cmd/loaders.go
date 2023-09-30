package cmd

import (
  "fmt"
  "github.com/alphamystic/odin/plugins"
  "github.com/alphamystic/odin/lib/utils"
  "github.com/spf13/cobra"
)

var cmdLoadScanner = &cobra.Command {
	Use: "load-scanner",
	Short: "Load a vulnerability scanner and test it",
	Run: func(cmd *cobra.Command,args []string){
		/*name,err := cmd.Flags().GetString("plugin")
		if err != nil {
			utils.PrintTextInASpecificColorInBold("red","	Unable to load plugin")
			return
		}
		pluginHandler(name)*/
	},
}

var cmdLoadDropper = &cobra.Command {
	Use: "load-dropper",
	Short: "Load a dropper/custom made C2",
	Run: func(cmd *cobra.Command,args []string){
		/*name,err := cmd.Flags().GetString("plugin")
		if err != nil {
			utils.PrintTextInASpecificColorInBold("red","	Unable to load plugin")
			return
		}
    err,drp,implementor := LoadDropper(name)
    if err != nil{
      utils.Logerror(err)
      return
    }
    for _,cmnd := range drp{
      for _,arg := range cmnd.Arguements{
        utils.PrintTextInASpecificColorInBold("yellow",fmt.Sprintf(" ",arg))
      }
    }
    var iarg string
    for {
      fmt.Scanln(&iarg)
      implementor("",iarg)
      /*START:
        fmt.Scanln(&iarg)
        switch iarg {
        case "back":
  				goto END
        default:
          implementor(iarg,iarg)
        }
      END:
        fmt.Println("[-+-]    Switching back to ODIN shell.")
        return
      }*/
	},
}

var cmdLoadPoncupine = &cobra.Command {
	Use: "load-poncupine",
	Short: "Load a dropper/custom made C2",
	Run: func(cmd *cobra.Command,args []string){
    iF,err := cmd.Flags().GetString("iF")
		if err != nil {
			utils.PrintTextInASpecificColor("red","Input file can not be empty")
			return
		}
    oF,err := cmd.Flags().GetString("oF")
		if err != nil {
			utils.PrintTextInASpecificColor("red","Output file can not be empty")
      utils.Logerror(err)
			return
		}
    frmt,err := cmd.Flags().GetString("f")
		if err != nil {
			utils.PrintTextInASpecificColor("red","Input file can not be empty")
			return
		}
    plgns := &plugins.PluginLoader{
      Name: "poncupine.so",
    }
    fmt.Println("Data to beloaded : ",iF,oF,frmt)
    err = plgns.LoadRodent(iF,oF,frmt)
    if err != nil{
      utils.Logerror(err);return
    }
    utils.PrintTextInASpecificColor("white","Gracefully exited poncupine ..............")
    fmt.Println("There was no error")
	},
}
func init(){
  //Poncupine vars
  cmdLoadPoncupine.Flags().String("f","iso","Format for poncupine to implement.")
  cmdLoadPoncupine.Flags().String("iF","payload.exe","Input file to modify.")
  cmdLoadPoncupine.Flags().String("oF","modified_payload.exe","Output file to be written into.")
  RootCmd.AddCommand(cmdLoadPoncupine)
  RootCmd.AddCommand(cmdLoadDropper)
  RootCmd.AddCommand(cmdLoadScanner)
}
