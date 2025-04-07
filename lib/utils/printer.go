package utils

import (
  "os"
  "fmt"
  "errors"
  "strings"
  "github.com/fatih/color"
)

//log error of errors
func Logerror(e error){
  if e != nil {
    red := color.New(color.FgRed)
    whiteBackground := red.Add(color.BgWhite)
    whiteBackground.Println(e)
  }
  fmt.Println("")
}

func Danger(e error){
  if e != nil{
    red := color.New(color.FgRed)
    boldRed := red.Add(color.Bold)
    boldRed.Println(e)
  }
}

func CustomError(text string,e error){
  if e != nil{
    red := color.New(color.FgRed)
    whiteBackground := red.Add(color.BgWhite)
    whiteBackground.Println(errors.New(fmt.Sprintf("[-]    %s: %s",text,e)))
  }
}

func DangerPanic(e error){
  if e != nil{
    red := color.New(color.FgRed)
    boldRed := red.Add(color.Bold)
    boldRed.Println(e)
    os.Exit(1)
  }
}
func Notice(text string){
  notice := color.New(color.Bold, color.FgBlue).PrintlnFunc()
  notice("[+]   NOTICE:    " + text + "\n")
}

func NoticeError(text string){
  redd := color.New(color.FgRed).PrintfFunc()
  redd("[-]   ERROR:  ")
  redd(" %s\n", text)
}

func Warning(text string){
  redd := color.New(color.FgCyan).PrintfFunc()
  redd("[-] WARNING:   ")
  redd(" %s\n", text +" \n")
}

func Terminal(){
  white := color.New(color.FgWhite).PrintfFunc()
  white("[WHEAGLE] :$  ")
}

func Odin(){
  white := color.New(color.FgWhite).PrintfFunc()
  white("[ODIN] :$  ")
}

func PrintInformation(text string){
  e := color.New(color.FgYellow, color.Bold)
  fmt.Printf("[+]   ")
  e.Printf(text+"\n")
}

func Interactor(text string,admin bool){
  e := color.New(color.FgMagenta,color.Bold)
  if admin {
    e.Printf("INTERACTING WITH ADMIN: \n")
  } else {
    e.Printf("INTERACTING WITH MULE: \n")
  }
  interactor := color.New(color.FgCyan,color.Bold).PrintfFunc()
  interactor("   |--  " + text + ":> ")
}

func PrintTextInASpecificColorInBold(colorName,text string){
  switch strings.ToLower(colorName) {
  case "yellow":
    e := color.New(color.FgYellow, color.Bold)
    fmt.Printf("[+]   ")
    e.Printf(text + "\n")
  case "red":
    e := color.New(color.FgRed, color.Bold)
    fmt.Printf("[+]   ")
    e.Printf(text + "\n")
  case "green":
    e := color.New(color.FgGreen, color.Bold)
    fmt.Printf("[+]   ")
    e.Printf(text + "\n")
  case "magenta":
    e := color.New(color.FgMagenta, color.Bold)
    fmt.Printf("[+]   ")
    e.Printf(text + "\n")
  case "white":
    e := color.New(color.FgWhite, color.Bold)
    fmt.Printf("[+]   ")
    e.Printf(text + "\n")
  case "blue":
    e := color.New(color.FgBlue, color.Bold)
    fmt.Printf("[+]   ")
    e.Printf(text + "\n")
  default:
    e := color.New(color.FgCyan, color.Bold)
    fmt.Printf("[+]   ")
    e.Printf(text + "\n")
  }
}
/*
redd := color.New(color.FgCyan).PrintfFunc()
redd("[-] WARNING:   ")
redd(" %s\n", text +" \n")*/
func PrintTextInASpecificColor(colorName,text string){
  switch strings.ToLower(colorName) {
  case "yellow":
    e := color.New(color.FgYellow)
    fmt.Printf("[+]   ")
    e.Printf(text + "\n")
  case "red":
    e := color.New(color.FgRed)
    fmt.Printf("[+]   ")
    e.Printf(text + "\n")
  case "green":
    e := color.New(color.FgGreen)
    fmt.Printf("[+]   ")
    e.Printf(text + "\n")
  case "magenta":
    e := color.New(color.FgMagenta)
    fmt.Printf("[+]   ")
    e.Printf(text + "\n")
  case "white":
    e := color.New(color.FgWhite)
    fmt.Printf("[+]   ")
    e.Printf(text + "\n")
  case "blue":
    e := color.New(color.FgBlue)
    fmt.Printf("[+]   ")
    e.Printf(text + "\n")
  default:
    e := color.New(color.FgCyan)
    fmt.Printf("[+]   ")
    e.Printf(text + "\n")
  }
}

func NoNewLine(colorName,text string){
  switch strings.ToLower(colorName) {
  case "yellow":
    e := color.New(color.FgYellow, color.Bold)
    e.Printf(text)
  case "red":
    e := color.New(color.FgRed, color.Bold)
    e.Printf(text)
  case "green":
    e := color.New(color.FgGreen, color.Bold)
    e.Printf(text )
  case "magenta":
    e := color.New(color.FgMagenta, color.Bold)
    e.Printf(text)
  case "white":
    e := color.New(color.FgWhite, color.Bold)
    e.Printf(text)
  case "blue":
    e := color.New(color.FgBlue, color.Bold)
    e.Printf(text)
  default:
    e := color.New(color.FgCyan, color.Bold)
    e.Printf(text)
  }
}
