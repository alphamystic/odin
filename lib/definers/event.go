package definers
/*
  * This package entails a basic event description, However defining all sysmon x auditd logs can be found in their respective files at windows_event.go and linux_event.go
*/

import (
  "time"
  "github.com/alphamystic/odin/lib/utils"
)

type Event struct {
  EId string // a hash of the reported event
  OS OperatingSystem
  Handled bool
  Level EventLevel
  Data string // should be marshaled up into an NTEvent or a LinuxEvent/Log
  utils.TimeStamps
}

type OperatingSystem int

const (
    UnknownOS OperatingSystem = iota
    Windows
    Linux
    MacOS
    iOS
    Android
    OtherOS
)

type EventLevel int
const (
    Info EventLevel = iota
    Low
    Medium
    High
    Critical
)

type NTEvent struct {
  ID int `json: "id"`
  Hash string `json:"hash"`
  Handled bool  `json:"handled"`
  Timestamp time.Time `json:"timestamp"`
  ProcessGUID string `json:"process_guid"`
  ProcessID int `json:"process_id"`
  Image string `json:"image"`
  CommandLine string `json:"command_line"`
  ParentProcess string `json:"parent_process"`
  Details interface{}  `json:"details"`
}
