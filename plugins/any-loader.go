package plugins

import(
  "os"
  "fmt"
  "errors"
  "plugin"
  "io/ioutil"
)

func AnyLoader(name,entryPoint string) error{
  var (
    files []os.FileInfo
    err error
    p *plugin.Plugin
    n plugin.Symbol
  )
  if files,err = ioutil.ReadDir(PluginsDir); err != nil{
    return errors.New(fmt.Sprintf("Error reading plugins Dir: %s \n ERROR: %s"))
  }
  for idx := range files {
    if name == files[idx].Name(){
      if p,err = plugin.Open(PluginsDir +"/" + files[idx].Name()); err != nil{
        return fmt.Errorf("Error opening file.\nERROR: %v",err)
      }
      if n,err = p.Lookup(entryPoint); err != nil{
        return fmt.Errorf("Found no plugin with entry point %s.\n ERROR: %v",entryPoint,err)
      }
      newFunc,ok := n.(func())
      if !ok {
        return errors.New("Plugin entry point is No good.")
      }
      //call the function
      newFunc()
    }
  }
  return nil
}
