package main

import (
  "os"
  "fmt"
  "io"
  "archive/zip"
)
var Zipper = func(inputFile,outputFile string) error {
  input,err := os.Open(inputFile)
  if err != nil {
    return fmt.Errorf("Error reading from input file %q. %v",inputFile,err)
  }
  defer input.Close()
  zipFile,err := os.Create(outputFile)
  if err != nil {
    return fmt.Errorf("Error creating output zip file.\n    ERROR: %q",err)
  }
  defer zipFile.Close()
  writer := zip.NewWriter(zipFile)
  entry,err := writer.Create(input.Name())
  if err != nil {
    return fmt.Errorf("Error creating name incide binary file.\n    ERROR: %q",err)
  }
  _,err = io.Copy(entry,input)
  if err != nil{
    fmt.Errorf("Error copying input into entry.\nERROR: %q",err)
  }
  if err := writer.Close(); err != nil {
    return fmt.Errorf("Error clossing writer.\n    ERROR: %q",err)
  }
  fmt.Println("Wrote to output zip file")
  return nil
}

func main(){
  if err := Zipper("home","/home/sam/Documents/3l0racle/odin/bin/temp/home.zip"); err !=nil{
    panic(err)
  }
  fmt.Println("Done zipping it up...")
}
