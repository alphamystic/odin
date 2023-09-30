package main

import (
  "os"
  "fmt"
  "bufio"
  "strings"
  _ "image/jpeg"
	_ "image/png"
  "github.com/alphamystic/odin/plugins"
  "github.com/alphamystic/odin/lib/utils"
  "github.com/alphamystic/odin/wheagle/cmd"

  "github.com/qeesung/image2ascii/convert"

)

type WheagleCNC struct{}

func NewDropper() plugins.Dropper{
  return new(WheagleCNC)
}

func (wc *WheagleCNC) Wheagle(wb bool)(func()error){
  return func()error{
    START:
    if !wb {
      goto END
    }
    convertOptions := convert.DefaultOptions
  	convertOptions.FixedWidth = 100
  	convertOptions.FixedHeight = 50

  	// Create the image converter
  	converter := convert.NewImageConverter()
  	fmt.Print(converter.ImageFile2ASCIIString("./wheagle.png", &convertOptions))
    END:
    fmt.Println("[WHEAGLE]  Starting WHEAGLE commandline")
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
        return err
      }
    }
  }
}
