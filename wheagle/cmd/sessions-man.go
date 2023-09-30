package cmd

import (
  "fmt"
  //"time"

  "github.com/alphamystic/odin/lib/c2"
  "github.com/alphamystic/odin/lib/db"
  "github.com/alphamystic/odin/lib/utils"

  "github.com/spf13/cobra"
)
/*should manage the operator cli communicating to implant/adminc2
1. When Generate admin/mule is called, should initialize it into the sessions
2. When sending commands to mules, determine which adminC2 is being written to
3. When interactive determine the C2 in Question comm to the mule
4. If delete mule, clearout it's session in particular
5. if close MS, give a new MS or sleep out mules to MS too
*/

// get the db driver
var SessionsDriver = new(db.Driver)
var ConnectorsDriver = new(db.Driver)

func init(){
  var err error
  SessionsDriver,err = db.Old(c2.SessionsPath,0644)
  if err != nil {
    utils.Logerror(fmt.Errorf("Error loading a db driver for sessions"))
    panic(err)
  }
  ConnectorsDriver,err = db.Old(c2.ConnectorsPath,0644)
  if err != nil {
    utils.Logerror(fmt.Errorf("Error loading a db driver for connectors"))
    panic(err)
  }
}
// this are sessions belongning to the minions
var  RunningSessions = c2.InitializeNewSessionManager()

// This are sessoins/connections to the various C2 Mother ships
var Conns = c2.InitializeNewConnectionMan()

// initialize  connectors on startup too
func init(){
  // load sessions
  err := RunningSessions.LoadSessions(SessionsDriver)
  if err != nil {
    //if err log it out and set sessions to an empty area
    utils.Warning(fmt.Sprintf("%s",err))
  }
  //print all available sessions
  RunningSessions.ListSessions()

  fmt.Println(" ")
  fmt.Println(" ")

  //Load connectors
  err =  Conns.LoadConnectors(ConnectorsDriver)
  if err != nil {
    utils.Warning(fmt.Sprintf("%s",err))
  }
  //Print avaiable connectors on startup
  Conns.ListConnectors()

  fmt.Println(" ")
  fmt.Println(" ")
}

var cmdLister = &cobra.Command{
  Use: "list",
  Short: "List all the active mules or connections",
  Long: "",
  Run: func(cmd *cobra.Command, args []string){
    for _,arg := range args {
      switch arg {
      case "minions":
          no := len(RunningSessions.Sessions)
          utils.PrintTextInASpecificColorInBold("magenta",fmt.Sprintf("Total number of mule sessions are %d",no))
          RunningSessions.ListSessions()
          fmt.Println("")
        case "admins":
          ano := len(Conns.Connections)
          utils.PrintTextInASpecificColorInBold("magenta",fmt.Sprintf("Total number of admin connections are %d",ano))
          Conns.ListConnectors()
          fmt.Println("")
        default:
        utils.PrintTextInASpecificColor("blue","Please specify what listing you want (list admins or list mules)")
      }
    }
  },
}

var cmdDeleteMS  = &cobra.Command{
  Use: "delete-ms",
  Short: "Delete a given mothership. This only deletes it from the db and connection manager not from the location(file/server)",
  Run: func(cmd *cobra.Command, args []string){
    id,err := cmd.Flags().GetString("id")
		if err != nil {
			utils.PrintTextInASpecificColor("red","	Error id cannot be empty.")
			return
		}
    if err := ConnectorsDriver.Delete("connectors",id); err != nil{
      utils.PrintTextInASpecificColor("red",fmt.Sprintf("Error deleting ms with id %s.\n ERROR: %s",id,err))
      return
    }
    utils.PrintTextInASpecificColor("yellow",fmt.Sprintf("Deleted MotherShip with id %s",id))
  },
}

var cmdDeleteMinion = &cobra.Command{
  Use: "delete-minion",
  Short: "Delete a given minion. This only deletes it from the db and session manager not from the location(implant)",
  Run: func(cmd *cobra.Command, args []string){
    id,err := cmd.Flags().GetString("id")
		if err != nil {
			utils.PrintTextInASpecificColor("red","	Error id cannot be empty.")
			return
		}
    if err := SessionsDriver.Delete("sessions",id); err != nil{
      utils.PrintTextInASpecificColor("red",fmt.Sprintf("Error deleting session with id %s.\n ERROR: %s",id,err))
      return
    }
    utils.PrintTextInASpecificColor("yellow",fmt.Sprintf("Deleted session with id %s",id))
  },
}

func init(){
  // session flags
  //cmdViewSessions.Flags().Bool("mule",false,"List active mule sessions.")
  // Root command
  cmdDeleteMS.Flags().String("id","dfgl0i9u87uty","ID for Mother Ship")
  cmdDeleteMinion.Flags().String("id","dfgl0i9u87uty","ID for minion")
  RootCmd.AddCommand(cmdDeleteMinion)
  RootCmd.AddCommand(cmdDeleteMS)
  RootCmd.AddCommand(cmdLister)
}
/*
var PrintCurrSessions := RunningSessions.NewSession(id, 365 * time.Days, "userData")

var RemoveMule
*/
