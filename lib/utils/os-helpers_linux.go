// +build !windows
package utils

import (
  "syscall"
  "os/exec"
)

var RunExecutable = func(pathToExec string)(int,error){
  cmd := exec.Command(pathToExec)
  err := cmd.Start()
  if err != nil{
    return 0,err
  }
  return cmd.Process.Pid,nil
}

var KillExec = func(pid int)error{
  return syscall.Kill(pid,syscall.SIGTERM)
}
