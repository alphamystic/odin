package builder

import (
  //"fmt"
  "bytes"

  "odin/vendors/donut"
)


type Donut struct{
  Arch,Type string
}
func(d *Donut)SetArch()donut.DonutArch{
  switch d.Arch {
    case  "x64":
      return donut.X64
    case "x32":
      return donut.X32
    case "x84":
      return donut.X84
    default:
      return donut.X64
  }
}

func (d *Donut) SetModuleType() donut.ModuleType{
  switch d.Type {
    case "dll":
      return donut.DONUT_MODULE_NET_DLL
    case "exe":
      return donut.DONUT_MODULE_NET_EXE
    case "un_dll":
      return donut.DONUT_MODULE_DLL
    case "un_exe":
      return donut.DONUT_MODULE_EXE
    case "xsl":
      return donut.DONUT_MODULE_XSL
    case "js":
      return donut.DONUT_MODULE_JS
    case "vbs":
      return donut.DONUT_MODULE_VBS
    default:
      return donut.DONUT_MODULE_NET_EXE
  }
}
func (d *Donut) Encode(input []byte)([]byte,error){
  arch := d.SetArch()
  osType := d.SetModuleType()
  config := &donut.DonutConfig{
    Arch: arch,
    Type: osType,
    InstType: donut.DONUT_INSTANCE_PIC,
    Entropy: donut.DONUT_ENTROPY_DEFAULT,
    Compress: 1,
		Format:   1,
		Bypass:   3,
  }
  sc,err := donut.ShellcodeFromBytes(bytes.NewBuffer(input),config)
  if err != nil{
    return nil,err
  }
  return sc.Bytes(),nil
}
