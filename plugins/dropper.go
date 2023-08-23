package plugins

import(
  "os"
  "fmt"
  "errors"
  "plugin"
  "io/ioutil"
  "odin/lib/utils"
)

type Dropper interface {
  Drop()(Implementor)
  Wheagle(bool) (func()error)  
}


type Implementor func(string, ...string)error

func (pl *PluginLoader) LoadImplementor()(error,func(string, ...string)error){
  var (
    files []os.FileInfo
    err error
    p *plugin.Plugin
    n plugin.Symbol
    dz Dropper
    implementor Implementor
  )
  if files,err = ioutil.ReadDir(PluginsDir); err != nil{
    return err,implementor
  }
  for idx := range files {
    utils.PrintTextInASpecificColorInBold("cyan",fmt.Sprintf("Found Plugin: "+ files[idx].Name()))
    if files[idx].Name() == pl.Name {
      if p,err = plugin.Open(PluginsDir + "/" + files[idx].Name()); err != nil{
        return err,implementor
      }
      if n,err = p.Lookup("NewDropper"); err != nil{
        return err,implementor
      }
      drpr,ok := n.(func() Dropper)
      if !ok {
        return fmt.Errorf("Invalid plugin entry point is no good."),implementor
      }
      dz = drpr()
      implementor = dz.Drop()
    } else {
      return errors.New(fmt.Sprintf("No plugin with the specified name: %s",pl.Name)),implementor
    }
  }
  return nil,implementor
}

func (pl *PluginLoader) LoadWheagle()func()error{
  return func()error{
    return nil
  }
}
