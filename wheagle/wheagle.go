package main

import (
  "os"
  "fmt"
  "flag"
  "bufio"
  "strings"
  _ "image/jpeg"
	_ "image/png"
  "github.com/alphamystic/odin/lib/utils"
  "github.com/alphamystic/odin/wheagle/cmd"

  "github.com/qeesung/image2ascii/convert"

)

func printBanner(){
  convertOptions := convert.DefaultOptions
	convertOptions.FixedWidth = 100
	convertOptions.FixedHeight = 50
	converter := convert.NewImageConverter()
	fmt.Print(converter.ImageFile2ASCIIString("./wheagle.png", &convertOptions))
}

func main(){
  var (
    banner = flag.Bool("ban",true,"Print banner while starting")
  )
  flag.Parse()
  if *banner == true {
    printBanner()
  }
  utils.PrintTextInASpecificColorInBold("cyan","**************************************************************************************")
  utils.NoNewLine("blue","    WriterTwitter: ")
  utils.NoNewLine("white"," "+ cmd.WriterTwitter + "\n")
  utils.NoNewLine("blue","    WriterContact: ")
  utils.NoNewLine("white"," "+ cmd.WriterContact)
  fmt.Println("")
  fmt.Println("[WHEAGLE]  Starting commandline")
  // start cli
  reader := bufio.NewReader(os.Stdin)
  for {
    utils.Terminal()
    input,_ := reader.ReadString('\n')
    input = strings.TrimSuffix(input,"\n")
    args := strings.Fields(input)
    if len(args) == 0 {
      continue
    }
    cmd.RootCmd.SetArgs(args)
    err := cmd.RootCmd.Execute()
    if err != nil {
      fmt.Println("Error: ",err)
    }
  }
}

/* Run something async
func main() {
    ticker := time.NewTicker(5 * time.Minute)
    done := make(chan bool)

    go func() {
        for {
            select {
            case <-done:
                return
            case <-ticker.C:
                // Call the function that you want to run asynchronously
                go pingSomething()
            }
        }
    }()

    // Your main code goes here, which can handle requests
    // ...

    // To stop the background function, you can close the "done" channel
    done <- true
}

func pingSomething() {
    // This is the function that you want to run asynchronously
    // ...
}

*/
