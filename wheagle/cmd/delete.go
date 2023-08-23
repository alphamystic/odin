package cmd

import (
  "odin/lib/utils"

  "github.com/spf13/cobra"
)

var cmdMinionDelete = &cobra.Command {
  Use: "delete-minion",
  Short: "Search for a specific minion",
  Run:func(cmd *cobra.Command,args []string){
    id,err := cmd.Flags().GetString("id")
		if err != nil {
			utils.PrintTextInASpecificColor("red","	Error id cannot be empty.")
			return
		}
    if utils.CheckifStringIsEmpty(id){
      if err := RunningSessions.DeleteSession(id,SessionsDriver); err != nil{
        utils.Logerror(err)
      }
      return
    } else {
      utils.PrintTextInASpecificColor("cyan"," Try search-minion --id [id] or search-admin --id [id]")
      return
    }
  },
}

var cmdAdminDelete = &cobra.Command {
  Use: "delete-admin",
  Short: "Search for a specific minion",
  Run:func(cmd *cobra.Command,args []string){
    id,err := cmd.Flags().GetString("id")
		if err != nil {
			utils.PrintTextInASpecificColor("red","	Error id cannot be empty.")
			return
		}
    if utils.CheckifStringIsEmpty(id){
      if err := Conns.RemoveConnection(id,ConnectorsDriver); err != nil{
        utils.Logerror(err)
      }
      return
    } else{
      utils.PrintTextInASpecificColor("cyan"," Try search-minion --id [id] or search-admin --id [id]")
      return
    }
  },
}

func init(){
  cmdAdminDelete.Flags().String("id","dfgl0i9u87uty","ID for Mother Ship")
  cmdMinionDelete.Flags().String("id","dfgl0i9u87uty","ID for minion")
  RootCmd.AddCommand(cmdAdminDelete)
  RootCmd.AddCommand(cmdMinionDelete)
}
