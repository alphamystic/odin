package main

import (
  "os"
  "fmt"
  "bufio"
  "strings"
  "odin/cli"
  "odin/lib/utils"

  "github.com/common-nighthawk/go-figure"
)

func main(){
  myFigure := figure.NewFigure("Odin", "isometric1", true)
  //myFigure := figure.NewFigure("Odin", "basic", true).Scroll(10000, 200, "right")
  myFigure.Print()
  utils.PrintTextInASpecificColorInBold("white","Initializing ODIN.....")
  fmt.Println("[ODIN]  Starting commandline")
  // start cli
  reader := bufio.NewReader(os.Stdin)
  for {
    utils.Odin()
    input,_ := reader.ReadString('\n')
    input = strings.TrimSuffix(input,"\n")
    args := strings.Fields(input)
    if len(args) == 0 {
      continue
    }
    cli.RootCmd.SetArgs(args)
    err := cli.RootCmd.Execute()
    if err != nil {
      fmt.Println("Error: ",err)
    }
  }
}
