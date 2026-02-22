package main


import (
  "os"
  "fmt"
  "bufio"
  "strings"
  "github.com/alphamystic/odin/cli"
  "github.com/alphamystic/odin/lib/utils"

  "github.com/common-nighthawk/go-figure"
)

func main(){
  myFigure := figure.NewFigure("Odin", "isometric1", true)
  myFigure.Print()
  utils.PrintTextInASpecificColorInBold("white","Initializing ODIN.....")
  fmt.Println("[ODIN]  Starting commandline")
  // Authenticate the USer first to the API
  reader := bufio.NewReader(os.Stdin)

	// fmt.Print("[+] Enter API base URL: ")
	// baseURL, _ := reader.ReadString('\n')
	// baseURL = strings.TrimSpace(baseURL)
  //
	// // Initialize client with base URL: https://localhost:8080
	globalClient := utils.GetClient("portfolio")
  //
	// // Prompt for login credentials
	// fmt.Print("[+] Enter email: ")
	// email, _ := reader.ReadString('\n')
	// email = strings.TrimSpace(email)
  //
	// fmt.Print("[+] Enter password: ")
	// password, _ := reader.ReadString('\n')
	//password = strings.TrimSpace(password)
  //end of login

	// Attempt login
  email := "sam@mail.com"
  password := "12345"
	err := globalClient.Login(email, password)
	if err != nil {
		fmt.Println("Login failed:", err)
		os.Exit(1) // Exit if login fails
	}

	utils.Notice("Login successful!")
	cli.GlobalClient =  globalClient
  // start cli
  //reader = bufio.NewReader(os.Stdin)
  for {
    utils.Odin()
    input,_ := reader.ReadString('\n')
    input = strings.TrimSuffix(input,"\n")
    args := strings.Fields(input)
    if len(args) == 0 {
      continue
    }
    cli.RootCmd.SetArgs(args)
    err := cli.RootCmd.Execute()
    if err != nil {
      fmt.Println("Error: ",err)
    }
  }
}
