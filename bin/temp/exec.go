package main

import "fmt"
import "syscall"

func main(){
err := syscall.Exec("whoami")
if err != nil{
fmt.Sprintf("%s",err.Error())
}
}
