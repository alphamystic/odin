package ph


import (
  "os"
  "fmt"
  "plugin"
  "io/ioutil"
  "odin/lib/utils"
  "odin/lib/handlers"
)

const PluginsDir = "../bin/"

func LoadRico(rd handlers.ReconData) []handlers.Vulnerabilities{
  var (
    files []os.FileInfo
    err error
    p *plugin.Plugin
    symb plugin.Symbol
    scanner handlers.Scanner
    vulns []handlers.Vulnerabilities
  )
  if files,err = ioutil.ReadDir(PluginsDir);err != nil {
    utils.Logerror(err)
  }
  for idx := range files {
    if p,err = plugin.Open(PluginsDir + "/" + files[idx].Name()); err != nil {
      utils.Logerror(err)
    }
    if symb,err = p.Lookup("New"); err != nil{
      utils.Logerror(err)
    }
    newFunc, ok := symb.(func() handlers.Scanner)
    if !ok {
      utils.Logerror(fmt.Errorf("Private has an invalid  entry point. Expecting: func New() handlers.Scanner{ ... }"))
    }
    scanner = newFunc()
    vulns = scanner.Scan(rd)
  }
  return vulns
}
