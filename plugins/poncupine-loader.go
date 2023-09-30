package plugins

import(
  "os"
  "fmt"
  //"errors"
  "plugin"
  "io/ioutil"
  "github.com/alphamystic/odin/lib/utils"
)

type Rodent interface{
  Porcupine(string,string,string) *Spike
}

type Spike struct{
  Err error
}

func (pl *PluginLoader) LoadRodent(iF,oF,format string)error{
  var (
    files []os.FileInfo
    err error
    p *plugin.Plugin
    n plugin.Symbol
    rodent Rodent
    spike *Spike
  )
  if files,err = ioutil.ReadDir(PluginsDir); err != nil{
    return fmt.Errorf("Error reading plugins Directory. %v",err)
  }
  for idx := range files {
    fmt.Println("Found plugins with name: ",files[idx].Name())
    if files[idx].Name() == pl.Name{
      if p,err = plugin.Open(PluginsDir+"/"+pl.Name);err != nil{
        return fmt.Errorf("Error opening plugin. %v",err)
      }
      fmt.Println("Found plugin with name: ",files[idx].Name())
      if n,err = p.Lookup("NewDropper"); err != nil {
        return fmt.Errorf("No NewDropper entry point found. %v",err)
      }
      rodentNew,ok := n.(func() Rodent)
      if !ok{
        return fmt.Errorf("Plugin entry point is not good. Expecting: func NewDropper() plugins.Spike{...}\n ERROR: %v",err)
      }
      rodent = rodentNew()
      spike = rodent.Porcupine(iF,oF,format)
      if spike.Err != nil{
        utils.Logerror(spike.Err)
      }
    } else {
      return fmt.Errorf("Found no plugin with name %s",pl.Name)
    }
  }
  return nil
}

//load-poncupine --iF ../bin/temp/test.exe --oF ../bin/test_modified.exe --f dropper
