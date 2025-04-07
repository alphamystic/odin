package main


import (
  "os"
  "fmt"
  "bufio"
  "strings"
  "github.com/alphamystic/odin/cli"
  "github.com/alphamystic/odin/loki"
  "github.com/alphamystic/odin/loki/ui/router"
  "github.com/alphamystic/odin/lib/utils"

  "github.com/common-nighthawk/go-figure"
)

func main(){
  myFigure := figure.NewFigure("Odin", "isometric1", true)
  // Start the server to write and read to
  Loki := &loki.Loki {
    Address: "0.0.0.0",
    PortS: 3001,
    Port: 4000,
    TlsCert: "",
    TlsKey: "",
    Tls: false,
    ApiKey: "", // servers api keey to chat service at main
  }
  svr,_ := Loki.CreateServer()
  rtr := router.NewRouter(svr,svr)
  go func(){
    rtr.Run(true)
  }()
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
