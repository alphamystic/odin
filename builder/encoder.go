package builder

import (
  "os"
  "fmt"
  "encoding/hex"
  "github.com/alphamystic/odin/lib/utils"
)
type Encoder interface {
  Encode([]byte) ([]byte,error)
}


var Encoding = func(inputFile,outputFile string) error{
  input,err := os.ReadFile(inputFile)
  if err != nil {
    return fmt.Errorf("Error reading from input file %q. %v",inputFile,err)
  }
  var tech string
  utils.PrintTextInASpecificColor("cyan","Encoding techniques surported include: \n  xor  \nhex  \ndonut(only  windows as it uses godonut.)")
  fmt.Sprintf("Enter encoding technique. \n")
  utils.NoNewLine("yellow","  [HUNTER]:")
  fmt.Scanln(&tech)
  switch tech {
    case "xor":
      key := byte(0x3F)
      for i := 0; i < len(input); i++{
        input[i] ^= key
      }
      err = os.WriteFile(outputFile,input,0644)
      if err != nil { return fmt.Errorf("Error writing to encoded file. %v",err) }
      return nil
    case "hex":
      encd := hex.EncodeToString(input)
      if err := os.WriteFile(outputFile,[]byte(encd),0644); err != nil {
        return fmt.Errorf("Error writing to output file.\n ERROR: %s",err)
      }
      return nil
    case "doonut","dnt":
      var arch,modType string
      fmt.Sprintf("Enter architecture of payload. (x64 or x32 or x84)\n")
      utils.NoNewLine("yellow","  [HUNTER]:")
      fmt.Scanln(&arch)
      fmt.Sprintf("Enter donut module type. ( Any of this: dll,exe,un_dll,un_exe,xsl,js,vbs )\n")
      utils.NoNewLine("yellow","  [HUNTER]:")
      fmt.Scanln(&modType)
      var e Encoder
      e = &Donut {arch,modType}
      data,err := e.Encode(input)
      if err != nil {
        return fmt.Errorf("Error encoding with donut.\n %s",err)
      }
      if err := os.WriteFile(outputFile,data,0644); err != nil{
        return fmt.Errorf("Error writing output to file. \nERROR: %s",err)
      }
      return nil
    default:
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

/*
func BadCharRemover(){
  sub := strings.NewReplacer(
    "0x00","0x90",
    "0xfr","0x50",
  )
  \xff\xe4
  0x0D
  \x40
  \x0a

}
*/
