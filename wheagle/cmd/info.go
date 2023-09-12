package cmd

import (
  "fmt"
  "odin/lib/utils"

  "github.com/spf13/cobra"
)

const(
  WriterName = "Odhiambo Samuel"
  WriterTwitter = "https://twitter.com/3lOr4cle"
  WriterLinkedin = "https://linkedin/mylinkedin"
  WriterContact ="+254 717408948 or +254 717845889"
  WriterEmail = "sodhis001@gmail.com or eloracle.odhiambosamuel@gmail.com"
  Github = "github.com/alphamystic/"
  WheagleGithub = "github.com/alphamystic/odin/wheagle"
  Basic = `
  Wheagle is a simple C2 framework to easen up pentesting by making implant generation and C2 control easier.
  This is so by tying up each implant to a perticular C2 and avoiding regeneration/starting of a listener every time.
  For a C2 hosted on another platform, one can still interact with it as though it's an implant on it's own to kill it manage it or as a backdoor.
  Wheagle comes with a ransomware simulator that basically creates a folder on your desktop then encrypts it and adds a readme.wheagle with the decryption key.
  Currently we don't support multiple operators or tunneling as we are still in the build stage but future implementations might include them.
  Should be easy to detect as syscalls aren't hashed and most techniques for persistence are noisy like (cronjobs forlinux or registryand service for windows.)
  To get started on basic commands try [WHEAGLE]: basics
  To avoid detection, use droppers like .iso/.zip with your generated payload/c2 to mask it from a few AV's or use a generated library (.so/.dll) to load it up
  Due to some unavoidable circumstances, golang has an issue with plugins so that feature probably holds on for a while but take a sneak peak at
  he plugins directory at "./plugins to see what we might support and add yours if interested.
  Wheagle is part of Odin, my attempt to create a vulnerability scanner, the whole idea was to automate pentesting process to complete compromise.
  I'm still working on the scanner, it can only do recon as per now. Wheagle was created to actually make it easier to manage completely compromised shells as it willsave
  the session for each and specify a mothership for it.

  Let's see how far I'll get with this. Adios...........
  `
)

const DISCLAIMER = `
WHEAGLE SHOULD ONLY BE USED FOR TRAINING PURPORSES OR ON SYSTEMS IN WHICH YOU HAVE BEEN GIVEN PERMISION TO PENTEST.
WRITERS DO NOT ASSUME LIABLITY/CONSEQUENSES FOR IT'S MISUSE. (WITH GREAT POWER COMES GREAT RESPONCIBILITY.)
JUST BE A GOOD ADULT AND USE IT CORRECTLY. WOULDN'T REALLY HURT.
SOME OF THE PAYLOAD/C2 OPTIONS CAN ACT AS GREAT BACKDOORS, DON'T USE THEM FOR SUCH AS TEMPTING AS IT SOUNDS.
`

var (
  VERSION = "0.0.1"
)

var WheagleCommands = []string{
  `gen-admin --name home --iprotocol tcp --oprotocol tcp --oport 45566  --iport  44566 --arch x64 --os lin --f elf --addr 0.0.0.0 --pass Qwerty`,
  `generate --name home-minon --f elf --ops lin --lhost 127.0.0.1 --lport 44566  --arch x64  --msid 9b8669f83e7769473ed4dd4fc0875d27`,
  `generate --name live --f exe --ops lin --lhost 127.0.0.1 --lport 44566  --arch x64  --msid 5323a0ea5ae00a123ca42220ad255da7`,
  `Interact with admin ia --id [connector id]`,
  `Interact with minion ia --id [session id]`,
  `To view your C2's list-admins`,
  `To view your implants list-minions`,
  `To search for a paerticular minion search-minion [char] where char is the name of the  minion or possible strings used in it's name`,
  `To search for a paerticular admin search-admin [char] wherre char is the name of the admin or possible strings used in it's name.`,
  `To send a command to all miinions use al --cmnd kill (This just kills eveeryone of them)`,
  `Delete a minion delete-minion --id [minion id]`,
  `Delet an admin delete-admin --id [admin id]`,
  `Generate shellcode hunter --f sc --iF payload.exe --oF shellcode.bin`,
  `To create a dropper hunter --f drp --iF payload.exe --oF shellcode.bin`,
  `For system commands system, starts an interactive shell for your particular os`,
  `To start msf, type msf, his startsban msf shell ut without sudo perms, Future implementations will include easier session manager`,
  `MSF create payload to creat a payload with msf`,
  `To run a ransomware simulation ransom --key "WHEAGLE" (Only works when interacting with a target or C2)`,
  `To start or stop a file server at temp directory(Where payloads are stored) use anaconda-stop or anaconda-start`,
  `start-ms Start a c2 in the background. start-ms --name shark --d  edfgvhjkiuyt path ../bin/temp/shark`,
  `kill-ms Kill a running mothership. kill-ms --id dfghjkjhfghj`,
  `     The commands (kill-ms and start-ms):`,
  `           1. Do not support C2 running in whichever the persistance mode but only executables`,
  `           2. Work only if you haven't exited the current session else on close they are removed from memory.(No db storage but will still keep on running)`,
}

var cmdDisclaimer = &cobra.Command{
  Use: "disclaimer",
  Short: "Display the DISCLAIMER",
  Run: func(cmd *cobra.Command,args []string){
    utils.PrintTextInASpecificColorInBold("red",DISCLAIMER)
  },
}

var cmdVersion = &cobra.Command{
  Use: "version",
  Short: "Display the Version number",
  Run: func(cmd *cobra.Command,args []string){
    utils.PrintTextInASpecificColor("white","Your current running version of wheagle is "+ VERSION + ".")
  },
}

var cmdCommands = &cobra.Command{
  Use: "basics",
  Short: "Display wheagle possible/usecase commands",
  Run: func(cmd *cobra.Command,args []string){
    utils.PrintTextInASpecificColorInBold("magenta","***********************************************************************************************************************************************")
    fmt.Printf("        **************************")
    utils.NoNewLine("white","   POSSIBLE WHEAGLE COMANDS   ")
    fmt.Printf("        **************************" +"\n")
    for _, psbCmd := range WheagleCommands{
      utils.PrintTextInASpecificColorInBold("cyan",psbCmd)
    }
    utils.PrintTextInASpecificColorInBold("magenta","***********************************************************************************************************************************************")
  },
}

var cmdInfo = &cobra.Command{
  Use: "info",
  Short: "Display info about wheagle and it's creator",
  Run: func(cmd *cobra.Command,args []string){
    utils.PrintTextInASpecificColorInBold("magenta","******************************************************************************")
    fmt.Printf("        **************************")
    utils.NoNewLine("white","   WHEAGLE INFORMATION   ")
    fmt.Printf("        **************************"+"\n")
    utils.PrintTextInASpecificColorInBold("cyan","*****************************************************************************")
    utils.NoNewLine("blue","    WritersName: ")
    utils.NoNewLine("white"," "+ WriterName)
    fmt.Println("")
    utils.PrintTextInASpecificColorInBold("cyan","****************************************************************************")
    utils.NoNewLine("blue","    Writers Twitter: ")
    utils.NoNewLine("white"," "+ WriterTwitter)
    fmt.Println("")
    utils.PrintTextInASpecificColorInBold("cyan","****************************************************************************")
    utils.NoNewLine("blue","    WriterLinkedin: ")
    utils.NoNewLine("white"," "+ WriterLinkedin)
    fmt.Println("")
    utils.PrintTextInASpecificColorInBold("cyan","****************************************************************************")
    utils.NoNewLine("blue","    WriterContact: ")
    utils.NoNewLine("white"," "+ WriterContact)
    fmt.Println("")
    utils.PrintTextInASpecificColorInBold("cyan","****************************************************************************")
    utils.NoNewLine("blue","    WriterEmail: ")
    utils.NoNewLine("white"," "+ WriterEmail)
    fmt.Println("")
    utils.PrintTextInASpecificColorInBold("cyan","****************************************************************************")
    utils.NoNewLine("blue","    Github: ")
    utils.NoNewLine("white"," "+ Github)
    fmt.Println("")
    utils.PrintTextInASpecificColorInBold("cyan","****************************************************************************")
    utils.NoNewLine("blue","    WheagleGithub: ")
    utils.NoNewLine("white"," "+ WriterTwitter)
    fmt.Println("")
    utils.PrintTextInASpecificColorInBold("magenta","*********************************************************************************")
    utils.PrintTextInASpecificColorInBold("cyan",Basic)
  },
}

func init(){
  RootCmd.AddCommand(cmdDisclaimer)
  RootCmd.AddCommand(cmdVersion)
  RootCmd.AddCommand(cmdCommands)
  RootCmd.AddCommand(cmdInfo)
}
