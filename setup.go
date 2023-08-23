package main

/*
  * This file is for future Implementations when doe writing tests and making it workable also for windows. and mac
*/
import (
  "os"
  "fmt"
  "log"
  "strings"
  "runtime"
  "os/exec"
  "syscall"
)

type InstallConfig struct{
  IsGolangInstalled bool
  OperatingSystem string
  IsAdmin bool
  InstallDir string
}

func main(){
  fmt.Println("[+]  ODIN Instalation setup......")
  var err error
  // ensure golang is setup
  var config = new(InstallConfig)
  err = config.CheckGoInstallation()
  if err != nil{
    log.Fatal(err)
  }
  //get the os
  switch runtime.GOOS {
  case "windows":
    config.OperatingSystem = "windows"
  case "linux":
    config.OperatingSystem = "linux"
  case "netbsd":
    config.OperatingSystem = "netbsd"
  case "android":
    config.OperatingSystem = "android"
    log.Fatal("Download official apk or build one from custom install")
    return
  default:
    log.Fatal("Your operating system isn't surported, try custom building. \n(Odin and wheagle will run on anything that has golang installed but windows.... \nmight break due to directories and shit................)")
  }
  //get the operating system
  // ensure it's running as admin
  if config.OperatingSystem == "windows"{
    //winows check for admin
    if !IsWindowsAdmin(){
      log.Fatal("Not running as Admin. You need administrator permisions to install odin")
      return
    } else {
      config.IsAdmin = true
    }
  } else {
    if os.Getuid() != 0 {
      log.Fatal("Not running as Admin. You need sudo perms to install")
      return
    } else { config.IsAdmin = true }
  }
  //get the instalation directory
  if err = config.CreateInstalationDirectory(); err != nil{
    log.Fatal(err)
  }
  // get the admin operator accout/default configs
  //git clone or download latest zip from a given url
  //unzip
  //move into folder
  // go build odin
  // go build wheagle
  // set the default environment path/variable alias for both odin and wheagle
  fmt.Println("[+]  Successfully installed ODIN and Whaegle........")
  fmt.Println("[+]  Starting ODIN")
}

func (ic *InstallConfig) CheckGoInstallation()error{
  cmd :=  exec.Command("go","version")
  out,err := cmd.CombinedOutput()
  if err != nil {
    return err
  }
  if strings.Contains(out,[]byte("Error")){
    return fmt.Errorf("Error, golang probably not installed or not properly installed")
  }
  ic.IsGolangInstalled = true
  return nil
}

// create installation directory
func (ic *InstallConfig) CreateInstalationDirectory()error{
  var err error
  ic.InstallDir = os.UserHomeDir() + "./odin"
  if err = os.MkdirAll(ic.InstallDir,0750); err != nil && !os.IsExist(err) {
    return fmt.Errorf("Error creating home instalation directory.\n %q",err);
  }
  if err = os.MkdirAll(ic.InstallDir+"/bin",0750); err != nil && !os.IsExist(err) {
    return fmt.Errorf("Error creating bin directory.\n %q",err);
  }
  if err = os.MkdirAll(ic.InstallDir+"/bin/payloads",0750); err != nil && !os.IsExist(err) {
    return fmt.Errorf("Error creating payloads directory.\n %q",err);
  }
  if err = os.MkdirAll(ic.InstallDir+"/bin/temp",0750); err != nil && !os.IsExist(err) {
    return fmt.Errorf("Error creatig temp directory.\n %q",err);
  }
  if err = os.MkdirAll(ic.InstallDir+"/bin/screenshots",0750); err != nil && !os.IsExist(err) {
    return fmt.Errorf("Error creating screenshots directory.\n %q",err);
  }
  if err = os.MkdirAll(ic.InstallDir+"/bin/downloads",0750); err != nil && !os.IsExist(err) {
    return fmt.Errorf("Error creating downloads directory.\n %q",err);
  }
  if err = os.MkdirAll(ic.InstallDir+"/bin/plugin",0750); err != nil && !os.IsExist(err) {
    return fmt.Errorf("Error creating pluugin files directory.\n %q",err);
  }
  return nil
}

// isAdmin checks if the current user has administrator privileges
var IsWindowsAdmin = func()bool {
	// Get the current process's token
	var token syscall.Token
	err := syscall.OpenProcessToken(syscall.GetCurrentProcess(), TOKEN_QUERY, &token)
	if err != nil {
		return false
	}
	defer token.Close()

	// Look up the LUID for the SeDebugPrivilege privilege
	var luid uint64
	err = procLookupPrivilegeValue.Find()(0, syscall.StringToUTF16Ptr("SeDebugPrivilege"), (*uint64)(unsafe.Pointer(&luid)))
	if err != nil {
		return false
	}
	// Check if the token has the SeDebugPrivilege privilege
	var privs TOKEN_PRIVILEGES
	privs.PrivilegeCount = 1
	privs.Privileges[0].Luid = luid
	var size uint32
	err = procAdjustTokenPrivileges.Find()(token, false, &privs, uint32(unsafe.Sizeof(privs)), nil, &size)
	if err != nil {
		return false
	}

	// Administrator privileges are granted to users with the SeDebugPrivilege privilege
	return privs.Privileges[0].Attributes&SE_PRIVILEGE_ENABLED != 0
}
