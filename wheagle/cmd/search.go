package cmd

import (
  //"time"
  "github.com/alphamystic/odin/lib/utils"

  "github.com/spf13/cobra"
)

var cmdMinionSearch = &cobra.Command {
  Use: "search-minion",
  Short: "Search for a specific minion",
  Run:func(cmd *cobra.Command,args []string){
    id,err := cmd.Flags().GetString("id")
		if err != nil {
			utils.PrintTextInASpecificColor("red","	Error id cannot be empty.")
			return
		}
    if utils.CheckifStringIsEmpty(id){
      RunningSessions.SearchSession(id)
      return
    } else {
      utils.PrintTextInASpecificColor("cyan"," Try search-minion --id [id] or search-admin --id [id]")
      return
    }
  },
}

var cmdAdminSearch = &cobra.Command {
  Use: "search-admin",
  Short: "Search for a specific minion",
  Run:func(cmd *cobra.Command,args []string){
    id,err := cmd.Flags().GetString("id")
		if err != nil {
			utils.PrintTextInASpecificColor("red","	Error id cannot be empty.")
			return
		}
    if utils.CheckifStringIsEmpty(id){
      Conns.SearchConnection(id)
      return
    } else{
      utils.PrintTextInASpecificColor("cyan"," Try search-minion --id [id] or search-admin --id [id]")
      return
    }
  },
}

func init(){
  cmdAdminSearch.Flags().String("id","dfgl0i9u87uty","ID for Mother Ship")
  cmdMinionSearch.Flags().String("id","dfgl0i9u87uty","ID for minion")
  RootCmd.AddCommand(cmdAdminSearch)
  RootCmd.AddCommand(cmdMinionSearch)
}
