/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cli

import (
	"os"
	"fmt"
	"time"
	"io/ioutil"

	"odin/plugins"
	"odin/lib/db"
	"odin/lib/utils"
	"github.com/spf13/cobra"
	"github.com/cheggaaa/pb/v3"
)



// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use: "odin",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

/*
var cmdStartSSHServer
var cmdStartFtpServer
var cmdStartServer


var cmdStatusTarget
var cmdListTargets (-bb, -pentest,-blackops)
var cmdExploitTarget (should either pop a wheagle shell)
var cmdDoRecon

zeroDay
var cmdListZeroDays
var cmdViewZeroDay

user
var cmdCreateUser
var cmdViewUser

MS WHEAGLE
var cmdCreateMothership
var cmdlistMotherships
var cmdViewMotherships
var cmdStatusMothership

issue
var cmdViewIssue
var cmdListIssues (all-mine, -all, -complete,-incomplate)

appointment
var cmdViewAppointment
var cmdListAppointments (-all-mine,-all, -complete,-incomplate)

contacts
var cmdListContacts (-all-mine,-all, -complete,-incomplate)
var cmdViewContacts

bd Backdoors
var cmdViewBackdoors (Not nesceccarily  a backdoor but supported formats present)

project
var cmdViewProject
var cmdListProjects (-all,-mine, -active, -archived)
*/

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the RootCmd.
/*
*/
var cmdPluginTester = &cobra.Command{
	Use: "load",
	Short: "Mannually load a perticular plugin.",
	Long: "If you've made a plugin and want to test it or just load it directly. Load it with the name and entry point",
	Run: func(command *cobra.Command,args []string){
		pName,err := command.Flags().GetString("name")
		if err != nil{
			utils.NoticeError(fmt.Sprintf("%s",err))
			return
		}
		entry,err := command.Flags().GetString("entry")
		if err != nil{
			utils.NoticeError(fmt.Sprintf("%s",err))
			return
		}
		pType,err := command.Flags().GetString("pType")
		if err != nil{
			utils.NoticeError(fmt.Sprintf("%s",err))
			return
		}
		if !utils.CheckifStringIsEmpty(pName) {
			utils.Warning("Plugin name can not be empty.")
			return
		}
		if !utils.CheckifStringIsEmpty(entry) {
			utils.Warning("Entry point can not be empty.")
			return
		}
		utils.PrintTextInASpecificColorInBold("white",fmt.Sprintf("Loading %s with entry %s of type %s",pName,entry,pType))
		if err = plugins.AnyLoader(pName,entry); err != nil{
			utils.Logerror(err)
		}
	},
}
const PluginsDir = "../bin/plugin/"
var Sanitizer = func(name string)(bool,error){
	var (
		compromised bool
		err error
		files []os.FileInfo
	)
	if files,err = ioutil.ReadDir(PluginsDir); err != nil{
		return compromised,fmt.Errorf("Error reading dir %s\nERROR: %v",PluginsDir,err)
	}
	for plgn := range files {
		if files[plgn].Name() == name {
			//get the hash
			hash,err := utils.GetFileHash256(PluginsDir + name)
			if err != nil{
				return compromised,fmt.Errorf("Error hash for plugin: %v",err)
			}
			//check hash
			dbHash,err := GetDBPluginHash(PluginsDir + name)
			if err != nil {
				return compromised,fmt.Errorf("Error getting has from db.\n ERROR: %v",err)
			}
			if hash != dbHash {
				err = utils.DeleteFile(PluginsDir + name)
				if err != nil{
					utils.NoticeError("Error deleting illegal plugin")
					utils.Logerror(err)
				}
			}
			if hash == dbHash {
				compromised = true
			}
		}
	}
	return compromised,nil
}

type LocalPlugins struct {
	Name string
	Hash string
	PType string
	Entry string
	CreatedAt string
	UpdatedAt string
}
func StorePlugin(name,dir,pType,entry string)error{
	hash,err := utils.GetFileHash256(dir+name)
	if err != nil{
		return fmt.Errorf("Error creating hash for file. %v",err)
	}
	plgn := &LocalPlugins{
		Name: name,
		Hash: hash,
		PType: pType,
		Entry: entry,
		CreatedAt: time.Now().String(),
	}
	driver,err := db.Old(PluginsDir,0644)
	if err != nil{
		return err
	}
	if err = driver.Write("plugins",name,plgn); err != nil{
		return err
	}
	if err = utils.CopyFileToDirectory(dir + name,PluginsDir + name); err != nil{
		utils.Warning(fmt.Sprintf("Plugin %s saved to db but not copied to plugins directory %s",name,PluginsDir))
		return fmt.Errorf("Error copying plugins to plugins dir.\nERROR: %v",err)
	}
	utils.PrintTextInASpecificColorInBold("white",fmt.Sprintf("Plugin with name %s stored to db.",name))
	return nil
}

var GetDBPluginHash = func(name string)(string,error){
	driver,err := db.Old(PluginsDir,0644)
	if err != nil{
		return "",fmt.Errorf("Error creating a read driver. %v",err)
	}
	var plgn = new(LocalPlugins)
	err = driver.Read("plugins",name,plgn)
	if err != nil{
		return "",fmt.Errorf("Error reading %s.\nERROR: %v",name,err)
	}
	return plgn.Hash,nil
}

var cmdInstallPlugin = &cobra.Command{
	Use: "install",
	Short: "Install an already downloaded plugin",
	Run:func(cmd *cobra.Command, args []string){
		name,err := cmd.Flags().GetString("name")
		if err != nil{
			utils.NoticeError(fmt.Sprintf("Error in your plugin name.\nERROR:%s",err))
			return
		}
		dir,err := cmd.Flags().GetString("dir")
		if err != nil{
			utils.NoticeError(fmt.Sprintf("Error in your directory.\nERROR:%s",err))
			return
		}
		pType,err := cmd.Flags().GetString("pType")
		if err != nil{
			utils.NoticeError(fmt.Sprintf("Error in your plugin Type.\nERROR:%s",err))
			return
		}
		entry,err := cmd.Flags().GetString("entry")
		if err != nil{
			utils.NoticeError(fmt.Sprintf("Error in your entry point.%s",err))
			return
		}
		if err = StorePlugin(name,dir,pType,entry); err != nil{
			utils.Logerror(err)
		}
	},
}

var cmdOdinLister = &cobra.Command{
	Use: "list",
	Short: "List all your installed plugins",
	Run: func (cmd *cobra.Command, args []string){
		plgns,_ := cmd.Flags().GetString("name")
		scans,_ := cmd.Flags().GetString("name")
		if len(plgns) > 0{
			if !utils.CheckifStringIsEmpty(plgns) {
				utils.Notice("Plugins can't be empty");return
			}
			/*swith between name,all
			if plgns == "name"{
				//scan for the name
				var pName,entryoint string
				fmt.Printf("Enter plugin name to load: ")
				fmt.Scan(&pName)
				fmt.Printf("Enter entry point: ")
				fmt.Scan(&entryPoint)
				err := AnyLoader(pName)
			}*/
			if plgns == "all" {
				driver,err := db.Old(PluginsDir,0644)
				if err != nil{
					utils.Logerror(err)
					return
				}
				plgs,err := driver.ReadAll("plugins")
				if err != nil{
					utils.Logerror(err)
					return
				}
				fmt.Println(plgs)
			}
		}

		if len(scans) > 0 {
			//listing all scans or
			if scans == "all"{
				fmt.Println("Lising scans .............................................")
			}
			//scan for something else or assume it's a name
		}
	},
}

// add a download and install plugin
var cmdStartServer = &cobra.Command{
	Use: "startserver",
	Short: "Start http server",
	Run: func(command *cobra.Command,args []string){
		//utils.PrintTextInASpecificColorInBold("yellow","Quitting ODIN")
		//starting server
	},
}

var cmdQuit = &cobra.Command{
	Use: "quit",
	Short: "Quit the application and terminate all listeners/server",
	Run: func(command *cobra.Command,args []string){
		utils.PrintTextInASpecificColorInBold("yellow","Quitting ODIN")
		os.Exit(0)
	},
}

func init(){
	cmdPluginTester.Flags().String("name","name","Name of Plugin to load. (Add with extenstion say poncurpine.dll or or poncurpine.so)")
	cmdPluginTester.Flags().String("pType","pType","Type of plugin to load")
	cmdPluginTester.Flags().String("entry","entry","Entry point for your Plugin. ")
	cmdInstallPlugin.Flags().String("name","name","Name of Plugin to load. (Add with extenstion say poncurpine.dll or or poncurpine.so)")
	cmdInstallPlugin.Flags().String("pType","pType","Type of plugin you're installing")
	cmdInstallPlugin.Flags().String("entry","entry","Entry point for your Plugin. ")
	cmdInstallPlugin.Flags().String("dir","dir","Directory of Plugin. (`/home/user/Downloads/poncurpine` or `C://User/Downloads/poncurpine.dll`) ")
	RootCmd.AddCommand(cmdInstallPlugin)
	RootCmd.AddCommand(cmdPluginTester)
	RootCmd.AddCommand(cmdOdinLister)
	RootCmd.AddCommand(cmdQuit)
}
func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.odin.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	fmt.Println("[+]	Initalizing setup")
	count := 1000
	bar := pb.StartNew(count)
	for i := 0; i < count; i++ {
		bar.Increment()
		time.Sleep(time.Millisecond)
	}
	// finish bar
	bar.Finish()
	RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
