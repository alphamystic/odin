package cmd

import(
  "fmt"
  "os"
	"os/exec"
  "odin/lib/utils"
  "github.com/spf13/cobra"
)

var Conda = new(utils.Anaconda)
var currentOs = utils.GetCurrentOS()

var cmdStartFileServer =  &cobra.Command{
	Use: "anaconda-start",
	Short: "Run anaconda file server",
	Long: "",
	Run: func(cmd *cobra.Command, args []string){
		go Conda.Stop(true)
		Conda.AnacondaServe()
	},
}

var cmdStopFileServer =  &cobra.Command{
	Use: "anaconda-stop",
	Short: "Interact with a perticular mule/minion or a c2",
	Long: "",
	Run: func(cmd *cobra.Command, args []string){
		Conda.Stop(false)
	},
}

var systemCommand = &cobra.Command{
	Use: "system",
	Short: "Execute a system command while incide wheagle to cli",
	Long: "Execute a system command while incide wheagle to cli. Though I'm not sure if I can pipe it out.",
	//Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command,args []string){
		cmdType,_ := cmd.Flags().GetString("cmdType")
		switch currentOs {
			case "windows":
				if cmdType == "psh" {
					cbr := exec.Command("powershell.exe")
					cbr.Stdin = os.Stdin
					cbr.Stdout = os.Stdout
					cbr.Stderr = os.Stderr
					err := cbr.Run()
					if err != nil {
						utils.Logerror(fmt.Errorf("Error executing powershell: %v",err))
					}
				} else {
					cbr := exec.Command("cmd.exe")
					cbr.Stdin = os.Stdin
					cbr.Stdout = os.Stdout
					cbr.Stderr = os.Stderr
					err := cbr.Run()
					if err != nil {
						utils.Logerror(fmt.Errorf("Error executing cmd: %v",err))
					}
				}
		default:
			cbr := exec.Command("/bin/sh")
			cbr.Stdin = os.Stdin
			cbr.Stdout = os.Stdout
			cbr.Stderr = os.Stderr
			err := cbr.Run()
			if err != nil {
				utils.Logerror(fmt.Errorf("Error executing commandline: %v",err))
			}
		}
	},
}

var cmdStartMsf = &cobra.Command{
	Use: "msf",
	Short: "Start a metasploit commandline. (Probably doesn't work on windows)",
	Long: `Starts a metasploit interactive command line. Works well when running with sudo permisions
				Use cntrl + z to suspend the console and navigate backwith system then fg the normal way.`,
	Run: func(cmd *cobra.Command,args []string){
		if os.Geteuid() != 0 {
			utils.Notice("Not running wheagle as root, restart as admin to run msf properly.")
			utils.Notice("Proceeding to start msf with non sudo privileges")
		}
		cbr := exec.Command("msfconsole")
		cbr.Stdin = os.Stdin
		cbr.Stdout = os.Stdout
		cbr.Stderr = os.Stderr
		err := cbr.Run()
		if err != nil {
			utils.Logerror(fmt.Errorf("Error starting metasploit: %v",err))
		}
	},
}

var cmdHelp = &cobra.Command{
	Use:   "help",
	Short: "Help usage for wheagle C2",
	Long: `Wheagle is a golang based C2 to manage compromised targets.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
	//Args: cobra.MinimumArgs(1),
	Run: func(command *cobra.Command,args []string){
		for _,arg := range args{
			switch arg {
			case "c2":
				utils.PrintTextInASpecificColorInBold("magenta", "C2 is already running, quit to terminate")
			default:
				fmt.Println("[+]		Next it will be the help message on: "+args[0])
			}
		}
	},
}

var cmdQuit = &cobra.Command{
	Use: "quit",
	Short: "Quit the application and terminate all listeners/server",
	Run: func(command *cobra.Command,args []string){
		utils.PrintTextInASpecificColorInBold("yellow","Quitting wheagle")
		os.Exit(0)
	},
}
