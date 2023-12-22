package core

import (
  "fmt"
)

type ShelcodeGen struct {
  InitialData []byte
  CleanedData []byte
  SCFormat ShellCodeFormat
}

type ShellCodeFormatint

const (
  CType ShellCodeFormat = iota
  CSharp
)

// if true, cleaned shellcode is printed if not it prints the none cleaned
// thinking of having this stored incide sg
// might later on add a length checker
func (sg *ShelcodeGen) DisplayShellcodeStrings(scType bool) {
  var shellcode []byte
  // Display the shellcode in a msf-like format
  if scType {
    shellcode = sg.CleanedData
  } else {
    shellcode = sg.InitialData
  }
  for _, b := range shellcode {
    fmt.Printf("\\x%02x", b)
  }
  fmt.Println("\";")
}

func (sg *ShelcodeGen) SanitizeShellcode(shellcode []byte) []byte {
  var sanitizedShellcode []byte
  for _, b := range sg.InitialData {
    if b != 0x00 && b != 0x0a && b != 0x0d && b != 0x20 {
      sanitizedShellcode = append(sanitizedShellcode, b)
    }
  }
  sg.CleanedData = sanitizedShellcode
  return
}
