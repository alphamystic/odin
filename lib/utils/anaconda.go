package utils

import (
  "fmt"
  "net/http"
)
/*
  This is some stupid HTTP file server( will probably add https later on)
  * Use it to download payloads from somewhere/your payloads directory
  *Defaults to ./bin/temp at port 33333
  * Will add upload functionalities later
*/

type Anaconda struct{
  Dir string
  Port int
  Address string
  Run chan bool
}

func (a *Anaconda) AnacondaServe(){
  PrintTextInASpecificColorInBold("cyan","======================================================")
  NoNewLine("cyan","=======    ********************************    =======\n")
  NoNewLine("cyan","=======         ")
  NoNewLine("white"," ANACONDA FILE SERVER  ")
  NoNewLine("cyan","         =======\n")
  NoNewLine("cyan","========        ***************                =======\n")
  PrintTextInASpecificColorInBold("cyan","======================================================")

  runCh := make(chan bool)
  go func(){
    for{
      select{
      case start,ok := <- runCh:
        if !ok{
          PrintTextInASpecificColor("blue","Specify run value.");return
        }
        if start {
          http.Handle("/",http.StripPrefix("/",http.FileServer(http.Dir(a.Dir))))
          PrintTextInASpecificColor("blue",fmt.Sprintf("Anaconda serving files at %d on directory %s",a.Port,a.Dir))
          Logerror(http.ListenAndServe(a.Address + ":" + IntToString(a.Port),nil))
        } else {
          PrintTextInASpecificColor("blue","Clossing anaconda........")
          return
        }
      }
    }
    }()
  for {
    select{
    case start,ok := <- a.Run:
      if !ok{
        PrintTextInASpecificColor("blue","Specify run value.");return
      }
      if start{
        runCh <- true
      } else {
        runCh <- false
      }
    }
  }
}

func (a *Anaconda) Stop(val bool){
  a.Run <- val
}
