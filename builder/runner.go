package builder

import (
  "fmt"
  "bytes"
  "errors"
  "strings"
  "text/template"
  "odin/lib/utils"
  "odin/wheagle/server/lib"
)

var temp *template.Template

type Builder struct {
  Name string
  Architecture string
  OsType string
  Format string
  Encoding string
  BuildCommand string
  Dir string
  EntryPoint string
  Template string
}

func CreateBuilder(architecture,ops,format,basic,name,entry string,minion bool)(*Builder,error){
  var b Builder
  b.Name  = name
  if !utils.CheckifStringIsEmpty(entry){
    b.EntryPoint = strings.ToUpper(entry)
  } else {
    b.EntryPoint = "WHEAGLE"
  }
  b.SetOutputDirectory("../bin/temp/")
  err := b.SetOsType(ops)
  if err != nil{
    return nil,err
  }
  err = b.SetArchitecture(architecture)
  if err != nil{
    return nil,err
  }
  err = b.SetFormat(format)
  if err != nil{
    return nil,err
  }
  b.SetTemplate(minion)
  err = b.SetEncoding(basic)
  if err != nil{
    return nil,err
  }
  err = b.CreateBuildCommand()
  if err != nil {
    return nil,err
  }
  return &b,nil
}

func (b *Builder) SetOutputDirectory(dir string){
  b.Dir = dir
}
/*
   PLUGIN Builders
   Windows GOOS=windows GOARCH=amd64 go build -buildmode=c-archive -o tomcat-checker.dll tomcat-checker.go
   Linux go build -buildmode=plugin -o tomcat-checker.so tomcat-checker.go
*/
// if file type does not need building, the required copy or creator function should be called
func (b *Builder) CreateBuildCommand() error{
  switch b.Format {
    case "exe": //go build -buildmode=c-archive for windows plugins
      cmnd := `GOOS={{.OsType}} GOARCH={{.Architecture}}  go build  -ldflags="-s -w" -o {{.Name}}`
      var buf bytes.Buffer
      temp = template.Must(template.New("cmnd").Parse(cmnd))
      _ = temp.Execute(&buf,struct{OsType,Architecture,Name string}{b.OsType,b.Architecture,b.Dir + b.Name + ".exe"})
      cmnd = buf.String()
      b.BuildCommand = cmnd
    case "dll":
      cmnd := `GOOS={{.OsType}} GOARCH={{.Architecture}} CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc go build -ldflags="-s -w -H=windowsgui" -buildmode=c-shared --race  -o {{.Name}}`
      var buf bytes.Buffer
      temp = template.Must(template.New("cmnd").Parse(cmnd))
      _ = temp.Execute(&buf,struct{OsType,Architecture,Name string}{b.OsType,b.Architecture,b.Dir + b.Name + ".dll"})
      cmnd = buf.String()
      b.BuildCommand = cmnd
    case "elf":
      cmnd := `GOOS={{.OsType}} GOARCH={{.Architecture}} go build -ldflags="-s -w" --race  -o {{.Name}}`
      var buf bytes.Buffer
      temp = template.Must(template.New("cmnd").Parse(cmnd))
      _ = temp.Execute(&buf,struct{OsType,Architecture,Name string}{b.OsType,b.Architecture,b.Dir + b.Name})
      cmnd = buf.String()
      b.BuildCommand = cmnd
    case "so":
      cmnd := `GOOS={{.OsType}} GOARCH={{.Architecture}} go build -ldflags="-s -w" -buildmode=plugin --race   -o {{.Name}}`
      var buf bytes.Buffer
      temp = template.Must(template.New("cmnd").Parse(cmnd))
      _ = temp.Execute(&buf,struct{OsType,Architecture,Name string}{b.OsType,b.Architecture,b.Dir + b.Name + ".so"})
      cmnd = buf.String()
      b.BuildCommand = cmnd
    case "iso":
      return errors.New("Generate normal payload  then try [WHEAGLE]: load-poncupine --iF payload.exe --oF payload.iso --f iso")
    case "apk":
      cmnd := `GOOS={{.OsType}} GOARCH={{.Architecture}} CGO_ENABLED=1 go build -ldflags="-s -w" -buildmode=c-shared -o {{.Name}}`
      var buf bytes.Buffer
      temp = template.Must(template.New("cmnd").Parse(cmnd))
      _ = temp.Execute(&buf,struct{OsType,Architecture,Name string}{b.OsType,b.Architecture,b.Dir + b.Name + ".apk"})
      cmnd = buf.String()
      b.BuildCommand = cmnd
    case "bin":
      if b.OsType == "darwin"{
        cmnd := `GOOS={{.OsType}} GOARCH={{.Architecture}} go build  -ldflags="-s -w" --race  -o {{.Name}}`
        var buf bytes.Buffer
        temp = template.Must(template.New("cmnd").Parse(cmnd))
        _ = temp.Execute(&buf,struct{OsType,Architecture,Name string}{b.OsType,b.Architecture,b.Dir + b.Name + ".bin"})
        cmnd = buf.String()
        b.BuildCommand = cmnd
        return nil
      }
      return errors.New("Error: .bin files implementation only works for mac os.")
    case "gplgn":
      cmnd := `GOOS={{.OsType}} GOARCH={{.Architecture}} go build -ldflags="-s -w" -buildmode=plugin --race   -o {{.Name}}`
      var buf bytes.Buffer
      temp = template.Must(template.New("cmnd").Parse(cmnd))
      if b.OsType == "windows"{
        _ = temp.Execute(&buf,struct{OsType,Architecture,Name string}{b.OsType,b.Architecture,b.Dir + b.Name + ".dll"})
      } else {
        utils.PrintTextInASpecificColor("cyan","Assuming plugin for linux/unix style (.so) ")
        _ = temp.Execute(&buf,struct{OsType,Architecture,Name string}{b.OsType,b.Architecture,b.Dir + b.Name + ".so"})
      }
      cmnd = buf.String()
      b.BuildCommand = cmnd
    case "macbin":
      cmnd := `GOOS={{.OsType}} GOARCH={{.Architecture}} CGO_ENABLED=1 go build  -ldflags="-s -w" -buildmode=c-shared -o {{.Name}}`
      var buf bytes.Buffer
      temp = template.Must(template.New("cmnd").Parse(cmnd))
      _ = temp.Execute(&buf,struct{OsType,Architecture,Name string}{b.OsType,b.Architecture,b.Dir + b.Name + ".macbin"})
      cmnd = buf.String()
      b.BuildCommand = cmnd
    default:
      return fmt.Errorf("Unsurpoted format. We are in build this actually shouldn't happen unless it's a file type formmat")
  }
  return nil

  /*var cmdTemp = `GOOS={{.OsType}} GOARCH={{.Architecture}} go build -o {{.Name}}`
  GOOS=windows GOARCH=amd64 CGO_ENABLED=1 go build -buildmode=c-shared -o library.dll library.go
  GOOS=linux GOARCH=arm CGO_ENABLED=1 go build -buildmode=c-shared -o library.so library.go
  GOOS=darwin GOARCH=amd64 CGO_ENABLED=1 go build -buildmode=c-shared -o library.dylib library.go
  GOOS=ios GOARCH=arm64 CGO_ENABLED=1 go build -buildmode=c-shared -o library.so library.go
  GOOS=android GOARCH=arm64 CGO_ENABLED=1 go build -buildmode=c-shared -o library.so library.go
  for iso genisoimage -o /home/sam/payload.iso -r /home/sam/payload/*/
}

func (b *Builder) SetArchitecture(arch string)error{
  if arch == ""{
    return errors.New("architecture can not be nil")
  }
  switch arch {
  case "x32","386","x86","86":
      b.Architecture = "386"
    case "64","x64":
      b.Architecture = "amd64"
    case "arm":
      b.Architecture = "arm"
    case "arm64":
      b.Architecture = "arm64"
    case "wasm":
      b.Architecture = "wasm"
    case "mips64":
      b.Architecture = "mips64"
    case "ppc64":
      b.Architecture = "ppc64"
    case "mips":
      b.Architecture = "mips"
    case "mips64le":
      b.Architecture = "mips64le"
    case "riscv64":
      b.Architecture = "riscv64"
    case "s390x":
      b.Architecture = "s390x"
    default:
      return errors.New("Unsurpported architecture")
  }
  return nil
}

func (b *Builder) SetOsType(osName string) error{
  switch osName {
    case "windows","win":
      b.OsType = "windows"
    case "android","and","andr":
      b.OsType = "android"
    case "apple","appl","ios":
      b.OsType = "ios"
    case "linux","lin","lnx":
      b.OsType = "linux"
    case "aix":
      b.OsType = "aix"
    case "freebsd","fbsd":
      b.OsType = "freebsd"
    case "solaris","sls":
      b.OsType = "solaris"
    case "netbsd","nbsd":
      b.OsType = "netbsd"
    case "openbsd","opbsd":
      b.OsType = "openbsd"
    case "plan9","p9":
      b.OsType = "plan9"
    case "illumos","ilm":
      b.OsType = "illumos"
    case "darwin","dwrn","macos","osx","mac":
      b.OsType = "darwin"
    case "dragonfly","drgf","dgr","drgn":
      b.OsType = "dragonfly"
    default:
      return errors.New("Error, Unsurported operating system")
  }
  return nil
}

func (b *Builder) SetFormat(fmrt string)error{
  switch fmrt {
    case "apk":
      b.Format = "apk"
    case "exe":
      b.Format = "exe"
    case "dll":
      b.Format = "dll"
    case "iso":
      b.Format = "iso"
    case "lin","elf":
      b.Format = "elf"
    case "so":
      b.Format = "so"
    case "mac":
      b.Format = "macbin"
    case "bin":
      b.Format = "bin"//default shellcode generator format
    case "go-plugin":
      b.Format = "gplgn"
    default:
      return errors.New("Unsurpoted payload format or not yet implemented")
  }
  return nil
}

func (b *Builder) SetTemplate(minion bool) {
  switch b.Format {
    case "gplgn":
      if minion {
        b.Template = lib.MinionLib
      } else{
        b.Template = lib.AdminLib
      }
    case "dll":
      if minion {
        b.Template = lib.DLLLoaderMinion
      } else {
        b.Template = lib.DLLLoaderAdmin
      }
    default:
      if minion{
        b.Template = lib.Mule
      } else {
        b.Template = lib.AdminGRPC
      }
  }
}

func(b *Builder)  SetEncoding(enc string) error{
  switch enc {
    case "basic":
      b.Encoding = "basic"
    case "donut":
      if b.OsType != " windows"{
        return errors.New("donut only supports encoding for windows files/executables")
      }
      b.Encoding = "donut"
    default:
      return errors.New("Encoding not surpported.")
  }
  return nil
}

func (b *Builder) Encoder()error{
  return nil
}

var SupportedArch = func(){
  var surportedArch = `GOOS=aix GOARCH=ppc64
  GOOS=android GOARCH=386
  GOOS=android GOARCH=amd64
  GOOS=android GOARCH=arm
  GOOS=android GOARCH=arm64
  GOOS=darwin GOARCH=amd64
  GOOS=darwin GOARCH=arm64
  GOOS=dragonfly GOARCH=amd64
  GOOS=freebsd GOARCH=386
  GOOS=freebsd GOARCH=amd64
  GOOS=freebsd GOARCH=arm
  GOOS=freebsd GOARCH=arm64
  GOOS=illumos GOARCH=amd64
  GOOS=ios GOARCH=amd64
  GOOS=ios GOARCH=arm64
  GOOS=js GOARCH=wasm
  GOOS=linux GOARCH=386
  GOOS=linux GOARCH=amd64
  GOOS=linux GOARCH=arm
  GOOS=linux GOARCH=arm64
  GOOS=linux GOARCH=loong64
  GOOS=linux GOARCH=mips
  GOOS=linux GOARCH=mips64
  GOOS=linux GOARCH=mips64le
  GOOS=linux GOARCH=mipsle
  GOOS=linux GOARCH=ppc64
  GOOS=linux GOARCH=ppc64le
  GOOS=linux GOARCH=riscv64
  GOOS=linux GOARCH=s390x
  GOOS=netbsd GOARCH=386
  GOOS=netbsd GOARCH=amd64
  GOOS=netbsd GOARCH=arm
  GOOS=netbsd GOARCH=arm64
  GOOS=openbsd GOARCH=386
  GOOS=openbsd GOARCH=amd64
  GOOS=openbsd GOARCH=arm
  GOOS=openbsd GOARCH=arm64
  GOOS=openbsd GOARCH=mips64
  GOOS=plan9 GOARCH=386
  GOOS=plan9 GOARCH=amd64
  GOOS=plan9 GOARCH=arm
  GOOS=solaris GOARCH=amd64
  GOOS=windows GOARCH=386
  GOOS=windows GOARCH=amd64
  GOOS=windows GOARCH=arm
  GOOS=windows GOARCH=arm64`
  fmt.Println(surportedArch)
}
//gen Iso
/*
#!/bin/bash

# specify the source directory and the output ISO file
src_dir=./my_files
out_file=my_files.iso

# create the ISO file
genisoimage -o $out_file -r $src_dir

*/
