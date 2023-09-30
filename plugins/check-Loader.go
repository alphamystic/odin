package plugins


import(
  "os"
  "fmt"
  //"errors"
  "plugin"
  "io/ioutil"
  "github.com/alphamystic/odin/lib/utils"
)

type Checker interface {
  Check(host string,ssl bool,port uint64) *Result
}

type Result struct{
  Vulnerable bool
  Details string
}

func (pl *PluginLoader) LoadChecker(host string,port int)error{
  var (
    files []os.FileInfo
    err error
    p *plugin.Plugin
    n plugin.Symbol
    check Checker
    res *Result
  )
  if files,err = ioutil.ReadDir(PluginsDir); err != nil{
    return fmt.Errorf("Error reading plugins Directory. %v",err)
  }
  for idx := range files {
    if pl.Name == files[idx].Name(){
      if p,err = plugin.Open(PluginsDir+"/"+pl.Name);err != nil{
        return fmt.Errorf("")
      }
      if n,err = p.Lookup("New"); err != nil {
        return fmt.Errorf("No New entry point found. %v",err)
      }
      newFunc,ok := n.(func() Checker)
      if !ok{
        return fmt.Errorf("Plugin entry point is not good. Expecting: func New() plugins.Checker{...}\n ERROR: %v",err)
      }
      check = newFunc()
      res = check.Check(host,false,uint64(port))
      if res.Vulnerable{
        utils.PrintTextInASpecificColor("yellow",fmt.Sprintf("Host is Vulnerable: %s",res.Details))
      } else {
        utils.PrintTextInASpecificColor("white", "Too bad budddy........:) Host not Vulnerable.")
      }
    }
  }
  return nil
}
