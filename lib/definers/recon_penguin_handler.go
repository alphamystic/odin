package definers

import (
  "github.com/alphamystic/odin/lib/utils"
)


type UniScans struct {
  ScanID string `json:"scanid,omitempty"`
  Name string `json:"name,omitempty"`
  ScanType string `json:"scantype,omitempty"`
  OwnerID string `json:"owner_id,omitempty"`
  utils.TimeStamps
}


type VulScanner struct {
    Name string
    Directory string
    VulnTypes []string // an array names of vulnerabilities to scan for
    //this way we can see if not in list do not scan(Could also just be one)
}
