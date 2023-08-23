//go:build windows
// +build windows

package bt

import (
  "os"
  "fmt"
	"time"
	"os/exec"
	"encoding/json"
)

type NTEvent struct {
  ID EventsID `json: "id"`
  Hash string `json:"hash"`
  Handled bool  `json:"handled"`
  Timestamp time.Time `json:"timestamp"`
  ProcessGUID string `json:"process_guid"`
  ProcessID int `json:"process_id"`
  Image string `json:"image"`
  CommandLine string `json:"command_line"`
  ParentProcess string `json:"parent_process"`
  Details interface{}  `json:"details"`
  CreatedAt string `json:"created_at"`
  UpdatedAt string  `json:"updated_at"`
}


/*
   A typical event should include it's asset it, useridprocess id, a dump of the in process memory, the pid,Action
   (quarantine or removed)
*/

type EventsID int
const (
  ProcessCreation EventsID = iota  //1
  ProcessChangedAFileCreationTime
  NetworkConnection
  SysmonServiceStateChanged
  ProcessTerminated
  DriverLoaded
  ImageLoaded
  CreateRemoteThread
  RawAccessRead
  ProcessAccess
  FileCreate
  RegistryEventCreationAndDeletion
  RegistryEventValueSet
  RegistryEventKeyAndValueRename
  FileCreateStreamHash
  ServiceConfigurationChange
  PipeEventCreation
  PipeEventConnected
  WmiEventWmiEventFilterActivity
  WmiEventWmiEventConsumerActivity
  WmiEventWmiEventConsumerToFilterActivity
  DNSEventDNSQuery
  FileDelete
  ClipboardChange
  ProcessTampering
  FileDeleteDetected
  FileBlockExecutable
  FileBlockShredding
  Error
)

// So I reaaly don't know if this works
// I want to print this out and see if I can Have a strict for each EventID
func readSysmonEvents() ([]NTEvent, error) {
  cmd := exec.Command("wevtutil", "qe", "Microsoft-Windows-Sysmon/Operational", "/f:json", "/q:*[System/EventID=1]")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	var events []NTEvent
	err = json.Unmarshal(output, &events)
	if err != nil {
		return nil, err
	}

	return events, nil
}
