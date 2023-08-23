package plugins
/*
import(
  "os"
  "fmt"
  "errors"
  "plugin"
  "io/ioutil"
  "odin/lib/utils"
)

type Backdoor interface {
  BD(targetFile,bdFile string) error
}

//loads the bd plugin and do the vodoo thingy
func (pl *PluginLoader) LoadSkunkBD(targetFile,bdFile string)error{
  var(
    files []os.FileInfo
    err error
    p *plugin.Plugin
    n plugin.Symbol
    skunk Backdoor
  )
  if files,err = ioutil.ReadDir(PluginsDir); err != nil{
    return errors.New(fmt.Sprintf("Error reading plugins Dir: %s \n ERROR: %s"))
  }
  for idx := range files {
    utils.PrintTextInASpecificColorInBold("cyan",fmt.Sprintf("Found Plugin: "+ files[idx].Name()))
    if p,err = plugin.Open(PluginsDir + "/" + files[idx].Name());err != nil{
      return err
    }
    if n,err = p.Lookup("BDSkunk");err != nil{
      return err
    }
    SkunkBD,ok := n.(func()error)
    if !ok {
      return errors.New("Plugin Entry point is no good. Expecting func BDSkunk()error")
    }
    skunk = SkunkBD()
    err = skunk.BD(targetFile,bdFile)
    if err != nil{
      utils.Logerror(err)
    }
    utils.PrintTextInASpecificColorInBold("yellow","Successfuly created a backdoor of file")
  }
  return nil
}
*/
