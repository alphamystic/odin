//go:build windows
// +build windows

// https://github.com/0daysimpson/Get-SysmonLogs/blob/master/Get-SysmonLogs.ps1
// https://syedhasan010.medium.com/sysmon-how-to-setup-configure-and-analyze-the-system-monitors-events-930e9add78d
package nt

import (
	"fmt"
	"golang.org/x/sys/windows"
	"os"
	"syscall"
	"unsafe"
)

const (
	EVENTLOG_SEQUENTIAL_READ = 0x0001
	SYSMON_EVENT_LOG_NAME    = "Microsoft-Windows-Sysmon/Operational"
)

type SysmonEvent struct {
	EventID       uint32
	ProviderName  string
	TimeGenerated windows.Filetime
	Description   string
}

func Agent() {
	// Open the Sysmon event log.
	hEventLog, err := windows.OpenEventLog("", syscall.StringToUTF16Ptr(SYSMON_EVENT_LOG_NAME))
	if err != nil {
		fmt.Printf("Failed to open Sysmon event log: %v\n", err)
		os.Exit(1)
	}
	defer windows.CloseEventLog(hEventLog)

	// Set up variables to read event records.
	var bytesRead uint32
	var buffer [4096]byte // Event log record buffer size

	// Read Sysmon event records.
	for {
		err = windows.ReadEventLog(hEventLog, EVENTLOG_SEQUENTIAL_READ, 0, &buffer, uint32(len(buffer)), &bytesRead, nil)
		if err != nil {
			fmt.Printf("Failed to read Sysmon event log: %v\n", err)
			break
		}

		// Process the Sysmon event data.
		for offset := uint32(0); offset < bytesRead; {
			record := (*windows.EVENTLOGRECORD)(unsafe.Pointer(&buffer[offset]))
			eventData := buffer[offset+uintptr(record.StringOffset):]
			sysmonEvent := SysmonEvent{
				EventID:       record.EventID,
				ProviderName:  windows.UTF16ToString(eventData[:record.DataLength/2]),
				TimeGenerated: record.TimeGenerated,
				Description:   string(eventData[:record.DataLength]),
			}
			fmt.Printf("EventID: %d, ProviderName: %s, TimeGenerated: %v, Description: %s\n",
				sysmonEvent.EventID, sysmonEvent.ProviderName, sysmonEvent.TimeGenerated, sysmonEvent.Description)
			offset += record.Length
		}
	}
}
